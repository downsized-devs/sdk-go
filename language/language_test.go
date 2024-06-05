package language

import (
	"net/http"
	"testing"
)

func TestHTTPStatusText(t *testing.T) {
	type args struct {
		lang string
		code int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "en",
			args: args{lang: English, code: http.StatusContinue},
			want: "Continue",
		},
		{
			name: "id",
			args: args{lang: Indonesian, code: http.StatusAccepted},
			want: "Diterima",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HTTPStatusText(tt.args.lang, tt.args.code); got != tt.want {
				t.Errorf("HTTPStatusText() = %v, want %v", got, tt.want)
			}
		})
	}
}
