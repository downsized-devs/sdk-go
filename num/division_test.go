package num

import "testing"

func TestSafeDivision(t *testing.T) {
	type args struct {
		numerator   float64
		denominator float64
		zeroValue   bool
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "denominator = 0 returns 0",
			args: args{
				numerator:   10,
				denominator: 0,
				zeroValue:   true,
			},
			want: 0,
		},
		{
			name: "denominator = 0 returns original value",
			args: args{
				numerator:   10,
				denominator: 0,
				zeroValue:   false,
			},
			want: 10,
		},
		{
			name: "denominator != 0",
			args: args{
				numerator:   10,
				denominator: 100,
				zeroValue:   true,
			},
			want: 0.1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SafeDivision(tt.args.numerator, tt.args.denominator, tt.args.zeroValue); got != tt.want {
				t.Errorf("DivisionChecker() = %v, want %v", got, tt.want)
			}
		})
	}
}
