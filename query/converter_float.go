package query

import (
	"math"

	"github.com/downsized-devs/sdk-go/null"
)

const (
	MaxFloat64 float64 = math.MaxFloat64
	MaxFloat32 float32 = math.MaxFloat32
	MinFloat64 float64 = math.MaxFloat64 * -1
	MinFloat32 float32 = math.MaxFloat32 * -1
)

func convertFloatArgs(_f interface{}) (primitiveType int8, isMany bool, args interface{}) {
	switch f := _f.(type) {
	case []float64:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Float64Arr

	case []*float64:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Float64Arr

	case float64:
		if f > 0 {
			args = f
		}
		primitiveType = Float64

	case []float32:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Float32Arr

	case []*float32:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Float32Arr

	case float32:
		if f > 0 {
			args = f
		}
		primitiveType = Float32

	case []null.Float64:
		if len(f) > 0 {
			var _args []float64
			for _, v := range f {
				if v.Valid && v.Float64 >= MinFloat64 {
					_args = append(_args, v.Float64)
				}
			}
			isMany = true
			args = _args
		}
		primitiveType = Float64Arr

	case []*null.Float64:
		if len(f) > 0 {
			var _args []float64
			for _, v := range f {
				if v != nil {
					if v.Valid && v.Float64 >= MinFloat64 {
						_args = append(_args, v.Float64)
					}
				}
			}
			isMany = true
			args = _args
		}
		primitiveType = Float64Arr

	case null.Float64:
		if f.Valid {
			args = f.Float64
		}
		primitiveType = Float64
	}

	return primitiveType, isMany, args
}
