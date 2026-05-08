package errors

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"runtime"
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

	// Capture the error and its creation line at runtime so the test remains
	// correct regardless of future edits that shift line numbers in this file.
	sdkErr := NewWithCode(codes.CodeBadRequest, "bad request")
	_, _, sdkErrLine, _ := runtime.Caller(0)
	sdkErrLine-- // NewWithCode was called on the line above

	type args struct {
		err error
	}
	tests := []struct {
		name            string
		args            args
		want            string
		want1           int
		checkLineNumber bool
		want2           string
		wantErr         bool
	}{
		{
			name:            "ok",
			args:            args{err: sdkErr},
			want:            pwd + "/errors_test.go",
			want1:           sdkErrLine,
			checkLineNumber: true,
			want2:           "bad request",
			wantErr:         false,
		},
		{
			name:    "not ok",
			args:    args{err: fmt.Errorf("")},
			want:    "",
			want1:   0,
			want2:   "",
			wantErr: true,
		},
		{
			name:    "wrapped ok",
			args:    args{err: fmt.Errorf("context: %w", NewWithCode(codes.CodeBadRequest, "wrapped bad request"))},
			want:    pwd + "/errors_test.go",
			want2:   "wrapped bad request",
			wantErr: false,
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
			if tt.checkLineNumber && got1 != tt.want1 {
				t.Errorf("GetCaller() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("GetCaller() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestGetCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want codes.Code
	}{
		{
			name: "direct stacktrace",
			err:  NewWithCode(codes.CodeBadRequest, "bad request"),
			want: codes.CodeBadRequest,
		},
		{
			name: "wrapped stacktrace",
			err:  fmt.Errorf("context: %w", NewWithCode(codes.CodeBadRequest, "bad request")),
			want: codes.CodeBadRequest,
		},
		{
			name: "non-stacktrace error",
			err:  fmt.Errorf("plain error"),
			want: codes.NoCode,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCode(tt.err); got != tt.want {
				t.Errorf("GetCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_Unwrap(t *testing.T) {
	inner := NewWithCode(codes.CodeBadRequest, "inner error")
	app := &App{sys: inner}
	if got := app.Unwrap(); got != inner {
		t.Errorf("App.Unwrap() = %v, want %v", got, inner)
	}
}

func TestWrapWithCode(t *testing.T) {
	cause := fmt.Errorf("original error")
	wrapped := WrapWithCode(cause, codes.CodeBadRequest, "wrapped: %s", "msg")
	if GetCode(wrapped) != codes.CodeBadRequest {
		t.Errorf("WrapWithCode() code = %v, want %v", GetCode(wrapped), codes.CodeBadRequest)
	}
	// errors.Is must still find the cause
	if !Is(wrapped, cause) {
		t.Errorf("WrapWithCode() wrapped error should satisfy errors.Is for the cause")
	}
}

func TestIs(t *testing.T) {
	sentinel := fmt.Errorf("sentinel")
	wrapped := fmt.Errorf("wrap: %w", sentinel)
	if !Is(wrapped, sentinel) {
		t.Errorf("Is() should find sentinel in wrapped chain")
	}
	if Is(wrapped, fmt.Errorf("other")) {
		t.Errorf("Is() should not match unrelated error")
	}
}

func TestAs(t *testing.T) {
	inner := NewWithCode(codes.CodeBadRequest, "bad request")
	wrapped := fmt.Errorf("wrap: %w", inner)
	var st *stacktrace
	if !As(wrapped, &st) {
		t.Errorf("As() should find *stacktrace in wrapped chain")
	}
	if st == nil {
		t.Errorf("As() should set target to the *stacktrace value")
	}
}

func TestStacktrace_Unwrap(t *testing.T) {
	cause := fmt.Errorf("root cause")
	err := WrapWithCode(cause, codes.CodeBadRequest, "outer")
	// Unwrap through the stacktrace should expose the cause
	if !Is(err, cause) {
		t.Errorf("stacktrace.Unwrap() should expose cause via errors.Is")
	}
}

func TestCompile_Indonesian(t *testing.T) {
	got, _ := Compile(NewWithCode(codes.CodeAuthFailure, "auth failed"), "id")
	if got != 401 {
		t.Errorf("Compile(Indonesian) got status %d, want 401", got)
	}
}

func TestNewWithCode_Formatting(t *testing.T) {
	err := NewWithCode(codes.CodeBadRequest, "value is %d", 42)
	if err.Error() != "value is 42" {
		t.Errorf("NewWithCode() message = %q, want %q", err.Error(), "value is 42")
	}
}
