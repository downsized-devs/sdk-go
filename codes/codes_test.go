package codes

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/downsized-devs/sdk-go/language"
)

func TestCompile(t *testing.T) {
	type args struct {
		code Code
		lang string
	}
	tests := []struct {
		name string
		args args
		want DisplayMessage
	}{
		{
			name: "default success message in english",
			args: args{code: CodeSuccess, lang: language.English},
			want: DisplayMessage{StatusCode: http.StatusOK, Title: "OK", Body: "Request successful"},
		},
		{
			name: "default success message in indonesian",
			args: args{code: CodeSuccess, lang: language.Indonesian},
			want: DisplayMessage{StatusCode: http.StatusOK, Title: "OK", Body: "Request berhasil"},
		},
		{
			name: "code not in ApplicationMessages uses default",
			args: args{code: CodeBadRequest, lang: language.English},
			want: DisplayMessage{StatusCode: http.StatusOK, Title: "OK", Body: "Request successful"},
		},
		{
			name: "code in ApplicationMessages - english",
			args: args{code: CodeAccepted, lang: language.English},
			want: DisplayMessage{StatusCode: http.StatusAccepted, Title: "Accepted", Body: "Request accepted"},
		},
		{
			name: "code in ApplicationMessages - indonesian",
			args: args{code: CodeAccepted, lang: language.Indonesian},
			want: DisplayMessage{StatusCode: http.StatusAccepted, Title: "Diterima", Body: "Request diterima"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Compile(tt.args.code, tt.args.lang); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Compile() = %v, want %v", got, tt.want)
			}
		})
	}
}
