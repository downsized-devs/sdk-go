package sql

import (
	"database/sql"
	"testing"

	mock_log "github.com/downsized-devs/sdk-go/tests/mock/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newMemSqliteDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// TestLeaderFollowerRoutingWithMockDB verifies that when both leader and
// follower MockDBs are injected, Leader() returns the leader mock and
// Follower() returns the follower mock.
func TestLeaderFollowerRoutingWithMockDB(t *testing.T) {
	leaderDB := newMemSqliteDB(t)
	followerDB := newMemSqliteDB(t)

	ctrl := gomock.NewController(t)
	log := mock_log.NewMockInterface(ctrl)
	log.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Fatal(gomock.Any(), gomock.Any()).AnyTimes()

	db := Init(Config{
		Driver: "sqlmock",
		Name:   "test",
		Leader: ConnConfig{
			Host:   "leader-host",
			Port:   1111,
			MockDB: leaderDB,
		},
		Follower: ConnConfig{
			Host:   "follower-host",
			Port:   2222,
			MockDB: followerDB,
		},
	}, log, nil)

	leaderCmd := db.Leader()
	followerCmd := db.Follower()
	require.NotNil(t, leaderCmd)
	require.NotNil(t, followerCmd)

	leaderInner, ok := leaderCmd.(*command)
	require.True(t, ok)
	followerInner, ok := followerCmd.(*command)
	require.True(t, ok)

	assert.Same(t, leaderDB, leaderInner.db.DB, "Leader() should wrap the leader MockDB")
	assert.Same(t, followerDB, followerInner.db.DB, "Follower() should wrap the follower MockDB")
}

// TestIsFollowerEnabled covers the host/port detection edge cases that 1c
// addresses, including "both host and port differ".
func TestIsFollowerEnabled(t *testing.T) {
	tests := []struct {
		name   string
		leader ConnConfig
		fol    ConnConfig
		want   bool
	}{
		{"no follower host", ConnConfig{Host: "a", Port: 1}, ConnConfig{Host: "", Port: 1}, false},
		{"identical addr", ConnConfig{Host: "a", Port: 1}, ConnConfig{Host: "a", Port: 1}, false},
		{"different host only", ConnConfig{Host: "a", Port: 1}, ConnConfig{Host: "b", Port: 1}, true},
		{"different port only", ConnConfig{Host: "a", Port: 1}, ConnConfig{Host: "a", Port: 2}, true},
		{"both host and port differ", ConnConfig{Host: "a", Port: 1}, ConnConfig{Host: "b", Port: 2}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sqlDB{cfg: Config{Leader: tt.leader, Follower: tt.fol}}
			assert.Equal(t, tt.want, s.isFollowerEnabled())
		})
	}
}
