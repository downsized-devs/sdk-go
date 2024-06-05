//go:build integration
// +build integration

package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceBindvarsWithArgs(t *testing.T) {
	type args struct {
		str  string
		args interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				str:  "SELECT * FROM farm WHERE id = ?",
				args: "1",
			},
			want: "SELECT * FROM farm WHERE id = 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := replaceBindvarsWithArgs(tt.args.str, tt.args.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
