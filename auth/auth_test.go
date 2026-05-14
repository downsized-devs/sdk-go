package auth

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	firebase_auth "firebase.google.com/go/auth"
	"github.com/downsized-devs/sdk-go/null"
	mock_logger "github.com/downsized-devs/sdk-go/tests/mock/logger"
	mock_parser "github.com/downsized-devs/sdk-go/tests/mock/parser"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// These unit tests cover the auth package code paths that don't require a
// real Firebase backend: the SkipFirebaseInit short-circuit branches, the
// pure-Go helpers (SetUserAuthInfo / GetUserAuthInfo / assignUser), and the
// exchangeRefreshToken HTTP path via an httptest server.

func newSkippedAuth(t *testing.T) *auth {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	log := mock_logger.NewMockInterface(ctrl)
	jsn := mock_parser.NewMockJsonInterface(ctrl)
	a := Init(Config{SkipFirebaseInit: true}, log, jsn, nil)
	return a.(*auth)
}

func TestInit_SkipFirebase(t *testing.T) {
	a := newSkippedAuth(t)
	assert.True(t, a.conf.SkipFirebaseInit)
	assert.Nil(t, a.firebase)
}

func TestAuth_SkipFirebaseInit_AllMethodsReturnNotImplemented(t *testing.T) {
	a := newSkippedAuth(t)
	ctx := context.Background()

	t.Run("SignInWithPassword", func(t *testing.T) {
		_, err := a.SignInWithPassword(ctx, UserLogin{Email: "a@b.com"})
		assert.Error(t, err)
	})
	t.Run("RefreshToken", func(t *testing.T) {
		_, err := a.RefreshToken(ctx, "tok")
		assert.Error(t, err)
	})
	t.Run("GetUser", func(t *testing.T) {
		users, err := a.GetUser(ctx, FirebaseUserParam{ID: "1"})
		assert.Error(t, err)
		assert.Empty(t, users)
	})
	t.Run("RegisterUser", func(t *testing.T) {
		_, err := a.RegisterUser(ctx, FirebaseUser{Email: "x"})
		assert.Error(t, err)
	})
	t.Run("UpdateUser", func(t *testing.T) {
		_, err := a.UpdateUser(ctx, FirebaseUser{ID: "1"})
		assert.Error(t, err)
	})
	t.Run("DeleteUser", func(t *testing.T) {
		err := a.DeleteUser(ctx, "uid")
		assert.Error(t, err)
	})
	t.Run("VerifyToken", func(t *testing.T) {
		_, err := a.VerifyToken(ctx, "tok")
		assert.Error(t, err)
	})
	t.Run("RevokeUserRefreshToken", func(t *testing.T) {
		err := a.RevokeUserRefreshToken(ctx, "uid")
		assert.Error(t, err)
	})
	t.Run("VerifyPassword", func(t *testing.T) {
		_, err := a.VerifyPassword(ctx, "a@b.com", "pw")
		assert.Error(t, err)
	})
	t.Run("GetUsers", func(t *testing.T) {
		users, err := a.GetUsers(ctx, []FirebaseUserParam{{Email: "x"}})
		assert.Error(t, err)
		assert.Empty(t, users)
	})
}

func TestAuth_UserAuthInfo_RoundTrip(t *testing.T) {
	a := newSkippedAuth(t)
	ctx := context.Background()

	token := &firebase_auth.Token{UID: "abc"}
	cred := &UserCredential{ID: 99, AccessToken: "at"}
	user := User{ID: 1, Name: "jack"}

	ctx2 := a.SetUserAuthInfo(ctx, UserAuthParam{User: user, FirebaseToken: token, UserCredential: cred})
	info, err := a.GetUserAuthInfo(ctx2)
	assert.NoError(t, err)
	assert.Equal(t, user, info.User)
	assert.Equal(t, *token, info.FirebaseToken)
	assert.Equal(t, *cred, info.UserCredential)
}

func TestAuth_SetUserAuthInfo_NilOptionalsAreZeroed(t *testing.T) {
	a := newSkippedAuth(t)
	ctx := a.SetUserAuthInfo(context.Background(), UserAuthParam{User: User{ID: 7}})
	info, err := a.GetUserAuthInfo(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(7), info.User.ID)
	assert.Equal(t, firebase_auth.Token{}, info.FirebaseToken)
	assert.Equal(t, UserCredential{}, info.UserCredential)
}

func TestAuth_GetUserAuthInfo_Missing(t *testing.T) {
	a := newSkippedAuth(t)
	_, err := a.GetUserAuthInfo(context.Background())
	assert.Error(t, err)
}

func TestAuth_GetUserAuthInfo_WrongTypeInContext(t *testing.T) {
	a := newSkippedAuth(t)
	ctx := context.WithValue(context.Background(), userAuthInfo, "not a UserAuthInfo")
	_, err := a.GetUserAuthInfo(ctx)
	assert.Error(t, err)
}

func TestAuth_assignUser(t *testing.T) {
	a := newSkippedAuth(t)
	rec := &firebase_auth.UserRecord{
		UserInfo: &firebase_auth.UserInfo{
			UID:         "uid-1",
			Email:       "jack@example.com",
			PhoneNumber: "+1234",
			DisplayName: "Jack",
			PhotoURL:    "http://img",
		},
		EmailVerified: true,
		Disabled:      false,
		UserMetadata: &firebase_auth.UserMetadata{
			CreationTimestamp:  3_000,
			LastLogInTimestamp: 5_000,
		},
	}
	got := a.assignUser(rec)
	assert.Equal(t, "uid-1", got.ID)
	assert.Equal(t, "jack@example.com", got.Email)
	assert.Equal(t, null.BoolFrom(true), got.IsEmailVerified)
	assert.Equal(t, null.BoolFrom(false), got.IsDisabled)
	assert.Equal(t, "+1234", got.PhoneNumber)
	assert.Equal(t, "Jack", got.DisplayName)
	assert.Equal(t, "http://img", got.PhotoURL)
	// Constructor stores nanoseconds but divides by 1000, so 3000 -> 3.
	assert.Equal(t, int64(3), got.CreationTimestamp)
	assert.Equal(t, int64(5), got.LastLoginTimestamp)
}

// --- exchangeRefreshToken tests using a fake HTTP transport ---

// roundTripFunc lets us inject a synthetic response without spinning up an
// httptest server (so the production URL constant is unchanged).
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func newAuthWithTransport(t *testing.T, rt http.RoundTripper, jsn *mock_parser.MockJsonInterface, log *mock_logger.MockInterface) *auth {
	t.Helper()
	return &auth{
		log:        log,
		json:       jsn,
		httpClient: &http.Client{Transport: rt},
		conf:       Config{Firebase: FirebaseConf{ApiKey: "key"}},
	}
}

func TestExchangeRefreshToken_MarshalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	jsn := mock_parser.NewMockJsonInterface(ctrl)
	log := mock_logger.NewMockInterface(ctrl)

	jsn.EXPECT().Marshal(gomock.Any()).Return(nil, errors.New("boom"))

	a := newAuthWithTransport(t, nil, jsn, log)
	_, err := a.exchangeRefreshToken(context.Background(), RefreshTokenRequest{GrantType: "refresh_token"})
	assert.Error(t, err)
}

func TestExchangeRefreshToken_HTTPError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	jsn := mock_parser.NewMockJsonInterface(ctrl)
	log := mock_logger.NewMockInterface(ctrl)

	jsn.EXPECT().Marshal(gomock.Any()).Return([]byte(`{}`), nil)

	a := newAuthWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("transport failed")
	}), jsn, log)

	_, err := a.exchangeRefreshToken(context.Background(), RefreshTokenRequest{})
	assert.Error(t, err)
}

func TestExchangeRefreshToken_Non200(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	jsn := mock_parser.NewMockJsonInterface(ctrl)
	log := mock_logger.NewMockInterface(ctrl)

	jsn.EXPECT().Marshal(gomock.Any()).Return([]byte(`{}`), nil)
	log.EXPECT().Error(gomock.Any(), gomock.Any()).Times(1)

	a := newAuthWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(strings.NewReader(`{"error":"bad"}`)),
			Header:     http.Header{},
		}, nil
	}), jsn, log)

	_, err := a.exchangeRefreshToken(context.Background(), RefreshTokenRequest{})
	assert.Error(t, err)
}

func TestExchangeRefreshToken_UnmarshalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	jsn := mock_parser.NewMockJsonInterface(ctrl)
	log := mock_logger.NewMockInterface(ctrl)

	jsn.EXPECT().Marshal(gomock.Any()).Return([]byte(`{}`), nil)
	jsn.EXPECT().Unmarshal(gomock.Any(), gomock.Any()).Return(errors.New("bad json"))

	a := newAuthWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`not json`)),
			Header:     http.Header{},
		}, nil
	}), jsn, log)

	_, err := a.exchangeRefreshToken(context.Background(), RefreshTokenRequest{})
	assert.Error(t, err)
}

func TestExchangeRefreshToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	jsn := mock_parser.NewMockJsonInterface(ctrl)
	log := mock_logger.NewMockInterface(ctrl)

	jsn.EXPECT().Marshal(gomock.Any()).Return([]byte(`{}`), nil)
	jsn.EXPECT().Unmarshal(gomock.Any(), gomock.Any()).DoAndReturn(func(_ []byte, dest interface{}) error {
		r := dest.(*RefreshTokenResponse)
		r.IDToken = "id-token"
		r.RefreshToken = "new-refresh"
		return nil
	})

	a := newAuthWithTransport(t, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		// Quick sanity check that the production constant URL is hit
		// (query string includes the API key).
		if !strings.HasPrefix(r.URL.String(), ExchangeRefreshTokenURL) {
			t.Errorf("unexpected URL: %s", r.URL)
		}
		if r.Header.Get(ContentType) != ApplicationJson {
			t.Errorf("missing or wrong content-type: %q", r.Header.Get(ContentType))
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"id_token":"id-token"}`)),
			Header:     http.Header{},
		}, nil
	}), jsn, log)

	got, err := a.exchangeRefreshToken(context.Background(), RefreshTokenRequest{RefreshToken: "old"})
	assert.NoError(t, err)
	assert.Equal(t, "id-token", got.IDToken)
	assert.Equal(t, "new-refresh", got.RefreshToken)
}

// Optional sanity: the constants are not accidentally renamed/reset.
func TestExchangeRefreshToken_Constants(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer srv.Close()
	assert.Equal(t, "application/json", ApplicationJson)
	assert.Equal(t, "Content-Type", ContentType)
	assert.True(t, strings.HasPrefix(ExchangeRefreshTokenURL, "https://"))
}
