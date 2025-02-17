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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Compile(tt.args.code, tt.args.lang); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Compile() = %v, want %v", got, tt.want)
			}
		})
	}
}
