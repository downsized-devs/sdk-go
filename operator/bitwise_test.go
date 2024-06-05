package operator

import "testing"

func TestCheckBitOnPosition(t *testing.T) {
	type args struct {
		number   int
		position int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			/*
				number = 75 and position = 4
				temp = 1 << (position-1) = 1 << 3 = 8
				Binary Representation of temp = 0..00001000
				Binary Representation of position = 0..01001011
				Since bitwise AND of number and temp is non-zero, result is SET.
			*/
			name: "success check bit 4th on number 75 is SET or NOT SET",
			args: args{
				number:   75,
				position: 4,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckBitOnPosition(tt.args.number, tt.args.position); got != tt.want {
				t.Errorf("CheckKthBitSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
