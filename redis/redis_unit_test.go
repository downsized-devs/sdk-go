package redis

import (
	"context"
	"testing"
	"time"

	"github.com/bsm/redislock"
	mock_log "github.com/downsized-devs/sdk-go/tests/mock/logger"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// Unit-test counterpart to redis_test.go (which is build-tagged and requires
// a live Redis on localhost:6379). These tests construct go-redis clients
// pointing at a closed/blackhole address so commands fail fast, exercising
// the SDK code paths without spinning up a real server.
//
// 127.0.0.1:1 is a port unlikely to have anything listening; the connect
// attempt fails almost immediately.
const blackholeAddr = "127.0.0.1:1"

func newBlackholeClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:        blackholeAddr,
		DialTimeout: 50 * time.Millisecond,
		MaxRetries:  -1,
	})
}

func newMockLogger(t *testing.T) *mock_log.MockInterface {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	log := mock_log.NewMockInterface(ctrl)
	log.EXPECT().Fatal(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	return log
}

func TestCache_GetDefaultTTL(t *testing.T) {
	ttl := 5 * time.Minute
	c := cache{conf: Config{DefaultTTL: ttl}}
	assert.Equal(t, ttl, c.GetDefaultTTL(context.Background()))
}

func TestCache_Get_Error(t *testing.T) {
	c := cache{rdb: newBlackholeClient()}
	_, err := c.Get(context.Background(), "missing")
	assert.Error(t, err)
}

func TestCache_SetEX_Error(t *testing.T) {
	c := cache{rdb: newBlackholeClient(), conf: Config{DefaultTTL: time.Minute}}
	err := c.SetEX(context.Background(), "k", "v", 0)
	assert.Error(t, err)
}

func TestCache_Lock_Error(t *testing.T) {
	rdb := newBlackholeClient()
	c := cache{rdb: rdb, rlock: redislock.New(rdb)}
	_, err := c.Lock(context.Background(), "k", time.Minute)
	assert.Error(t, err)
}

func TestCache_LockRelease_NilLock(t *testing.T) {
	c := cache{rdb: newBlackholeClient()}
	err := c.LockRelease(context.Background(), nil)
	assert.NoError(t, err)
}

func TestCache_Del_Error(t *testing.T) {
	c := cache{rdb: newBlackholeClient(), log: newMockLogger(t)}
	err := c.Del(context.Background(), "key*")
	assert.Error(t, err)
}

func TestCache_FlushAndPing_Errors(t *testing.T) {
	c := cache{rdb: newBlackholeClient()}
	ctx := context.Background()
	assert.Error(t, c.FlushAll(ctx))
	assert.Error(t, c.FlushAllAsync(ctx))
	assert.Error(t, c.FlushDB(ctx))
	assert.Error(t, c.FlushDBAsync(ctx))
	assert.Error(t, c.Ping(ctx))
}

func TestInit_ConnectFailure_DoesNotPanicWithMockedLogger(t *testing.T) {
	// Init calls connect() which calls log.Fatal on failure. The mock
	// logger does not exit, so Init returns a usable struct with an rdb
	// pointing at the bad address.
	log := newMockLogger(t)
	c := Init(Config{
		Protocol: "tcp",
		Host:     "127.0.0.1",
		Port:     "1",
	}, log)
	assert.NotNil(t, c)
}

func TestInit_TLS_Enabled(t *testing.T) {
	// Exercises the TLS-config branch in connect() (including the
	// InsecureSkipVerify warning).
	log := newMockLogger(t)
	c := Init(Config{
		Protocol: "tcp",
		Host:     "127.0.0.1",
		Port:     "1",
		TLS: TLSConfig{
			Enabled:            true,
			InsecureSkipVerify: true,
		},
	}, log)
	assert.NotNil(t, c)
}

func TestConstants(t *testing.T) {
	// Smoke check that the exported sentinel values pass through.
	assert.NotNil(t, ErrNotObtained)
	assert.Equal(t, redis.Nil, Nil)
}

func TestCRC16(t *testing.T) {
	// Verify the empty-string identity and that the function is
	// deterministic / non-zero for non-empty inputs.
	assert.Equal(t, uint16(0), CRC16(""))
	assert.NotEqual(t, uint16(0), CRC16("foo"))
	assert.Equal(t, CRC16("foo"), CRC16("foo"))
	// The widely-published check value 0x31C3 corresponds to the
	// "123456789" reference input under CRC16-CCITT/XModem.
	assert.Equal(t, uint16(0x31C3), CRC16("123456789"))
}
