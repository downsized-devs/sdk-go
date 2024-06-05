package query

import (
	"testing"
)

func Test_getOffset(t *testing.T) {
	type args struct {
		p int64
		l int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "p = 0",
			args: args{p: 1, l: 1},
			want: 0,
		},
		{
			name: "p <= 0",
			args: args{p: 0, l: 1},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getOffset(tt.args.p, tt.args.l); got != tt.want {
				t.Errorf("getOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isPage(t *testing.T) {
	type args struct {
		paramTagValue string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ok",
			args: args{paramTagValue: "page"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPage(tt.args.paramTagValue); got != tt.want {
				t.Errorf("isPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isLimit(t *testing.T) {
	type args struct {
		paramTagValue string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ok",
			args: args{paramTagValue: "limit"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLimit(tt.args.paramTagValue); got != tt.want {
				t.Errorf("isLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isSortBy(t *testing.T) {
	type args struct {
		paramTagValue string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "error",
			args: args{paramTagValue: ""},
			want: false,
		},
		{
			name: "ok",
			args: args{paramTagValue: "sortby"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSortBy(tt.args.paramTagValue); got != tt.want {
				t.Errorf("isSortBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateLimit(t *testing.T) {
	type args struct {
		l int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "ok",
			args: args{l: 0},
			want: 10,
		},
		{
			name: "ok",
			args: args{l: 1},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateLimit(tt.args.l); got != tt.want {
				t.Errorf("validateLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validatePage(t *testing.T) {
	type args struct {
		p int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validatePage(tt.args.p); got != tt.want {
				t.Errorf("validatePage() = %v, want %v", got, tt.want)
			}
		})
	}
}
