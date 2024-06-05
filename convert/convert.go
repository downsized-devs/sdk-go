package convert

import (
	"reflect"
	"strings"

	"github.com/cstockton/go-conv"
	"github.com/downsized-devs/sdk-go/codes"
	"github.com/downsized-devs/sdk-go/errors"
)

func ToInt64(i interface{}) (int64, error) {
	val, err := conv.Int64(i)
	if err != nil {
		return val, errors.NewWithCode(codes.CodeInvalidValue, "failed to parse value: %v. with err: %v", i, err)
	}

	return val, nil
}

func ToArrInt64(i interface{}) ([]int64, error) {
	val := make([]int64, 0)
	// if the input is not a slice, convert the input to slice with one member
	if reflect.TypeOf(i).Kind() != reflect.Slice {
		singleVal, err := conv.Int64(i)
		if err != nil {
			return []int64{}, errors.NewWithCode(codes.CodeInvalidValue, "failed to parse value: %v. with err: %v", i, err)
		}

		return append(val, singleVal), nil
	}

	s := reflect.ValueOf(i)
	for i := 0; i < s.Len(); i++ {
		singleVal, err := conv.Int64(s.Index(i).Interface())
		if err != nil {
			return []int64{}, errors.NewWithCode(codes.CodeInvalidValue, "failed to parse value: %v. with err: %v", i, err)
		}

		val = append(val, singleVal)
	}

	return val, nil
}

func ToFloat64(i interface{}) (float64, error) {
	v := reflect.ValueOf(i)
	v = reflect.Indirect(v)

	if !v.IsValid() {
		return 0, errors.NewWithCode(codes.CodeInvalidValue, "failed to parse value: %v. Invalid value", i)
	}

	floatType := reflect.TypeOf(float64(0))
	if v.Type().ConvertibleTo(floatType) {
		fv := v.Convert(floatType)
		return fv.Float(), nil
	}

	// if value is not convertible in Go vanilla, convert using library
	val, err := conv.Float64(i)
	if err != nil {
		return val, errors.NewWithCode(codes.CodeInvalidValue, "failed to parse value: %v. with err: %v", i, err)
	}

	return val, nil
}

// Convert any data to string. This function will always return nil error
func ToString(i interface{}) (string, error) {
	// conv.String will always return nil error. See here for reference:
	// https://github.com/cstockton/go-conv/blob/57aa63165ab27764a034b79724d9fc1b876a9bba/internal/refconv/string.go#L18
	val, _ := conv.String(i)

	return val, nil
}

// Convert integer to its character representation (eg. 1=A 3=C).
// Support multiple character for number greater than 26 (eg. 27=AA 703=AAA)
func IntToChar(i int) string {
	char := ""
	i--

	if i < 0 {
		return char
	}

	if firstcharint := i / 26; firstcharint > 0 {
		char += IntToChar(firstcharint)
		char += string(rune('A' + (i % 26)))
	} else {
		char += string(rune('A' + i))
	}

	return char
}

// Convert Pascal Case string (PascalCase) to Camel Case string (camelCase)
func PascalCaseToCamelCase(s string) string {
	if len(s) < 1 {
		return ""
	}

	first := strings.ToLower(s[:1])
	rest := s[1:]
	return first + rest
}

// Convert Text to camel case string with lower case in beginning
// ex: "Additional Item" --> "additionalItem"
func ToCamelCase(input string) string {
	return toCamelInitCase(input, false)
}

// Convert text to camle case string with upper case in beginning
// ex: "additional item" --> "AdditionalItem"
func ToPascalCase(input string) string {
	return toCamelInitCase(input, true)
}

func RomanToInt64(roman string) int64 {
	romansMap := map[string]int64{
		"i": 1,
		"v": 5,
		"x": 10,
		"l": 50,
		"c": 100,
		"d": 500,
		"m": 1000,
	}

	roman = strings.ToLower(roman)
	roman = strings.ReplaceAll(roman, "iv", "iiii")
	roman = strings.ReplaceAll(roman, "ix", "viiii")
	roman = strings.ReplaceAll(roman, "xl", "xxxx")
	roman = strings.ReplaceAll(roman, "xc", "lxxxx")
	roman = strings.ReplaceAll(roman, "cd", "cccc")
	roman = strings.ReplaceAll(roman, "cm", "dcccc")

	var result int64
	for _, romanChar := range roman {
		result += romansMap[string(romanChar)]
	}

	return result
}

type romanPair struct {
	value int64
	roman string
}

func Int64ToRoman(num int64) string {
	romanPairs := []romanPair{
		{1000, "M"},
		{900, "CM"},
		{500, "D"},
		{400, "CD"},
		{100, "C"},
		{90, "XC"},
		{50, "L"},
		{40, "XL"},
		{10, "X"},
		{9, "IX"},
		{5, "V"},
		{4, "IV"},
		{1, "I"},
	}

	var result string
	for _, romanPair := range romanPairs {
		for num >= romanPair.value {
			result += romanPair.roman
			num -= romanPair.value
		}
	}

	return result
}

// Converts a string to CamelCase
func toCamelInitCase(s string, initCase bool) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	n := strings.Builder{}
	n.Grow(len(s))
	capNext := initCase
	for i, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'
		if capNext {
			if vIsLow {
				v += 'A'
				v -= 'a'
			}
		} else if i == 0 {
			if vIsCap {
				v += 'a'
				v -= 'A'
			}
		}
		if vIsCap || vIsLow {
			n.WriteByte(v)
			capNext = false
		} else if vIsNum := v >= '0' && v <= '9'; vIsNum {
			n.WriteByte(v)
			capNext = true
		} else {
			capNext = v == '_' || v == ' ' || v == '-' || v == '.'
		}
	}
	return n.String()
}
