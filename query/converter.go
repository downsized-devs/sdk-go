package query

import (
	"fmt"
	"reflect"
	"time"

	"github.com/downsized-devs/sdk-go/null"
)

type builderFunction func(primitiveType int8, isLike, isMany, isSqlNull bool, fieldName, paramTag, dbTag string, args interface{})

const (
	Int int8 = iota
	IntArr
	Int64
	Int64Arr
	Int32
	Int32Arr
	Int16
	Int16Arr
	Int8
	Int8Arr
	Uint
	UintArr
	Uint64
	Uint64Arr
	Uint32
	Uint32Arr
	Uint16
	Uint16Arr
	Uint8
	Uint8Arr
	String
	StringArr
	Float32
	Float32Arr
	Float64
	Float64Arr
	Bool
	Time
	TimeArr
)

func traverseOnParam(paramTagName, dbTagName, fieldTagName, fieldName, paramTagValue, dbTagValue string, aliasMap map[string]string, p reflect.Value, builderFunc builderFunction) {
	switch p.Kind() {

	// on pointer/ interface
	case reflect.Ptr, reflect.Interface:

		if !p.Elem().IsValid() || p.IsNil() {
			return
		}

		// handle if is not time type, null type, and struct
		// continue to traverse
		if !isTimeType(p.Elem()) && !isNullType(p.Elem()) && p.Elem().Kind() == reflect.Struct {
			traverseOnParam(paramTagName, dbTagName, fieldTagName, fieldName+"."+p.Elem().Type().Name(), paramTagValue, dbTagValue, aliasMap, p.Elem(), builderFunc)
		}

		// else convert on types
		convertOnTypes(paramTagValue, dbTagValue, fieldName, p.Elem(), builderFunc)
		return

	// on struct
	case reflect.Struct:
		if isTimeType(p) {
			convertOnTypes(paramTagValue, dbTagValue, fieldName+"."+getNameFromStructTagOrOriginalName(fieldTagName, p, 0), p, builderFunc)
			return
		}

		for i := 0; i < p.NumField(); i++ {
			// only exported struct that can be Traversed
			if p.Field(i).CanSet() {
				paramTagValue = p.Type().Field(i).Tag.Get(paramTagName)
				dbTagValue = p.Type().Field(i).Tag.Get(dbTagName)

				if dbTagValue == "-" {
					continue
				}
				var address string
				if p.CanAddr() {
					address = fmt.Sprint(p.Addr().Pointer())
				}
				alias := aliasMap[address]
				if alias != "" && dbTagValue != "" && address != "" {
					dbTagValue = alias + "." + dbTagValue
				}

				if isNullType(p.Field(i)) {
					convertOnTypes(paramTagValue, dbTagValue, fieldName+"."+getNameFromStructTagOrOriginalName(fieldTagName, p, i), p.Field(i), builderFunc)
					continue
				}

				traverseOnParam(paramTagName, dbTagName, fieldTagName, fieldName+"."+getNameFromStructTagOrOriginalName(fieldTagName, p, i), paramTagValue, dbTagValue, aliasMap, p.Field(i), builderFunc)
			}
		}

	default:
		convertOnTypes(paramTagValue, dbTagValue, fieldName, p, builderFunc)
		return
	}
}

func convertOnTypes(paramTagValue, dbTagValue, fieldName string, e reflect.Value, builderFunc builderFunction) {
	var (
		args                      interface{}
		isMany, isLike, isSqlNull bool
		primitiveType             int8
	)

	switch f := e.Interface().(type) {
	// Integer Fields
	case []int64,
		[]*int64,
		[]uint64,
		[]*uint64,
		[]null.Int64,
		[]*null.Int64,
		[]int32,
		[]*int32,
		[]uint32,
		[]*uint32,
		[]int16,
		[]*int16,
		[]uint16,
		[]*uint16,
		[]int8,
		[]*int8,
		[]uint8,
		[]*uint8,
		[]int,
		[]*int,
		[]uint,
		[]*uint,
		int64,
		uint64,
		null.Int64,
		int32,
		uint32,
		int16,
		uint16,
		int8,
		uint8,
		int,
		uint:
		primitiveType, isMany, isSqlNull, args = convertIntArgs(f)
		builderFunc(primitiveType, isLike, isMany, isSqlNull, fieldName, paramTagValue, dbTagValue, args)
		return
	// Float fields
	case []float64,
		[]*float64,
		[]float32,
		[]*float32,
		[]null.Float64,
		[]*null.Float64,
		float64,
		float32,
		null.Float64:
		primitiveType, isMany, args = convertFloatArgs(f)
		builderFunc(primitiveType, isLike, isMany, false, fieldName, paramTagValue, dbTagValue, args)
		return
	// String fields
	case []string,
		[]*string,
		[]null.String,
		[]*null.String,
		string,
		null.String:
		primitiveType, isMany, isLike, isSqlNull, args = convertStringArgs(f)
		builderFunc(primitiveType, isLike, isMany, isSqlNull, fieldName, paramTagValue, dbTagValue, args)
		return
	// Bool
	case []bool,
		[]*bool,
		[]null.Bool,
		[]*null.Bool,
		bool,
		null.Bool:
		primitiveType, isMany, args = convertBoolArgs(f)
		builderFunc(primitiveType, isLike, isMany, false, fieldName, paramTagValue, dbTagValue, args)
		return
	// time
	case []time.Time,
		[]*time.Time,
		[]null.Time,
		[]*null.Time,
		[]null.Date,
		[]*null.Date,
		time.Time,
		null.Time,
		null.Date:
		primitiveType, isMany, isSqlNull, args = convertTimeArgs(f)
		builderFunc(primitiveType, isLike, isMany, isSqlNull, fieldName, paramTagValue, dbTagValue, args)
		return
	}
}

func isTimeType(e reflect.Value) bool {
	return e.Kind() == reflect.Struct && (e.Type().String() == "null.Time" || e.Type().String() == "null.Date" || e.Type().String() == "time.Time")
}

func isNullType(e reflect.Value) bool {
	return e.Kind() == reflect.Struct &&
		(e.Type().String() == "null.String" ||
			e.Type().String() == "null.Bool" ||
			e.Type().String() == "null.Float64" ||
			e.Type().String() == "null.Int64" ||
			e.Type().String() == "null.Time" ||
			e.Type().String() == "null.Date")
}

func getNameFromStructTagOrOriginalName(fieldName string, v reflect.Value, i int) string {
	name := v.Type().Field(i).Tag.Get(fieldName)
	if len(name) > 0 {
		return name
	}
	return v.Type().Field(i).Name
}
