package num

import (
	"reflect"
	"testing"
)

func TestEmptyStringSlice(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty string slice generated",
			args: args{
				length: 2,
			},
			want: []string{"", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EmptyStringSlice(tt.args.length); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EmptyStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
