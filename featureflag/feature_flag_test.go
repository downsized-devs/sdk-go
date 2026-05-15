package featureflag

import (
	"context"
	"testing"
	"time"

	"github.com/downsized-devs/sdk-go/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// staticRetriever returns the same []byte every call — minimal implementation of
// retriever.Retriever for the go-feature-flag client.
type staticRetriever struct {
	payload []byte
}

func (s *staticRetriever) Retrieve(_ context.Context) ([]byte, error) {
	return s.payload, nil
}

// Flag config in the legacy v0.x format the library expects.
const flagYAML = `
my-flag:
  percentage: 100
  true: true
  false: false
  default: true
`

func newDisabled(t *testing.T) Interface {
	t.Helper()
	return Init(Config{Enabled: false}, logger.Init(logger.Config{}), nil)
}

func newEnabled(t *testing.T) Interface {
	t.Helper()
	r := &staticRetriever{payload: []byte(flagYAML)}
	return Init(Config{Enabled: true, PollingInterval: time.Hour}, logger.Init(logger.Config{}), r)
}

// ------------------------- Disabled path ------------------------- //

func TestCheckUserFlags_DisabledReturnsFalse(t *testing.T) {
	ff := newDisabled(t)
	got, err := ff.CheckUserFlags("my-flag", "user-1", userLoginType, nil)
	require.NoError(t, err)
	assert.False(t, got)
}

func TestGetAllUserFlags_DisabledReturnsEmpty(t *testing.T) {
	ff := newDisabled(t)
	got, err := ff.GetAllUserFlags("user-1", userLoginType, nil)
	require.NoError(t, err)
	assert.Equal(t, []byte{}, got)
}

func TestRefresh_DisabledIsNoOp(t *testing.T) {
	ff := newDisabled(t)
	assert.NotPanics(t, func() { ff.Refresh() })
}

// ------------------------- Enabled path ------------------------- //

func TestCheckUserFlags_EnabledEvaluatesFlag(t *testing.T) {
	ff := newEnabled(t)
	got, err := ff.CheckUserFlags("my-flag", "user-1", userLoginType, nil)
	require.NoError(t, err)
	assert.True(t, got, "flag defaults to true via defaultRule:variation:Enabled")
}

func TestCheckUserFlags_EnabledUnknownFlagFallsBackToDefault(t *testing.T) {
	ff := newEnabled(t)
	got, err := ff.CheckUserFlags("does-not-exist", "user-1", userLoginType, nil)
	// go-feature-flag returns the fallback (false) and an error for unknown flags.
	_ = err
	assert.False(t, got)
}

func TestGetAllUserFlags_EnabledMarshalsJSON(t *testing.T) {
	ff := newEnabled(t)
	got, err := ff.GetAllUserFlags("user-1", userLoginType, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, got)
	assert.Contains(t, string(got), "my-flag")
}

// ------------------------- User-type branch coverage ------------------------- //

func Test_getUserType_Anonymous(t *testing.T) {
	ff := &featureFlag{}
	u := ff.getUserType("u-1", userAnonymousType, nil)
	assert.Equal(t, "u-1", u.GetKey())
	assert.True(t, u.IsAnonymous())
}

func Test_getUserType_Login(t *testing.T) {
	ff := &featureFlag{}
	u := ff.getUserType("u-2", userLoginType, nil)
	assert.Equal(t, "u-2", u.GetKey())
	assert.False(t, u.IsAnonymous())
}

func Test_getUserType_Custom(t *testing.T) {
	ff := &featureFlag{}
	u := ff.getUserType("u-3", userCustomType, map[string]any{"plan": "pro"})
	assert.Equal(t, "u-3", u.GetKey())
	assert.Equal(t, "pro", u.GetCustom()["plan"])
}

func Test_getUserType_Unknown(t *testing.T) {
	ff := &featureFlag{}
	u := ff.getUserType("u-4", "garbage", nil)
	// Unknown user type leaves the zero User.
	assert.Empty(t, u.GetKey())
}

func Test_registerCustomUser_EmptyMap(t *testing.T) {
	ff := &featureFlag{}
	u := ff.registerCustomUser("u-5", map[string]any{})
	assert.Equal(t, "u-5", u.GetKey())
	assert.Empty(t, u.GetCustom())
}
