package dates

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDifference(t *testing.T) {
	mockTime := time.Now()
	mockTimeTommorow := time.Now().Add(24 * time.Hour)
	mockTime2 := time.Now()

	type args struct {
		a time.Time
		b time.Time
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "a after b",
			args: args{a: mockTime, b: mockTimeTommorow},
			want: 1,
		},
		{
			name: "b after a",
			args: args{a: mockTimeTommorow, b: mockTime},
			want: 1,
		},
		{
			name: "time is the same",
			args: args{a: mockTime, b: mockTime2},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Difference(tt.args.a, tt.args.b); got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
