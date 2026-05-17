package auth

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	firebase_auth "firebase.google.com/go/auth"
	mock_logger "github.com/downsized-devs/sdk-go/tests/mock/logger"
	mock_parser "github.com/downsized-devs/sdk-go/tests/mock/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// fakeFirebaseAuth captures calls and returns canned responses.
type fakeFirebaseAuth struct {
	getUsersRes *firebase_auth.GetUsersResult
	getUsersErr error
	getUsersIds []firebase_auth.UserIdentifier

	createRes    *firebase_auth.UserRecord
	createErr    error
	createParams *firebase_auth.UserToCreate

	updateRes    *firebase_auth.UserRecord
	updateErr    error
	updateUID    string
	updateParams *firebase_auth.UserToUpdate

	deleteErr error
	deleteUID string

	verifyTokenRes *firebase_auth.Token
	verifyTokenErr error
	verifyTokenIn  string

	revokeErr error
	revokeUID string
}

func (f *fakeFirebaseAuth) GetUsers(_ context.Context, identifiers []firebase_auth.UserIdentifier) (*firebase_auth.GetUsersResult, error) {
	f.getUsersIds = identifiers
	return f.getUsersRes, f.getUsersErr
}

func (f *fakeFirebaseAuth) CreateUser(_ context.Context, p *firebase_auth.UserToCreate) (*firebase_auth.UserRecord, error) {
	f.createParams = p
	return f.createRes, f.createErr
}

func (f *fakeFirebaseAuth) UpdateUser(_ context.Context, uid string, p *firebase_auth.UserToUpdate) (*firebase_auth.UserRecord, error) {
	f.updateUID = uid
	f.updateParams = p
	return f.updateRes, f.updateErr
}

func (f *fakeFirebaseAuth) DeleteUser(_ context.Context, uid string) error {
	f.deleteUID = uid
	return f.deleteErr
}

func (f *fakeFirebaseAuth) VerifyIDTokenAndCheckRevoked(_ context.Context, tok string) (*firebase_auth.Token, error) {
	f.verifyTokenIn = tok
	return f.verifyTokenRes, f.verifyTokenErr
}

func (f *fakeFirebaseAuth) RevokeRefreshTokens(_ context.Context, uid string) error {
	f.revokeUID = uid
	return f.revokeErr
}

func newAuthWithFake(t *testing.T, fake *fakeFirebaseAuth) *auth {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	return &auth{
		log:      mock_logger.NewMockInterface(ctrl),
		json:     mock_parser.NewMockJsonInterface(ctrl),
		firebase: fake,
		conf:     Config{},
	}
}

// ---------------- GetUser ---------------- //

func TestGetUser_ByID(t *testing.T) {
	rec := sampleUserRecord("uid-1", "alice@example.com")
	fake := &fakeFirebaseAuth{
		getUsersRes: &firebase_auth.GetUsersResult{Users: []*firebase_auth.UserRecord{rec}},
	}
	a := newAuthWithFake(t, fake)

	got, err := a.GetUser(context.Background(), FirebaseUserParam{ID: "uid-1"})
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "uid-1", got[0].ID)

	require.Len(t, fake.getUsersIds, 1)
	_, ok := fake.getUsersIds[0].(firebase_auth.UIDIdentifier)
	assert.True(t, ok)
}

func TestGetUser_ByEmail(t *testing.T) {
	fake := &fakeFirebaseAuth{
		getUsersRes: &firebase_auth.GetUsersResult{Users: []*firebase_auth.UserRecord{sampleUserRecord("uid-2", "bob@example.com")}},
	}
	a := newAuthWithFake(t, fake)

	_, err := a.GetUser(context.Background(), FirebaseUserParam{Email: "bob@example.com"})
	require.NoError(t, err)
	require.Len(t, fake.getUsersIds, 1)
	_, ok := fake.getUsersIds[0].(firebase_auth.EmailIdentifier)
	assert.True(t, ok)
}

func TestGetUser_ByPhone(t *testing.T) {
	fake := &fakeFirebaseAuth{
		getUsersRes: &firebase_auth.GetUsersResult{Users: []*firebase_auth.UserRecord{sampleUserRecord("uid-3", "")}},
	}
	a := newAuthWithFake(t, fake)

	_, err := a.GetUser(context.Background(), FirebaseUserParam{PhoneNumber: "+62812"})
	require.NoError(t, err)
	require.Len(t, fake.getUsersIds, 1)
	_, ok := fake.getUsersIds[0].(firebase_auth.PhoneIdentifier)
	assert.True(t, ok)
}

func TestGetUser_PropagatesError(t *testing.T) {
	fake := &fakeFirebaseAuth{getUsersErr: errors.New("backend down")}
	a := newAuthWithFake(t, fake)

	_, err := a.GetUser(context.Background(), FirebaseUserParam{ID: "x"})
	assert.Error(t, err)
}

// ---------------- GetUsers (plural) ---------------- //

func TestGetUsers_AllVariants(t *testing.T) {
	fake := &fakeFirebaseAuth{
		getUsersRes: &firebase_auth.GetUsersResult{Users: []*firebase_auth.UserRecord{
			sampleUserRecord("u1", "a@x"),
			sampleUserRecord("u2", "b@x"),
		}},
	}
	a := newAuthWithFake(t, fake)

	got, err := a.GetUsers(context.Background(), []FirebaseUserParam{
		{ID: "u1"},
		{Email: "b@x"},
		{PhoneNumber: "+1"},
	})
	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Len(t, fake.getUsersIds, 3)
}

func TestGetUsers_PropagatesError(t *testing.T) {
	fake := &fakeFirebaseAuth{getUsersErr: errors.New("oh no")}
	a := newAuthWithFake(t, fake)
	_, err := a.GetUsers(context.Background(), []FirebaseUserParam{{Email: "x"}})
	assert.Error(t, err)
}

// ---------------- RegisterUser ---------------- //

func TestRegisterUser_Success(t *testing.T) {
	fake := &fakeFirebaseAuth{
		createRes: sampleUserRecord("uid-new", "x@y"),
	}
	a := newAuthWithFake(t, fake)

	in := FirebaseUser{Email: "x@y", Password: "pw", DisplayName: "X", PhoneNumber: "+1", PhotoURL: "http://"}
	in.IsEmailVerified.Bool = true
	in.IsEmailVerified.Valid = true
	in.IsDisabled.Bool = false
	in.IsDisabled.Valid = true

	got, err := a.RegisterUser(context.Background(), in)
	require.NoError(t, err)
	assert.Equal(t, "uid-new", got.ID)
}

func TestRegisterUser_ForwardsIsDisabled(t *testing.T) {
	cases := []struct {
		name     string
		disabled bool
	}{
		{"disabled=true", true},
		{"disabled=false", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeFirebaseAuth{createRes: sampleUserRecord("uid", "x@y")}
			a := newAuthWithFake(t, fake)

			in := FirebaseUser{Email: "x@y", Password: "pw", DisplayName: "X"}
			// Email verified flipped opposite to ensure the wrong field isn't forwarded.
			in.IsEmailVerified.Valid = true
			in.IsEmailVerified.Bool = !tc.disabled
			in.IsDisabled.Valid = true
			in.IsDisabled.Bool = tc.disabled

			_, err := a.RegisterUser(context.Background(), in)
			require.NoError(t, err)

			got, ok := userToCreateParams(fake.createParams)["disabled"].(bool)
			require.True(t, ok, "disabled param should be set as a bool")
			assert.Equal(t, tc.disabled, got)
		})
	}
}

func TestRegisterUser_Error(t *testing.T) {
	fake := &fakeFirebaseAuth{createErr: errors.New("create failed")}
	a := newAuthWithFake(t, fake)

	_, err := a.RegisterUser(context.Background(), FirebaseUser{Email: "x"})
	assert.Error(t, err)
}

// ---------------- UpdateUser ---------------- //

func TestUpdateUser_Success(t *testing.T) {
	fake := &fakeFirebaseAuth{
		updateRes: sampleUserRecord("uid-1", "new@example.com"),
	}
	a := newAuthWithFake(t, fake)

	in := FirebaseUser{ID: "uid-1", Email: "new@example.com", PhoneNumber: "+1", DisplayName: "N", PhotoURL: "http://", Password: "pw"}
	in.IsEmailVerified.Valid = true
	in.IsEmailVerified.Bool = true
	in.IsDisabled.Valid = true

	got, err := a.UpdateUser(context.Background(), in)
	require.NoError(t, err)
	assert.Equal(t, "uid-1", got.ID)
	assert.Equal(t, "uid-1", fake.updateUID)
}

func TestUpdateUser_ForwardsIsDisabled(t *testing.T) {
	cases := []struct {
		name     string
		disabled bool
	}{
		{"disabled=true", true},
		{"disabled=false", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeFirebaseAuth{updateRes: sampleUserRecord("uid-1", "x@y")}
			a := newAuthWithFake(t, fake)

			in := FirebaseUser{ID: "uid-1", Email: "x@y"}
			in.IsEmailVerified.Valid = true
			in.IsEmailVerified.Bool = !tc.disabled
			in.IsDisabled.Valid = true
			in.IsDisabled.Bool = tc.disabled

			_, err := a.UpdateUser(context.Background(), in)
			require.NoError(t, err)

			got, ok := userToUpdateParams(fake.updateParams)["disableUser"].(bool)
			require.True(t, ok, "disableUser param should be set as a bool")
			assert.Equal(t, tc.disabled, got)
		})
	}
}

func TestUpdateUser_Error(t *testing.T) {
	fake := &fakeFirebaseAuth{updateErr: errors.New("update failed")}
	a := newAuthWithFake(t, fake)
	_, err := a.UpdateUser(context.Background(), FirebaseUser{ID: "x"})
	assert.Error(t, err)
}

// ---------------- DeleteUser ---------------- //

func TestDeleteUser_Success(t *testing.T) {
	fake := &fakeFirebaseAuth{}
	a := newAuthWithFake(t, fake)
	require.NoError(t, a.DeleteUser(context.Background(), "uid-1"))
	assert.Equal(t, "uid-1", fake.deleteUID)
}

func TestDeleteUser_Error(t *testing.T) {
	fake := &fakeFirebaseAuth{deleteErr: errors.New("nope")}
	a := newAuthWithFake(t, fake)
	assert.Error(t, a.DeleteUser(context.Background(), "uid-1"))
}

// ---------------- VerifyToken ---------------- //

func TestVerifyToken_Success(t *testing.T) {
	fake := &fakeFirebaseAuth{verifyTokenRes: &firebase_auth.Token{UID: "abc"}}
	a := newAuthWithFake(t, fake)
	got, err := a.VerifyToken(context.Background(), "tok")
	require.NoError(t, err)
	assert.Equal(t, "abc", got.UID)
}

func TestVerifyToken_Expired(t *testing.T) {
	fake := &fakeFirebaseAuth{verifyTokenErr: errors.New(expiredTokenMessage + ": foo")}
	a := authWithLoggerCapture(t, fake)
	_, err := a.VerifyToken(context.Background(), "tok")
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "expired"))
}

func TestVerifyToken_Revoked(t *testing.T) {
	fake := &fakeFirebaseAuth{verifyTokenErr: errors.New(revokedTokenMessage + ": foo")}
	a := authWithLoggerCapture(t, fake)
	_, err := a.VerifyToken(context.Background(), "tok")
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "revoked"))
}

func TestVerifyToken_OtherError(t *testing.T) {
	fake := &fakeFirebaseAuth{verifyTokenErr: errors.New("malformed token")}
	a := authWithLoggerCapture(t, fake)
	_, err := a.VerifyToken(context.Background(), "tok")
	assert.Error(t, err)
}

// ---------------- RevokeUserRefreshToken ---------------- //

func TestRevokeUserRefreshToken_Success(t *testing.T) {
	fake := &fakeFirebaseAuth{}
	a := newAuthWithFake(t, fake)
	require.NoError(t, a.RevokeUserRefreshToken(context.Background(), "uid-7"))
	assert.Equal(t, "uid-7", fake.revokeUID)
}

func TestRevokeUserRefreshToken_Error(t *testing.T) {
	fake := &fakeFirebaseAuth{revokeErr: errors.New("backend")}
	a := newAuthWithFake(t, fake)
	assert.Error(t, a.RevokeUserRefreshToken(context.Background(), "uid"))
}

// ---------------- helpers ---------------- //

func sampleUserRecord(uid, email string) *firebase_auth.UserRecord {
	return &firebase_auth.UserRecord{
		UserInfo: &firebase_auth.UserInfo{
			UID:         uid,
			Email:       email,
			PhoneNumber: "+1",
			DisplayName: "Sample",
			PhotoURL:    "http://img",
		},
		EmailVerified: true,
		UserMetadata: &firebase_auth.UserMetadata{
			CreationTimestamp:  1_000,
			LastLogInTimestamp: 2_000,
		},
	}
}

// userToCreateParams reads the unexported params map from a firebase
// UserToCreate so tests can assert which fields were forwarded.
func userToCreateParams(u *firebase_auth.UserToCreate) map[string]interface{} {
	return readParamsField(u)
}

// userToUpdateParams reads the unexported params map from a firebase
// UserToUpdate so tests can assert which fields were forwarded.
func userToUpdateParams(u *firebase_auth.UserToUpdate) map[string]interface{} {
	return readParamsField(u)
}

func readParamsField(v any) map[string]interface{} {
	rv := reflect.ValueOf(v).Elem().FieldByName("params")
	rv = reflect.NewAt(rv.Type(), unsafe.Pointer(rv.UnsafeAddr())).Elem()
	if rv.IsNil() {
		return nil
	}
	return rv.Interface().(map[string]interface{})
}

// authWithLoggerCapture builds an auth with a real-enough logger mock that
// permits .Error() calls (VerifyToken logs on every error branch).
func authWithLoggerCapture(t *testing.T, fake *fakeFirebaseAuth) *auth {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	log := mock_logger.NewMockInterface(ctrl)
	log.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	return &auth{
		log:      log,
		json:     mock_parser.NewMockJsonInterface(ctrl),
		firebase: fake,
		conf:     Config{},
	}
}
