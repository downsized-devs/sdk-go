package operator

import (
	"reflect"
	"testing"
)

func TestTernaryInt(t *testing.T) {
	type args struct {
		condition bool
		a         int
		b         int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "int A",
			args: args{
				condition: true,
				a:         1,
				b:         2,
			},
			want: 1,
		},
		{
			name: "int B",
			args: args{
				condition: false,
				a:         3,
				b:         4,
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ternary(tt.args.condition, tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ternary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTernaryFloat(t *testing.T) {
	type args struct {
		condition bool
		a         float64
		b         float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "float A",
			args: args{
				condition: true,
				a:         1,
				b:         2,
			},
			want: 1,
		},
		{
			name: "float B",
			args: args{
				condition: false,
				a:         3,
				b:         4,
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ternary(tt.args.condition, tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ternary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTernaryString(t *testing.T) {
	type args struct {
		condition bool
		a         string
		b         string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "string A",
			args: args{
				condition: true,
				a:         "yonko",
				b:         "admiral",
			},
			want: "yonko",
		},
		{
			name: "string B",
			args: args{
				condition: false,
				a:         "yonko",
				b:         "admiral",
			},
			want: "admiral",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ternary(tt.args.condition, tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ternary() = %v, want %v", got, tt.want)
			}
		})
	}
}
