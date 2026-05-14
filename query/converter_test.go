package query

import (
	"reflect"
	"testing"
	"time"

	"github.com/downsized-devs/sdk-go/null"
)

func Test_traverseOnParam(t *testing.T) {
	type args struct {
		paramTagName  string
		dbTagName     string
		fieldTagName  string
		fieldName     string
		paramTagValue string
		dbTagValue    string
		aliasMap      map[string]string
		p             reflect.Value
		builderFunc   builderFunction
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "transverse on param default",
			args: args{
				paramTagName:  "",
				dbTagName:     "",
				fieldTagName:  "",
				fieldName:     "",
				paramTagValue: "",
				dbTagValue:    "",
				aliasMap:      map[string]string{},
				p:             reflect.ValueOf(reflect.Pointer),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			traverseOnParam(tt.args.paramTagName, tt.args.dbTagName, tt.args.fieldTagName, tt.args.fieldName, tt.args.paramTagValue, tt.args.dbTagValue, tt.args.aliasMap, tt.args.p, tt.args.builderFunc)
		})
	}
}

func Test_isTimeType(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  bool
	}{
		{
			name:  "time.Time is time type",
			value: time.Time{},
			want:  true,
		},
		{
			name:  "null.Time is time type",
			value: null.Time{},
			want:  true,
		},
		{
			name:  "null.Date is time type",
			value: null.Date{},
			want:  true,
		},
		{
			name:  "string is not time type",
			value: "",
			want:  false,
		},
		{
			name:  "int is not time type",
			value: 0,
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.value)
			if got := isTimeType(v); got != tt.want {
				t.Errorf("isTimeType(%T) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func Test_isNullType(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  bool
	}{
		{
			name:  "null.String is null type",
			value: null.String{},
			want:  true,
		},
		{
			name:  "null.Bool is null type",
			value: null.Bool{},
			want:  true,
		},
		{
			name:  "null.Float64 is null type",
			value: null.Float64{},
			want:  true,
		},
		{
			name:  "null.Int64 is null type",
			value: null.Int64{},
			want:  true,
		},
		{
			name:  "null.Time is null type",
			value: null.Time{},
			want:  true,
		},
		{
			name:  "null.Date is null type",
			value: null.Date{},
			want:  true,
		},
		{
			name:  "string is not null type",
			value: "",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.value)
			if got := isNullType(v); got != tt.want {
				t.Errorf("isNullType(%T) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func Test_getNameFromStructTagOrOriginalName(t *testing.T) {
	type Tagged struct {
		Name string `param:"custom_name"`
	}
	type Untagged struct {
		Name string
	}

	t.Run("returns tag value when tag exists", func(t *testing.T) {
		v := reflect.ValueOf(Tagged{})
		got := getNameFromStructTagOrOriginalName("param", v, 0)
		if got != "custom_name" {
			t.Errorf("getNameFromStructTagOrOriginalName() = %q, want %q", got, "custom_name")
		}
	})

	t.Run("returns field name when no tag", func(t *testing.T) {
		v := reflect.ValueOf(Untagged{})
		got := getNameFromStructTagOrOriginalName("param", v, 0)
		if got != "Name" {
			t.Errorf("getNameFromStructTagOrOriginalName() = %q, want %q", got, "Name")
		}
	})
}
