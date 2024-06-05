package files

import (
	"testing"
)

func TestGetExtension(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "filename with ext len < 1",
			args: args{
				filename: "",
			},
			want: "",
		},
		{
			name: "file extension = file name",
			args: args{
				filename: "akagami.shanks",
			},
			want: "shanks",
		},
		{
			name: "file extension = file name",
			args: args{
				filename: "yonko.onepiece.admiral",
			},
			want: "admiral",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetExtension(tt.args.filename); got != tt.want {
				t.Errorf("GetExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsExist(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "file name does not exists",
			args: args{
				filename: "",
			},
			want: false,
		},
		{
			name: "file name exists",
			args: args{filename: "test_folder/test_file"},
			want: true,
		},
		{
			name: "file name exists",
			args: args{filename: "test_folder"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsExist(tt.args.filename); got != tt.want {
				t.Errorf("IsExist() = %v, want %v", got, tt.want)
			}
		})
	}
}
