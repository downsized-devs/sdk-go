package num

import "testing"

func TestRoundingNumber(t *testing.T) {
	type args struct {
		num float64
		n   int
	}

	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "rounded decimal number up to 1 decimal place",
			args: args{
				num: 1.546789120718,
				n:   1,
			},
			want: 1.5,
		},
		{
			name: "rounded decimal number up to 3 decimal place",
			args: args{
				num: 1.546789120718,
				n:   3,
			},
			want: 1.547,
		},
		{
			name: "rounded decimal number up to 5 decimal place",
			args: args{
				num: 1.546789120718,
				n:   5,
			},
			want: 1.54679,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RoundingNumber(tt.args.num, tt.args.n)
			if got != tt.want {
				t.Errorf("RandomString() = %v, want %v", got, tt.want)
			}
		})
	}
}
