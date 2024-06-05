package errors

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/downsized-devs/sdk-go/codes"
	"github.com/downsized-devs/sdk-go/language"
)

func TestApp_Error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "OK",
			want: "invalid format",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &App{
				sys: NewWithCode(codes.CodeBadRequest, "invalid format"),
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("App.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompile(t *testing.T) {
	type args struct {
		err  error
		lang string
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 App
	}{
		{
			name: "ok",
			args: args{err: NewWithCode(codes.CodeAuthFailure, "auth failed"), lang: language.English},
			want: http.StatusUnauthorized,
			want1: App{
				Code:  codes.CodeAuthFailure,
				Title: "Unauthorized",
				Body:  "Unauthorized access. You are not authorized to access this resource.",
				sys:   NewWithCode(codes.CodeAuthFailure, "auth failed"),
			},
		},
		{
			name: "not ok",
			args: args{err: NewWithCode(codes.NoCode, "no code"), lang: language.English},
			want: http.StatusInternalServerError,
			want1: App{
				Code:  codes.NoCode,
				Title: "Service Error Not Defined",
				Body:  "Unknown error. Please contact admin",
				sys:   NewWithCode(codes.CodeAuthFailure, "auth failed"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Compile(tt.args.err, tt.args.lang)
			if got != tt.want {
				t.Errorf("Compile() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1.Code, tt.want1.Code) {
				t.Errorf("Compile() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetCaller(t *testing.T) {
	pwd, _ := os.Getwd()
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   int
		want2   string
		wantErr bool
	}{
		{
			name:    "ok",
			args:    args{err: NewWithCode(codes.CodeBadRequest, "bad request")},
			want:    pwd + "/errors_test.go",
			want1:   98,
			want2:   "bad request",
			wantErr: false,
		},
		{
			name:    "not ok",
			args:    args{err: fmt.Errorf("")},
			want:    "",
			want1:   0,
			want2:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := GetCaller(tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCaller() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCaller() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetCaller() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("GetCaller() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
