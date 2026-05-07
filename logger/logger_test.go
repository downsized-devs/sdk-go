package logger

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/downsized-devs/sdk-go/appcontext"
	"github.com/downsized-devs/sdk-go/codes"
	"github.com/downsized-devs/sdk-go/errors"
	"github.com/stretchr/testify/assert"
)

func Test_logger_Trace(t *testing.T) {
	type args struct {
		ctx context.Context
		obj interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "trace",
			args: args{
				ctx: context.Background(),
				obj: int64(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logger{}
			l.Trace(tt.args.ctx, tt.args.obj)
		})
	}
}

func Test_logger_Debug(t *testing.T) {
	type args struct {
		ctx context.Context
		obj interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Debug",
			args: args{
				ctx: context.Background(),
				obj: int64(2),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logger{}
			l.Debug(tt.args.ctx, tt.args.obj)
		})
	}
}

func Test_logger_Info(t *testing.T) {
	type args struct {
		ctx context.Context
		obj interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "info",
			args: args{
				ctx: context.Background(),
				obj: int64(3),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logger{}
			l.Info(tt.args.ctx, tt.args.obj)
		})
	}
}

func Test_logger_Warn(t *testing.T) {
	type args struct {
		ctx context.Context
		obj interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "warn",
			args: args{
				ctx: context.Background(),
				obj: int64(4),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logger{}
			l.Warn(tt.args.ctx, tt.args.obj)
		})
	}
}

func Test_logger_Error(t *testing.T) {
	type args struct {
		ctx context.Context
		obj interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Error",
			args: args{
				ctx: context.Background(),
				obj: int64(5),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logger{}
			l.Error(tt.args.ctx, tt.args.obj)
		})
	}
}

func Test_logger_Fatal(t *testing.T) {
	type args struct {
		ctx context.Context
		obj interface{}
	}
	tests := []struct {
		name string
		args args
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logger{}
			l.Fatal(tt.args.ctx, tt.args.obj)
		})
	}
}

func TestInit(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantNil bool
	}{
		{
			name:    "debug level",
			cfg:     Config{Level: "debug"},
			wantNil: false,
		},
		{
			name:    "info level with custom skip",
			cfg:     Config{Level: "info", CallerSkipFrameCount: 4},
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Init(tt.cfg)
			if (got == nil) != tt.wantNil {
				t.Errorf("Init() returned nil=%v, want nil=%v", got == nil, tt.wantNil)
			}
		})
	}
}

func Test_getCaller(t *testing.T) {
	pwd, _ := os.Getwd()
	type args struct {
		obj interface{}
	}

	sdkErr := errors.NewWithCode(codes.CodeBadRequest, "test")
	_, _, sdkErrLine, _ := runtime.Caller(0)
	sdkErrLine-- // sdkErr was created on the line above

	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "get caller",
			args: args{
				obj: string(""),
			},
			want: string(""),
		},
		{
			name: "get caller error",
			args: args{
				obj: os.ErrInvalid,
			},
			want: os.ErrInvalid,
		},
		{
			name: "get caller sdk error",
			args: args{obj: sdkErr},
			want: fmt.Sprintf("%s/logger_test.go:%#v --- test", pwd, sdkErrLine),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getCaller(tt.args.obj)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_getContextFields(t *testing.T) {
	mockTime := time.Now()

	mockTime2 := mockTime.Add(2 * time.Second)

	now = func() time.Time {
		return mockTime2
	}

	restoreAll := func() {
		now = time.Now
	}

	defer restoreAll()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "get context fields time elapsed",
			args: args{ctx: appcontext.SetRequestStartTime(context.Background(), mockTime)},
			want: map[string]interface{}{
				"request_id":      "",
				"service_version": "",
				"user_id":         0,
				"time_elapsed":    "2000ms",
				"user_agent":      "",
			},
		},
		{
			name: "get context fields context default",
			args: args{ctx: context.Background()},
			want: map[string]interface{}{
				"request_id":      "",
				"service_version": "",
				"time_elapsed":    "0ms",
				"user_agent":      "",
				"user_id":         0,
			},
		},
		{
			name: "get context fields app response",
			args: args{ctx: appcontext.SetAppResponseCode(context.Background(), codes.CodeInvalidValue)},
			want: map[string]interface{}{
				"app_resp_code":   codes.CodeInvalidValue,
				"request_id":      "",
				"service_version": "",
				"time_elapsed":    "0ms",
				"user_agent":      "",
				"user_id":         0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getContextFields(tt.args.ctx)
			assert.Equal(t, tt.want, got)
		})
	}
}
