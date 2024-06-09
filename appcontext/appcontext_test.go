package appcontext

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/downsized-devs/sdk-go/codes"
	"github.com/downsized-devs/sdk-go/language"
	"github.com/stretchr/testify/assert"
)

func TestSetAcceptLanguage(t *testing.T) {
	type args struct {
		ctx      context.Context
		language string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), language: language.English},
			want: context.WithValue(context.Background(), acceptLanguage, language.English),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetAcceptLanguage(tt.args.ctx, tt.args.language); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetAcceptLanguage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAcceptLanguage(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "not ok",
			args: args{ctx: context.Background()},
			want: language.English,
		},
		{
			name: "ok",
			args: args{ctx: context.WithValue(context.Background(), acceptLanguage, language.Indonesian)},
			want: language.Indonesian,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAcceptLanguage(tt.args.ctx); got != tt.want {
				t.Errorf("GetAcceptLanguage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetRequestId(t *testing.T) {
	type args struct {
		ctx context.Context
		rid string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), rid: "randomized-request-id"},
			want: context.WithValue(context.Background(), requestId, "randomized-request-id"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetRequestId(tt.args.ctx, tt.args.rid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetRequestId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRequestId(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "not ok",
			args: args{ctx: context.Background()},
			want: "",
		},
		{
			name: "ok",
			args: args{ctx: context.WithValue(context.Background(), requestId, "randomized-request-id")},
			want: "randomized-request-id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRequestId(tt.args.ctx); got != tt.want {
				t.Errorf("GetRequestId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetServiceVersion(t *testing.T) {
	type args struct {
		ctx     context.Context
		version string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), version: "v1.0.0"},
			want: context.WithValue(context.Background(), serviceVersion, "v1.0.0"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetServiceVersion(tt.args.ctx, tt.args.version); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetServiceVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetServiceVersion(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "not ok",
			args: args{ctx: context.Background()},
			want: "",
		},
		{
			name: "ok",
			args: args{ctx: context.WithValue(context.Background(), serviceVersion, "v1.0.0")},
			want: "v1.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetServiceVersion(tt.args.ctx); got != tt.want {
				t.Errorf("GetServiceVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetUserAgent(t *testing.T) {
	type args struct {
		ctx context.Context
		ua  string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), ua: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko)"},
			want: context.WithValue(context.Background(), userAgent, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko)"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetUserAgent(tt.args.ctx, tt.args.ua); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetUserAgent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUserAgent(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "not ok",
			args: args{ctx: context.Background()},
			want: "",
		},
		{
			name: "ok",
			args: args{ctx: context.WithValue(context.Background(), userAgent, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko)")},
			want: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUserAgent(tt.args.ctx); got != tt.want {
				t.Errorf("GetUserAgent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetRequestStartTime(t *testing.T) {
	mocktime := time.Now()
	type args struct {
		ctx context.Context
		t   time.Time
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), t: mocktime},
			want: context.WithValue(context.Background(), requestStartTime, mocktime),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetRequestStartTime(tt.args.ctx, tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetRequestStartTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRequestStartTime(t *testing.T) {
	mocktime := time.Now()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "not ok",
			args: args{ctx: context.Background()},
		},
		{
			name: "not ok",
			args: args{ctx: context.WithValue(context.Background(), requestStartTime, mocktime)},
			want: mocktime,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRequestStartTime(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRequestStartTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetAppResponseCode(t *testing.T) {
	type args struct {
		ctx  context.Context
		code codes.Code
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), code: codes.CodeInvalidValue},
			want: context.WithValue(context.Background(), appResponseCode, codes.CodeInvalidValue),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetAppResponseCode(tt.args.ctx, tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetApplicationCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppResponseCode(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want codes.Code
	}{
		{
			name: "not ok",
			args: args{ctx: context.Background()},
		},
		{
			name: "ok",
			args: args{ctx: context.WithValue(context.Background(), appResponseCode, codes.CodeInvalidValue)},
			want: codes.CodeInvalidValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAppResponseCode(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetApplicationCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetDeviceType(t *testing.T) {
	type args struct {
		ctx      context.Context
		platform string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), platform: "web"},
			want: context.WithValue(context.Background(), deviceType, "web"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetDeviceType(tt.args.ctx, tt.args.platform); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetDeviceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDeviceType(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "not ok",
			args: args{ctx: context.Background()},
			want: "web",
		},
		{
			name: "ok",
			args: args{ctx: context.WithValue(context.Background(), deviceType, "app")},
			want: "app",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDeviceType(tt.args.ctx); got != tt.want {
				t.Errorf("GetDeviceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetAppErrorMessage(t *testing.T) {
	type args struct {
		ctx    context.Context
		errMsg string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), errMsg: "Error Message"},
			want: context.WithValue(context.Background(), appErrorMessage, "Error Message"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetAppErrorMessage(tt.args.ctx, tt.args.errMsg); !assert.Equal(t, tt.want, got) {
				t.Errorf("SetAppErrorMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppErrorMessage(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "not ok",
			args: args{ctx: context.Background()},
			want: "",
		},
		{
			name: "ok",
			args: args{ctx: context.WithValue(context.Background(), appErrorMessage, "Error Message")},
			want: "Error Message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAppErrorMessage(tt.args.ctx); got != tt.want {
				t.Errorf("GetAppErrorMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetEventName(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), name: "create"},
			want: context.WithValue(context.Background(), eventName, "create"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetEventName(tt.args.ctx, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetEventName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEventName(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty event name",
			args: args{ctx: context.Background()},
			want: "",
		},
		{
			name: "ok",
			args: args{ctx: context.WithValue(context.Background(), eventName, "create")},
			want: "create",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEventName(tt.args.ctx); got != tt.want {
				t.Errorf("GetEventName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetEventDescription(t *testing.T) {
	type args struct {
		ctx  context.Context
		desc string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), desc: "create test"},
			want: context.WithValue(context.Background(), eventDescription, "create test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetEventDescription(tt.args.ctx, tt.args.desc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetEventDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEventDescription(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty event desc",
			args: args{ctx: context.Background()},
			want: "",
		},
		{
			name: "ok",
			args: args{ctx: context.WithValue(context.Background(), eventDescription, "create")},
			want: "create",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEventDescription(tt.args.ctx); got != tt.want {
				t.Errorf("GetEventDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetRequestBody(t *testing.T) {
	type args struct {
		ctx  context.Context
		body string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), body: "{'data': 'test'}"},
			want: context.WithValue(context.Background(), requestBody, "{'data': 'test'}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetRequestBody(tt.args.ctx, tt.args.body); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetRequestBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRequestBody(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "empty body",
			args: args{ctx: context.Background()},
			want: nil,
		},
		{
			name: "ok",
			args: args{ctx: context.WithValue(context.Background(), requestBody, "{'data': 'test'}")},
			want: "{'data': 'test'}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRequestBody(tt.args.ctx); got != tt.want {
				t.Errorf("GetRequestBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetRequestURI(t *testing.T) {
	type args struct {
		ctx context.Context
		uri string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), uri: "v1/test"},
			want: context.WithValue(context.Background(), requestURI, "v1/test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetRequestURI(tt.args.ctx, tt.args.uri); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetRequestURI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRequestURI(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty uri",
			args: args{ctx: context.Background()},
			want: "",
		},
		{
			name: "ok",
			args: args{ctx: context.WithValue(context.Background(), requestURI, "v1/test")},
			want: "v1/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRequestURI(tt.args.ctx); got != tt.want {
				t.Errorf("GetRequestURI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetRequestQuery(t *testing.T) {
	type args struct {
		ctx   context.Context
		query string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), query: "test=data"},
			want: context.WithValue(context.Background(), requestQuery, "test=data"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetRequestQuery(tt.args.ctx, tt.args.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetRequestQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRequestQuery(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty query",
			args: args{ctx: context.Background()},
			want: "",
		},
		{
			name: "ok",
			args: args{ctx: context.WithValue(context.Background(), requestQuery, "test=data")},
			want: "test=data",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRequestQuery(tt.args.ctx); got != tt.want {
				t.Errorf("GetRequestQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetRequestMethod(t *testing.T) {
	type args struct {
		ctx    context.Context
		method string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), method: "POST"},
			want: context.WithValue(context.Background(), requestMethod, "POST"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetRequestMethod(tt.args.ctx, tt.args.method); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetRequestMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRequestMethod(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty method",
			args: args{ctx: context.Background()},
			want: "",
		},
		{
			name: "ok",
			args: args{ctx: context.WithValue(context.Background(), requestMethod, "PUT")},
			want: "PUT",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRequestMethod(tt.args.ctx); got != tt.want {
				t.Errorf("GetRequestMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetRequestIP(t *testing.T) {
	type args struct {
		ctx context.Context
		ip  string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), ip: "127.0.0.1"},
			want: context.WithValue(context.Background(), clientIP, "127.0.0.1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetRequestIP(tt.args.ctx, tt.args.ip); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetRequestIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRequestIP(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty ip",
			args: args{ctx: context.Background()},
			want: "",
		},
		{
			name: "ok",
			args: args{ctx: context.WithValue(context.Background(), clientIP, "127.0.0.1")},
			want: "127.0.0.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRequestIP(tt.args.ctx); got != tt.want {
				t.Errorf("GetRequestIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetResponseHttpCode(t *testing.T) {
	type args struct {
		ctx  context.Context
		code int
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), code: 200},
			want: context.WithValue(context.Background(), responseHttpCode, 200),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetResponseHttpCode(tt.args.ctx, tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetResponseHttpCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetResponseHttpCode(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "empty code",
			args: args{ctx: context.Background()},
			want: 0,
		},
		{
			name: "ok",
			args: args{ctx: context.WithValue(context.Background(), responseHttpCode, 200)},
			want: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetResponseHttpCode(tt.args.ctx); got != tt.want {
				t.Errorf("GetResponseHttpCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetUserId(t *testing.T) {
	type args struct {
		ctx context.Context
		ui  int
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), ui: 1},
			want: context.WithValue(context.Background(), userId, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetUserId(tt.args.ctx, tt.args.ui); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetUserId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUserId(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "not ok",
			args: args{ctx: context.Background()},
			want: 0,
		},
		{
			name: "ok",
			args: args{ctx: context.WithValue(context.Background(), userId, 1)},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUserId(tt.args.ctx); got != tt.want {
				t.Errorf("GetUserId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetAuthToken(t *testing.T) {
	type args struct {
		ctx   context.Context
		token string
	}
	mockToken := "eyJhbGciOiJSUzI1NiIsImtpZCI6Ijg3YzFlN2Y4MDAzNGJiYzgxYjhmMmRiODM3OTIxZjRiZDI4N2YxZGYiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoiSm9oYW5uZXMgV2FydWh1IiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hL0FMbTV3dTN0THNYU2QweXA5bFFlOVVuOVltQW5jSXBKckxzYnR5ODBtWElwPXM5Ni1jIiwiaXNzIjoiaHR0cHM6Ly9zZWN1cmV0b2tlbi5nb29nbGUuY29tL2RlbG9zLWFxdWFoZXJvLXN0YWdpbmciLCJhdWQiOiJkZWxvcy1hcXVhaGVyby1zdGFnaW5nIiwiYXV0aF90aW1lIjoxNjgwNDg4MjY3LCJ1c2VyX2lkIjoiYTVtNU5aakxjOFllRldWTENLQjVDbVFXdWMyMyIsInN1YiI6ImE1bTVOWmpMYzhZZUZXVkxDS0I1Q21RV3VjMjMiLCJpYXQiOjE2ODA0OTczMTAsImV4cCI6MTY4MDUwMDkxMCwiZW1haWwiOiJqb2hhbm5lcy53YXJ1aHVAZGVsb3NhcXVhLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJmaXJlYmFzZSI6eyJpZGVudGl0aWVzIjp7Imdvb2dsZS5jb20iOlsiMTExNDU1NTYyMjk3ODA5NDQzMzkwIl0sImVtYWlsIjpbImpvaGFubmVzLndhcnVodUBkZWxvc2FxdWEuY29tIl19LCJzaWduX2luX3Byb3ZpZGVyIjoiZ29vZ2xlLmNvbSJ9fQ.ciC6RpQzz0Wzd8NZetmi1PW8jFrRmzEVaJZqPhCrcyQW-O0G3RABmmCWIte6s5-WdF422TQPWr2U2NleIVSqh4-TiTckGVo8RaevipT5rmeBvzqWyYvM5k9cqadcYEtTW6HHMRBHuG9YP1ShY3zqFShnm6qDyCrWTzxgjKxVT2gq6Yrm-Ljrwwhjv03Ho3G4ljo5QPo0hHmSg7c_5BClxHrGJVI702Axbt7HjAB_NUddbMZE8q_-3a3FTYZyRvM4yrn0PoOz-yRLTgw_lTb50ftmp1WhDgr9POIvs6ECx4F8mYs_GJHM5l14DXhVymqwpVEUt9GXgNvzP0ZHFFfnwA"
	mockContextWithValue := context.WithValue(context.Background(), authToken, mockToken)
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "success set auth token",
			args: args{
				ctx:   context.Background(),
				token: mockToken,
			},
			want: mockContextWithValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetAuthToken(tt.args.ctx, tt.args.token); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetAuthToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAuthToken(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	mockToken := "eyJhbGciOiJSUzI1NiIsImtpZCI6Ijg3YzFlN2Y4MDAzNGJiYzgxYjhmMmRiODM3OTIxZjRiZDI4N2YxZGYiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoiSm9oYW5uZXMgV2FydWh1IiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hL0FMbTV3dTN0THNYU2QweXA5bFFlOVVuOVltQW5jSXBKckxzYnR5ODBtWElwPXM5Ni1jIiwiaXNzIjoiaHR0cHM6Ly9zZWN1cmV0b2tlbi5nb29nbGUuY29tL2RlbG9zLWFxdWFoZXJvLXN0YWdpbmciLCJhdWQiOiJkZWxvcy1hcXVhaGVyby1zdGFnaW5nIiwiYXV0aF90aW1lIjoxNjgwNDg4MjY3LCJ1c2VyX2lkIjoiYTVtNU5aakxjOFllRldWTENLQjVDbVFXdWMyMyIsInN1YiI6ImE1bTVOWmpMYzhZZUZXVkxDS0I1Q21RV3VjMjMiLCJpYXQiOjE2ODA0OTczMTAsImV4cCI6MTY4MDUwMDkxMCwiZW1haWwiOiJqb2hhbm5lcy53YXJ1aHVAZGVsb3NhcXVhLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJmaXJlYmFzZSI6eyJpZGVudGl0aWVzIjp7Imdvb2dsZS5jb20iOlsiMTExNDU1NTYyMjk3ODA5NDQzMzkwIl0sImVtYWlsIjpbImpvaGFubmVzLndhcnVodUBkZWxvc2FxdWEuY29tIl19LCJzaWduX2luX3Byb3ZpZGVyIjoiZ29vZ2xlLmNvbSJ9fQ.ciC6RpQzz0Wzd8NZetmi1PW8jFrRmzEVaJZqPhCrcyQW-O0G3RABmmCWIte6s5-WdF422TQPWr2U2NleIVSqh4-TiTckGVo8RaevipT5rmeBvzqWyYvM5k9cqadcYEtTW6HHMRBHuG9YP1ShY3zqFShnm6qDyCrWTzxgjKxVT2gq6Yrm-Ljrwwhjv03Ho3G4ljo5QPo0hHmSg7c_5BClxHrGJVI702Axbt7HjAB_NUddbMZE8q_-3a3FTYZyRvM4yrn0PoOz-yRLTgw_lTb50ftmp1WhDgr9POIvs6ECx4F8mYs_GJHM5l14DXhVymqwpVEUt9GXgNvzP0ZHFFfnwA"
	mockContextWithValue := context.WithValue(context.Background(), authToken, mockToken)
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success get auth token",
			args: args{
				ctx: mockContextWithValue,
			},
			want: mockToken,
		},
		{
			name: "failed/ no auth token ",
			args: args{
				ctx: context.Background(),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAuthToken(tt.args.ctx); got != tt.want {
				t.Errorf("GetAuthToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetServiceName(t *testing.T) {
	type args struct {
		ctx context.Context
		val string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "set service name",
			args: args{
				ctx: context.Background(),
				val: "generic-service",
			},
			want: context.WithValue(context.Background(), serviceName, "generic-service"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetServiceName(tt.args.ctx, tt.args.val); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetServiceName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetServiceName(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get service error",
			args: args{
				ctx: context.Background(),
			},
			want: "",
		},
		{
			name: "success get service",
			args: args{
				ctx: context.WithValue(context.Background(), serviceName, "generic-service"),
			},
			want: "generic-service",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetServiceName(tt.args.ctx); got != tt.want {
				t.Errorf("GetServiceName() = %v, want %v", got, tt.want)
			}
		})
	}
}
