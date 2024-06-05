package character

import (
	"testing"

	"golang.org/x/text/language"
)

func TestCapitalizeFirstCharacter(t *testing.T) {
	type args struct {
		s   string
		lan language.Tag
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test with normal sentence in English",
			args: args{s: "hello world", lan: language.English},
			want: "Hello World",
		},
		{
			name: "Test with sentence in French",
			args: args{s: "this is a test", lan: language.French},
			want: "This Is A Test",
		},
		{
			name: "Test with empty input in Spanish",
			args: args{s: "", lan: language.Spanish},
			want: "",
		},
		{
			name: "Test with single word in Italian",
			args: args{s: "one", lan: language.Italian},
			want: "One",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CapitalizeFirstCharacter(tt.args.s, tt.args.lan); got != tt.want {
				t.Errorf("CapitalizeFirstCharacter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsStrongCharCombination(t *testing.T) {
	tests := []struct {
		name      string
		character string
		minLength int
		want      bool
	}{
		{
			name:      "Test with strong character combination",
			character: "Abcd123!",
			minLength: 8,
			want:      true,
		},
		{
			name:      "Test with strong character combination \\",
			character: "Abcd123\\",
			minLength: 8,
			want:      true,
		},
		{
			name:      "Test with strong character combination _",
			character: "Abcd123_",
			minLength: 8,
			want:      true,
		},
		{
			name:      "Test with backtic",
			character: "Abcd123`",
			minLength: 8,
			want:      true,
		},
		{
			name:      "Test with weak character combination (no uppercase)",
			character: "abcd123!",
			minLength: 8,
			want:      false,
		},
		{
			name:      "Test with weak character combination (short length)",
			character: "Abc1!",
			minLength: 8,
			want:      false,
		},
		{
			name:      "Test with weak character combination (no digit)",
			character: "Abcdefgh!",
			minLength: 8,
			want:      false,
		},
		{
			name:      "Test with weak character combination (no special character)",
			character: "Abcd1234",
			minLength: 8,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStrongCharCombination(tt.character, tt.minLength); got != tt.want {
				t.Errorf("IsStrongCharCombination() = %v, want %v", got, tt.want)
			}
		})
	}
}
