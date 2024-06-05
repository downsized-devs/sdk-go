package query

import (
	"testing"

	"github.com/downsized-devs/sdk-go/null"
	"github.com/stretchr/testify/assert"
)

func Test_convertBoolArgs(t *testing.T) {
	getBoolPtr := func(x bool) *bool {
		return &x
	}
	type args struct {
		_f interface{}
	}
	tests := []struct {
		name              string
		args              args
		wantPrimitiveType int8
		wantIsMany        bool
		wantArgs          interface{}
	}{
		{
			name: "bool",
			args: args{
				_f: bool(true),
			},
			wantPrimitiveType: int8(26),
			wantIsMany:        bool(false),
			wantArgs:          bool(true),
		},
		{
			name: "[]bool",
			args: args{
				_f: []bool{true, false},
			},
			wantPrimitiveType: int8(26),
			wantIsMany:        bool(true),
			wantArgs:          []bool{true, false},
		},
		{
			name: "[]null.bool",
			args: args{
				_f: []null.Bool{
					{Bool: true, Valid: true},
					{Bool: false, Valid: true},
				},
			},
			wantPrimitiveType: int8(26),
			wantIsMany:        bool(true),
			wantArgs:          []bool{true, false},
		},
		{
			name: "null.bool",
			args: args{
				_f: null.Bool{
					Bool:  true,
					Valid: true},
			},
			wantPrimitiveType: int8(26),
			wantIsMany:        bool(false),
			wantArgs:          bool(true),
		},
		{
			name: "[]*bool",
			args: args{
				_f: []*bool{getBoolPtr(true), getBoolPtr(false)},
			},
			wantPrimitiveType: int8(26),
			wantIsMany:        bool(true),
			wantArgs:          []*bool{getBoolPtr(true), getBoolPtr(false)},
		},
		{
			name: "[]*null.bool",
			args: args{
				_f: []*null.Bool{
					{Bool: *getBoolPtr(true), Valid: true},
					{Bool: *getBoolPtr(false), Valid: true},
				},
			},
			wantPrimitiveType: int8(26),
			wantIsMany:        bool(true),
			wantArgs:          []bool{true, false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrimitiveType, gotIsMany, gotArgs := convertBoolArgs(tt.args._f)
			if gotPrimitiveType != tt.wantPrimitiveType {
				t.Errorf("convertBoolArgs() gotPrimitiveType = %v, want %v", gotPrimitiveType, tt.wantPrimitiveType)
			}
			if gotIsMany != tt.wantIsMany {
				t.Errorf("convertBoolArgs() gotIsMany = %v, want %v", gotIsMany, tt.wantIsMany)
			}
			assert.Equal(t, tt.wantArgs, gotArgs)
		})
	}
}
