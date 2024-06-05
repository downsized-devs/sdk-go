package ratelimiter

import (
	"reflect"
	"testing"

	"github.com/downsized-devs/sdk-go/log"
	mock_log "github.com/downsized-devs/sdk-go/tests/mock/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func Test_Limiter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

	configTrue := Config{
		Enabled: true,
		Period:  "1s",
		Limit:   1,
		Paths: []ConfigPath{
			{
				Enabled: true,
				Period:  "1m",
				Limit:   2,
				Path:    "/test",
			},
		},
	}

	configFalse := Config{
		Enabled: false,
	}

	type args struct {
		log log.Interface
		cfg Config
	}

	tests := []struct {
		name string
		args args
		want gin.HandlerFunc
	}{
		{
			name: "rate limiter false",
			args: args{
				log: logger,
				cfg: configFalse,
			},
			want: Init(configFalse, logger).Limiter(),
		},
		{
			name: "rate limiter true",
			args: args{
				log: logger,
				cfg: configTrue,
			},
			want: Init(configTrue, logger).Limiter(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl := Init(tt.args.cfg, tt.args.log)

			if got := rl.Limiter(); reflect.ValueOf(got).Pointer() != reflect.ValueOf(tt.want).Pointer() {
				t.Errorf("Limiter() = %v, want %v", got, tt.want)
			}
		})
	}
}
