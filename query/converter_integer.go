package query

import (
	"github.com/downsized-devs/sdk-go/null"
)

const (
	MaxUint64 = ^uint64(0)
	MinUint64 = uint64(0)
	MaxUint32 = ^uint32(0)
	MinUint32 = uint32(0)
	MaxUint16 = ^uint16(0)
	MinUint16 = uint16(0)
	MaxUint8  = ^uint8(0)
	MinUint8  = uint8(0)
	MaxUint   = ^uint(0)
	MinUint   = uint(0)
	MaxInt64  = int64(MaxUint64 >> 1)
	MinInt64  = -MaxInt64 - 1
	MaxInt32  = int32(MaxUint32 >> 1)
	MinInt32  = -MaxInt32 - 1
	MaxInt16  = int16(MaxUint16 >> 1)
	MinInt16  = -MaxInt16 - 1
	MaxInt8   = int8(MaxUint8 >> 1)
	MinInt8   = -MaxInt8 - 1
	MaxInt    = int(MaxUint >> 1)
	MinInt    = -MaxInt - 1
)

func convertIntArgs(_f interface{}) (primitiveType int8, isMany bool, isSqlNull bool, args interface{}) {
	switch f := _f.(type) {
	case []int64:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Int64Arr
	case []*int64:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Int64Arr

	case int64:
		if f > 0 {
			args = f
		}
		primitiveType = Int64

	case []uint64:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Uint64Arr

	case []*uint64:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Uint64Arr

	case uint64:
		if f > 0 {
			args = f
		}
		primitiveType = Uint64

	case []null.Int64:
		if len(f) > 0 {
			var _args []int64
			for _, r := range f {
				if r.Valid && r.Int64 >= MinInt64 {
					_args = append(_args, r.Int64)
				}
			}
			isMany = true
			args = _args
		}
		primitiveType = Int64Arr

	case []*null.Int64:
		if len(f) > 0 {
			var _args []int64
			for _, r := range f {
				if r != nil {
					if r.Valid && r.Int64 >= MinInt64 {
						_args = append(_args, r.Int64)
					}
				}
			}
			isMany = true
			args = _args
		}
		primitiveType = Int64Arr

	case null.Int64:
		if f.SqlNull {
			isSqlNull = true
		}
		if f.Valid {
			args = f.Int64
		}
		primitiveType = Int64

	case []int32:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Int32Arr

	case []*int32:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Int32Arr

	case int32:
		if f > 0 {
			args = f
		}
		primitiveType = Int32

	case []uint32:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Uint32Arr

	case []*uint32:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Uint32Arr

	case uint32:
		if f > 0 {
			args = f
		}
		primitiveType = Uint32

	case []int16:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Int16Arr

	case []*int16:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Int16Arr

	case int16:
		if f > 0 {
			args = f
		}
		primitiveType = Int16

	case []uint16:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Uint16Arr

	case []*uint16:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Uint16Arr

	case uint16:
		if f > 0 {
			args = f
		}
		primitiveType = Uint16

	case []int8:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Int8Arr

	case []*int8:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Int8Arr

	case int8:
		if f > 0 {
			args = f
		}
		primitiveType = Int8

	case []uint8:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Uint8Arr

	case []*uint8:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Uint8Arr

	case uint8:
		if f > 0 {
			args = f
		}
		primitiveType = Uint8

	case []int:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = IntArr

	case []*int:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = IntArr

	case int:
		if f > 0 {
			args = f
		}
		primitiveType = Int

	case []uint:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = UintArr

	case []*uint:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = UintArr

	case uint:
		if f > 0 {
			args = f
		}
		primitiveType = Uint
	}

	return primitiveType, isMany, isSqlNull, args
}
