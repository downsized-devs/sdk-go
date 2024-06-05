package helper

import (
	"testing"
)

func TestConvertToUpperCase(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "all good",
			args: args{
				text: "testText",
			},
			want: "TestText",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertToUpperCase(tt.args.text); got != tt.want {
				t.Errorf("ConvertToUpperCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToSnakeCase(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "all good",
			args: args{
				text: "testText",
			},
			want: "test_text",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertToSnakeCase(tt.args.text); got != tt.want {
				t.Errorf("ConvertToSnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertToSpace(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "all good",
			args: args{
				text: "testText",
			},
			want: "test Text",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertToSpace(tt.args.text); got != tt.want {
				t.Errorf("convertToSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToUpperSpace(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "all good",
			args: args{
				text: "testText",
			},
			want: "Test Text",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertToUpperSpace(tt.args.text); got != tt.want {
				t.Errorf("ConvertToUpperSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToLowerSpace(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "all good",
			args: args{
				text: "testText",
			},
			want: "test text",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertToLowerSpace(tt.args.text); got != tt.want {
				t.Errorf("ConvertToLowerSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToCamelCase(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "all good",
			args: args{
				text: "TestText",
			},
			want: "testText",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertToCamelCase(tt.args.text); got != tt.want {
				t.Errorf("ConvertToCamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToLowerDash(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "all good",
			args: args{
				text: "testText",
			},
			want: "test-text",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertToLowerDash(tt.args.text); got != tt.want {
				t.Errorf("ConvertToLowerDash() = %v, want %v", got, tt.want)
			}
		})
	}
}
