//go:build integration
// +build integration

package redis

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/bsm/redislock"
	mock_log "github.com/downsized-devs/sdk-go/tests/mock/logger"
	"github.com/go-redis/redis/v8"
	"go.uber.org/mock/gomock"
)

func Test_redis_Get(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Fatal(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name          string
		args          args
		prepCacheMock func() *redis.Client
		want          interface{}
		wantErr       bool
	}{
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				key: "test1",
			},
			prepCacheMock: func() *redis.Client {
				db := redis.NewClient(&redis.Options{
					Addr:     "localhost:6379",
					Username: "",
					Password: "",
				})
				return db
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				key: "test1",
			},
			prepCacheMock: func() *redis.Client {
				db := redis.NewClient(&redis.Options{
					Addr:     "localhost:6379",
					Username: "",
					Password: "",
				})

				db.Set(context.Background(), "test1", "test1", time.Hour)

				return db
			},
			want:    "test1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		db := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Username: "",
			Password: "",
		})
		t.Run(tt.name, func(t *testing.T) {
			rdb := tt.prepCacheMock()
			c := cache{
				rdb: rdb,
			}
			got, err := c.Get(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("cache.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("cache.Get() = %v, want %v", got, tt.want)
			}
			db.Del(context.Background(), "test1")
		})
	}
}

func Test_redis_SetEX(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTime := time.Hour * 24

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Fatal(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	type args struct {
		ctx     context.Context
		key     string
		val     string
		expTime time.Duration
	}
	tests := []struct {
		name          string
		args          args
		prepCacheMock func() *redis.Client
		wantErr       bool
	}{
		{
			name: "error",
			args: args{
				ctx:     context.Background(),
				key:     "testset",
				val:     "yes",
				expTime: 0,
			},
			prepCacheMock: func() *redis.Client {
				db := redis.NewClient(&redis.Options{
					Addr:     "localhost:6378",
					Username: "",
					Password: "",
				})
				return db
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx:     context.Background(),
				key:     "whatever",
				val:     "yes",
				expTime: mockTime,
			},
			prepCacheMock: func() *redis.Client {
				db := redis.NewClient(&redis.Options{
					Addr:     "localhost:6379",
					Username: "",
					Password: "",
				})
				return db
			},
			wantErr: false,
		},
		{
			name: "success with default ttl",
			args: args{
				ctx: context.Background(),
				key: "whatever",
				val: "yes",
			},
			prepCacheMock: func() *redis.Client {
				db := redis.NewClient(&redis.Options{
					Addr:     "localhost:6379",
					Username: "",
					Password: "",
				})
				return db
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		db := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Username: "",
			Password: "",
		})
		t.Run(tt.name, func(t *testing.T) {
			rdb := tt.prepCacheMock()
			c := cache{
				conf: Config{
					DefaultTTL: mockTime,
				},
				rdb: rdb,
			}
			if err := c.SetEX(tt.args.ctx, tt.args.key, tt.args.val, tt.args.expTime); (err != nil) != tt.wantErr {
				t.Errorf("cache.SetEX() error = %v, wantErr %v", err, tt.wantErr)
			}
			db.Del(context.Background(), "testset")
		})
	}
}

func Test_redis_Lock(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTime := time.Minute * 24

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Fatal(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	type args struct {
		ctx     context.Context
		key     string
		expTime time.Duration
	}

	tests := []struct {
		name          string
		args          args
		prepCacheMock func() (*redis.Client, *redislock.Client)
		want          string
		wantErr       bool
	}{
		{
			name: "error",
			args: args{
				ctx:     context.Background(),
				key:     "testlock",
				expTime: 0,
			},
			prepCacheMock: func() (*redis.Client, *redislock.Client) {
				db := redis.NewClient(&redis.Options{
					Addr:     "localhost:6379",
					Username: "",
					Password: "",
				})

				locker := redislock.New(db)
				db.Set(context.Background(), "testlock", "testlock", time.Hour)

				return db, locker
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx:     context.Background(),
				key:     "galgadot",
				expTime: mockTime,
			},
			prepCacheMock: func() (*redis.Client, *redislock.Client) {
				db := redis.NewClient(&redis.Options{
					Addr:     "localhost:6379",
					Username: "",
					Password: "",
				})

				locker := redislock.New(db)
				return db, locker
			},
			want:    "galgadot",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdb, lck := tt.prepCacheMock()
			c := cache{
				rdb:   rdb,
				rlock: lck,
			}
			got, err := c.Lock(tt.args.ctx, tt.args.key, tt.args.expTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("cache.Lock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(got.Key(), tt.want) {
				t.Errorf("cache.Lock() = %v, want %v", got, tt.want)
			}
			got.Release(context.Background())
		})
	}
}

func Test_redis_LockRelease(t *testing.T) {
	mockTime := time.Minute * 24
	type args struct {
		ctx  context.Context
		lock *redislock.Lock
	}
	tests := []struct {
		name          string
		args          args
		prepCacheMock func() (*redis.Client, *redislock.Client)
		prepLockMock  func() *redislock.Lock
		wantErr       bool
	}{
		{
			name: "error",
			args: args{
				ctx:  context.Background(),
				lock: &redislock.Lock{},
			},
			prepCacheMock: func() (*redis.Client, *redislock.Client) {
				db := redis.NewClient(&redis.Options{
					Addr:     "localhost:6379",
					Username: "",
					Password: "",
				})

				locker := redislock.New(db)

				return db, locker
			},
			prepLockMock: func() *redislock.Lock {
				db := redis.NewClient(&redis.Options{
					Addr:     "localhost:6379",
					Username: "",
					Password: "",
				})

				locker := redislock.New(db)
				j, err := locker.Obtain(context.Background(), "tyty", mockTime, nil)
				fmt.Println(err)
				j.Release(context.Background())
				return j
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx:  context.Background(),
				lock: &redislock.Lock{},
			},
			prepCacheMock: func() (*redis.Client, *redislock.Client) {
				db := redis.NewClient(&redis.Options{
					Addr:     "localhost:6379",
					Username: "",
					Password: "",
				})

				locker := redislock.New(db)

				return db, locker
			},
			prepLockMock: func() *redislock.Lock {
				db := redis.NewClient(&redis.Options{
					Addr:     "localhost:6379",
					Username: "",
					Password: "",
				})

				locker := redislock.New(db)
				j, _ := locker.Obtain(context.Background(), "test", mockTime, nil)
				return j
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rdb, lck := tt.prepCacheMock()
			c := cache{
				rdb:   rdb,
				rlock: lck,
			}
			k := tt.prepLockMock()

			if err := c.LockRelease(tt.args.ctx, k); (err != nil) != tt.wantErr {
				t.Errorf("cache.LockRelease() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cache_Del(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name          string
		args          args
		prepCacheMock func() *redis.Client
		prepIterMock  func() string
		wantErr       bool
	}{
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
				key: "skey1",
			},
			prepCacheMock: func() *redis.Client {
				db := redis.NewClient(&redis.Options{
					Addr:     "localhost:6379",
					Username: "",
					Password: "",
				})
				return db
			},
			prepIterMock: func() string {
				db := redis.NewClient(&redis.Options{
					Addr:     "localhost:6379",
					Username: "",
					Password: "",
				})
				db.Set(context.Background(), "skey1", "key1val", time.Hour)
				x := db.Scan(context.Background(), 0, "skey1", 0).Iterator()
				x.Next(context.Background())
				y := x.Val()
				return y
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdb := tt.prepCacheMock()
			c := cache{
				rdb: rdb,
				log: logger,
			}
			str := tt.prepIterMock()
			if err := c.Del(tt.args.ctx, str); (err != nil) != tt.wantErr {
				t.Errorf("cache.Del() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_redis_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name          string
		args          args
		prepCacheMock func() *redis.Client
		wantErr       bool
	}{
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
			},
			prepCacheMock: func() *redis.Client {
				db := redis.NewClient(&redis.Options{
					Addr:     "localhost:6379",
					Username: "",
					Password: "",
				})
				return db
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdb := tt.prepCacheMock()

			c := cache{
				rdb: rdb,
				log: logger,
			}
			if err := c.Ping(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("cache.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}
