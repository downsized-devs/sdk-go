package query

import (
	"time"

	"github.com/downsized-devs/sdk-go/null"
)

func convertTimeArgs(_f interface{}) (primitiveType int8, isMany, isSqlNull bool, args interface{}) {
	switch f := _f.(type) {
	case time.Time:
		if !f.IsZero() {
			args = f
		}
		primitiveType = Time

	case []time.Time:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = TimeArr

	case []*time.Time:
		if len(f) > 0 {
			isMany = true
			args = f
		}
		primitiveType = TimeArr

	case []null.Time:
		if len(f) > 0 {
			var _args []time.Time
			for _, r := range f {
				if r.Valid {
					_args = append(_args, r.Time)
				}
			}
			isMany = true
			args = _args
		}
		primitiveType = TimeArr

	case []*null.Time:
		if len(f) > 0 {
			var _args []time.Time
			for _, r := range f {
				if r != nil {
					if r.Valid {
						_args = append(_args, r.Time)
					}
				}
			}
			isMany = true
			args = _args
		}
		primitiveType = TimeArr

	case []null.Date:
		if len(f) > 0 {
			var _args []time.Time
			for _, r := range f {
				if r.Valid {
					_args = append(_args, r.Time)
				}
			}
			isMany = true
			args = _args
		}
		primitiveType = TimeArr

	case []*null.Date:
		if len(f) > 0 {
			var _args []time.Time
			for _, r := range f {
				if r != nil {
					if r.Valid {
						_args = append(_args, r.Time)
					}
				}
			}
			isMany = true
			args = _args
		}
		primitiveType = TimeArr

	case null.Time:
		if f.SqlNull {
			isSqlNull = true
		}
		if f.Valid {
			args = f.Time
		}
		primitiveType = Time

	case null.Date:
		if f.SqlNull {
			isSqlNull = true
		}
		if f.Valid {
			args = f.Time
		}
		primitiveType = Time
	}

	return primitiveType, isMany, isSqlNull, args
}
