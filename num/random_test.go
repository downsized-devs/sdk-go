package num

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	mockTime := time.Date(2022, 5, 11, 0, 0, 0, 0, time.UTC)
	now = func() time.Time {
		return mockTime
	}

	defer func() {
		now = time.Now
	}()

	type args struct {
		n int
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "generated random string with 0 as len",
			args: args{
				n: 0,
			},
			want: "",
		},
		{
			name: "generated random string with 5 as len",
			args: args{
				n: 5,
			},
			want: "Sjmcz",
		},
		{
			name: "generated random string with 10 as len",
			args: args{
				n: 10,
			},
			want: "SjmczPXFsY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandomString(tt.args.n)
			if got != tt.want {
				t.Errorf("RandomString() = %v, want %v", got, tt.want)
			}
			assert.Equal(t, tt.args.n, len([]rune(got)))
		})
	}
}
