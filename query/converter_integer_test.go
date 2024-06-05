package query

import (
	"reflect"
	"testing"

	"github.com/downsized-devs/sdk-go/null"
)

func Test_convertIntArgs(t *testing.T) {

	getInt64Ptr := func(x int64) *int64 {
		return &x
	}

	getUInt64Ptr := func(x uint64) *uint64 {
		return &x
	}

	getInt32Ptr := func(x int32) *int32 {
		return &x
	}

	getUInt32Ptr := func(x uint32) *uint32 {
		return &x
	}

	getInt16Ptr := func(x int16) *int16 {
		return &x
	}

	getUInt16Ptr := func(x uint16) *uint16 {
		return &x
	}

	getInt8Ptr := func(x int8) *int8 {
		return &x
	}

	getUInt8Ptr := func(x uint8) *uint8 {
		return &x
	}

	getIntPtr := func(x int) *int {
		return &x
	}

	getUIntPtr := func(x uint) *uint {
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
		wantIsSqlNull     bool
		wantArgs          interface{}
	}{
		{
			name: "[]int64",
			args: args{
				_f: []int64{1, 2},
			},
			wantPrimitiveType: Int64Arr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []int64{1, 2},
		},
		{
			name: "[]*int64",
			args: args{
				_f: []*int64{getInt64Ptr(1), getInt64Ptr(2)},
			},
			wantPrimitiveType: Int64Arr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []*int64{getInt64Ptr(1), getInt64Ptr(2)},
		},
		{
			name: "int64",
			args: args{
				_f: int64(1),
			},
			wantPrimitiveType: Int64,
			wantIsMany:        false,
			wantIsSqlNull:     false,
			wantArgs:          int64(1),
		},
		{
			name: "[]uint64",
			args: args{
				_f: []uint64{1, 2},
			},
			wantPrimitiveType: Uint64Arr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []uint64{1, 2},
		},
		{
			name: "[]*uint64",
			args: args{
				_f: []*uint64{getUInt64Ptr(1), getUInt64Ptr(2)},
			},
			wantPrimitiveType: Uint64Arr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []*uint64{getUInt64Ptr(1), getUInt64Ptr(2)},
		},
		{
			name: "uint64",
			args: args{
				_f: uint64(1),
			},
			wantPrimitiveType: Uint64,
			wantIsMany:        false,
			wantIsSqlNull:     false,
			wantArgs:          uint64(1),
		},
		{
			name: "[]null.int64",
			args: args{
				_f: []null.Int64{
					{
						Int64: 1,
						Valid: true,
					},
					{
						Int64: 2,
						Valid: true,
					},
				},
			},
			wantPrimitiveType: Int64Arr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []int64{1, 2},
		},
		{
			name: "[]*null.int64",
			args: args{
				_f: []*null.Int64{
					{
						Int64: *getInt64Ptr(1),
						Valid: true,
					},
					{
						Int64: *getInt64Ptr(2),
						Valid: true,
					},
				},
			},
			wantPrimitiveType: Int64Arr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []int64{1, 2},
		},
		{
			name: "null.int64",
			args: args{
				_f: null.Int64{
					Int64: 1,
					Valid: true,
				},
			},
			wantPrimitiveType: Int64,
			wantIsMany:        false,
			wantIsSqlNull:     false,
			wantArgs:          int64(1),
		},
		{
			name: "[]int32",
			args: args{
				_f: []int32{1, 2},
			},
			wantPrimitiveType: Int32Arr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []int32{1, 2},
		},
		{
			name: "[]*int32",
			args: args{
				_f: []*int32{getInt32Ptr(1), getInt32Ptr(2)},
			},
			wantPrimitiveType: Int32Arr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []*int32{getInt32Ptr(1), getInt32Ptr(2)},
		},
		{
			name: "int32",
			args: args{
				_f: int32(1),
			},
			wantPrimitiveType: Int32,
			wantIsMany:        false,
			wantIsSqlNull:     false,
			wantArgs:          int32(1),
		},
		{
			name: "[]uint32",
			args: args{
				_f: []uint32{1, 2},
			},
			wantPrimitiveType: Uint32Arr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []uint32{1, 2},
		},
		{
			name: "[]*uint32",
			args: args{
				_f: []*uint32{getUInt32Ptr(1), getUInt32Ptr(2)},
			},
			wantPrimitiveType: Uint32Arr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []*uint32{getUInt32Ptr(1), getUInt32Ptr(2)},
		},
		{
			name: "uint32",
			args: args{
				_f: uint32(1),
			},
			wantPrimitiveType: Uint32,
			wantIsMany:        false,
			wantIsSqlNull:     false,
			wantArgs:          uint32(1),
		},
		{
			name: "[]int16",
			args: args{
				_f: []int16{1, 2},
			},
			wantPrimitiveType: Int16Arr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []int16{1, 2},
		},
		{
			name: "[]*int16",
			args: args{
				_f: []*int16{getInt16Ptr(1), getInt16Ptr(2)},
			},
			wantPrimitiveType: Int16Arr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []*int16{getInt16Ptr(1), getInt16Ptr(2)},
		},
		{
			name: "int16",
			args: args{
				_f: int16(1),
			},
			wantPrimitiveType: Int16,
			wantIsMany:        false,
			wantIsSqlNull:     false,
			wantArgs:          int16(1),
		},
		{
			name: "[]uint16",
			args: args{
				_f: []uint16{1, 2},
			},
			wantPrimitiveType: Uint16Arr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []uint16{1, 2},
		},
		{
			name: "[]*uint16",
			args: args{
				_f: []*uint16{getUInt16Ptr(1), getUInt16Ptr(2)},
			},
			wantPrimitiveType: Uint16Arr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []*uint16{getUInt16Ptr(1), getUInt16Ptr(2)},
		},
		{
			name: "uint16",
			args: args{
				_f: uint16(1),
			},
			wantPrimitiveType: Uint16,
			wantIsMany:        false,
			wantIsSqlNull:     false,
			wantArgs:          uint16(1),
		},
		{
			name: "[]int8",
			args: args{
				_f: []int8{1, 2},
			},
			wantPrimitiveType: Int8Arr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []int8{1, 2},
		},
		{
			name: "[]*int8",
			args: args{
				_f: []*int8{getInt8Ptr(1), getInt8Ptr(2)},
			},
			wantPrimitiveType: Int8Arr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []*int8{getInt8Ptr(1), getInt8Ptr(2)},
		},
		{
			name: "int8",
			args: args{
				_f: int8(1),
			},
			wantPrimitiveType: Int8,
			wantIsMany:        false,
			wantIsSqlNull:     false,
			wantArgs:          int8(1),
		},
		{
			name: "[]uint8",
			args: args{
				_f: []uint8{1, 2},
			},
			wantPrimitiveType: Uint8Arr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []uint8{1, 2},
		},
		{
			name: "[]*uint8",
			args: args{
				_f: []*uint8{getUInt8Ptr(1), getUInt8Ptr(2)},
			},
			wantPrimitiveType: Uint8Arr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []*uint8{getUInt8Ptr(1), getUInt8Ptr(2)},
		},
		{
			name: "uint8",
			args: args{
				_f: uint8(1),
			},
			wantPrimitiveType: Uint8,
			wantIsMany:        false,
			wantIsSqlNull:     false,
			wantArgs:          uint8(1),
		},
		{
			name: "[]int",
			args: args{
				_f: []int{1, 2},
			},
			wantPrimitiveType: IntArr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []int{1, 2},
		},
		{
			name: "[]*int",
			args: args{
				_f: []*int{getIntPtr(1), getIntPtr(2)},
			},
			wantPrimitiveType: IntArr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []*int{getIntPtr(1), getIntPtr(2)},
		},
		{
			name: "int",
			args: args{
				_f: int(1),
			},
			wantPrimitiveType: Int,
			wantIsMany:        false,
			wantIsSqlNull:     false,
			wantArgs:          int(1),
		},
		{
			name: "[]uint",
			args: args{
				_f: []uint{1, 2},
			},
			wantPrimitiveType: UintArr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []uint{1, 2},
		},
		{
			name: "[]*uint",
			args: args{
				_f: []*uint{getUIntPtr(1), getUIntPtr(2)},
			},
			wantPrimitiveType: UintArr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []*uint{getUIntPtr(1), getUIntPtr(2)},
		},
		{
			name: "uint",
			args: args{
				_f: uint(1),
			},
			wantPrimitiveType: Uint,
			wantIsMany:        false,
			wantIsSqlNull:     false,
			wantArgs:          uint(1),
		},
		{
			name: "null.int64: SET to NULL",
			args: args{
				_f: null.Int64{
					SqlNull: true,
				},
			},
			wantPrimitiveType: Int64,
			wantIsMany:        false,
			wantIsSqlNull:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrimitiveType, gotIsMany, gotIsSqlNull, gotArgs := convertIntArgs(tt.args._f)
			if gotPrimitiveType != tt.wantPrimitiveType {
				t.Errorf("convertIntArgs() gotPrimitiveType = %v, want %v", gotPrimitiveType, tt.wantPrimitiveType)
			}
			if gotIsMany != tt.wantIsMany {
				t.Errorf("convertIntArgs() gotIsMany = %v, want %v", gotIsMany, tt.wantIsMany)
			}
			if gotIsSqlNull != tt.wantIsSqlNull {
				t.Errorf("convertTimeArgs() gotIsSqlNull = %v, want %v", gotIsSqlNull, tt.wantIsSqlNull)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("convertIntArgs() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}
