package query

import (
	"reflect"
	"testing"

	"github.com/downsized-devs/sdk-go/null"
)

func Test_convertFloatArgs(t *testing.T) {

	getFloat64Ptr := func(x float64) *float64 {
		return &x
	}

	getFloat32Ptr := func(x float32) *float32 {
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
			name: "[]float64",
			args: args{
				_f: []float64{1, 2, 3},
			},
			wantPrimitiveType: int8(25),
			wantIsMany:        true,
			wantArgs:          []float64{1, 2, 3},
		},
		{
			name: "[]*float64",
			args: args{
				_f: []*float64{getFloat64Ptr(1), getFloat64Ptr(2), getFloat64Ptr(3)},
			},
			wantPrimitiveType: int8(25),
			wantIsMany:        true,
			wantArgs:          []*float64{getFloat64Ptr(1), getFloat64Ptr(2), getFloat64Ptr(3)},
		},
		{
			name: "float64",
			args: args{
				_f: float64(1),
			},
			wantPrimitiveType: int8(24),
			wantIsMany:        false,
			wantArgs:          float64(1),
		},
		{
			name: "[]float32",
			args: args{
				_f: []float32{1, 2, 3},
			},
			wantPrimitiveType: int8(23),
			wantIsMany:        true,
			wantArgs:          []float32{1, 2, 3},
		},
		{
			name: "[]*float32",
			args: args{
				_f: []*float32{getFloat32Ptr(1), getFloat32Ptr(2), getFloat32Ptr(3)},
			},
			wantPrimitiveType: int8(23),
			wantIsMany:        true,
			wantArgs:          []*float32{getFloat32Ptr(1), getFloat32Ptr(2), getFloat32Ptr(3)},
		},
		{
			name: "float32",
			args: args{
				_f: float32(1),
			},
			wantPrimitiveType: int8(22),
			wantIsMany:        false,
			wantArgs:          float32(1),
		},
		{
			name: "[]null.Float64",
			args: args{
				_f: []null.Float64{
					{Float64: 1, Valid: true},
					{Float64: 2, Valid: true},
				},
			},
			wantPrimitiveType: Float64Arr,
			wantIsMany:        true,
			wantArgs:          []float64{1, 2},
		},
		{
			name: "[]*null.Float64",
			args: args{
				_f: []*null.Float64{
					{Float64: *getFloat64Ptr(1), Valid: true},
					{Float64: *getFloat64Ptr(2), Valid: true},
				},
			},
			wantPrimitiveType: Float64Arr,
			wantIsMany:        true,
			wantArgs:          []float64{1, 2},
		},
		{
			name: "null.Float64",
			args: args{
				_f: null.Float64{
					Float64: 1, Valid: true,
				},
			},
			wantPrimitiveType: Float64,
			wantIsMany:        false,
			wantArgs:          float64(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrimitiveType, gotIsMany, gotArgs := convertFloatArgs(tt.args._f)
			if gotPrimitiveType != tt.wantPrimitiveType {
				t.Errorf("convertFloatArgs() gotPrimitiveType = %v, want %v", gotPrimitiveType, tt.wantPrimitiveType)
			}
			if gotIsMany != tt.wantIsMany {
				t.Errorf("convertFloatArgs() gotIsMany = %v, want %v", gotIsMany, tt.wantIsMany)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("convertFloatArgs() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}
