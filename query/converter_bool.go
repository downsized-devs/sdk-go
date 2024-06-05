package query

import (
	"github.com/downsized-devs/sdk-go/null"
)

func convertBoolArgs(_f interface{}) (primitiveType int8, isMany bool, args interface{}) {
	switch f := _f.(type) {
	case []bool:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Bool

	case []*bool:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = Bool

	case bool:
		if f {
			args = f
		}
		primitiveType = Bool

	case []null.Bool:
		if len(f) > 0 {
			var _args []bool
			for _, r := range f {
				if r.Valid {
					_args = append(_args, r.Bool)
				}
			}
			isMany = true
			args = _args
		}
		primitiveType = Bool

	case []*null.Bool:
		if len(f) > 0 {
			var _args []bool
			for _, r := range f {
				if r != nil {
					if r.Valid {
						_args = append(_args, r.Bool)
					}
				}
			}
			isMany = true
			args = _args
		}
		primitiveType = Bool

	case null.Bool:
		if f.Valid {
			args = f.Bool
		}
		primitiveType = Bool
	}

	return primitiveType, isMany, args
}
