package query

import (
	"reflect"
	"testing"
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
				p:             reflect.ValueOf(reflect.Ptr),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			traverseOnParam(tt.args.paramTagName, tt.args.dbTagName, tt.args.fieldTagName, tt.args.fieldName, tt.args.paramTagValue, tt.args.dbTagValue, tt.args.aliasMap, tt.args.p, tt.args.builderFunc)
		})
	}
}
