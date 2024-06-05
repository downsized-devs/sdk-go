package query

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatQueryForRows(t *testing.T) {
	type args struct {
		ctx    context.Context
		q      string
		inputs [][]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []interface{}
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
				q:   "test",
				inputs: [][]interface{}{
					{"test"}, {"test"},
				},
			},
			want:    "test (?), (?)",
			want1:   []interface{}{"test", "test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := FormatQueryForRows(tt.args.ctx, tt.args.q, tt.args.inputs)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatQueryForRows() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FormatQueryForRows() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("FormatQueryForRows() got1 = %v, want %v", got1, tt.want1)
			}
			assert.Equal(t, tt.want, got)
			// assert.Equal(t, tt.want1, got1...)
		})
	}
}
