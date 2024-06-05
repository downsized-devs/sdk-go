package query

import (
	"reflect"
	"testing"

	"github.com/downsized-devs/sdk-go/null"
)

func Test_convertStringArgs(t *testing.T) {
	getStringPtr := func(x string) *string {
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
		wantIsLike        bool
		wantIsSqlNull     bool
		wantArgs          interface{}
	}{
		{
			name: "[]string",
			args: args{
				_f: []string{"akagami", "shanks"},
			},
			wantPrimitiveType: StringArr,
			wantIsMany:        true,
			wantIsLike:        false,
			wantIsSqlNull:     false,
			wantArgs:          []string{"akagami", "shanks"},
		},
		{
			name: "[]*string",
			args: args{
				_f: []*string{getStringPtr("akagami"), getStringPtr("shanks")},
			},
			wantPrimitiveType: StringArr,
			wantIsMany:        true,
			wantIsLike:        false,
			wantIsSqlNull:     false,
			wantArgs:          []*string{getStringPtr("akagami"), getStringPtr("shanks")},
		},
		{
			name: "string",
			args: args{
				_f: string("%yonko"),
			},
			wantPrimitiveType: String,
			wantIsMany:        false,
			wantIsLike:        true,
			wantIsSqlNull:     false,
			wantArgs:          string("%yonko"),
		},
		{
			name: "[]null.string",
			args: args{
				_f: []null.String{
					{String: "akagami", Valid: true},
					{String: "shanks", Valid: true},
				},
			},
			wantPrimitiveType: StringArr,
			wantIsMany:        true,
			wantIsLike:        false,
			wantIsSqlNull:     false,
			wantArgs:          []string{"akagami", "shanks"},
		},
		{
			name: "[]*null.string",
			args: args{
				_f: []*null.String{
					{String: *getStringPtr("akagami"), Valid: true},
					{String: *getStringPtr("shanks"), Valid: true},
				},
			},
			wantPrimitiveType: StringArr,
			wantIsMany:        true,
			wantIsLike:        false,
			wantIsSqlNull:     false,
			wantArgs:          []string{"akagami", "shanks"},
		},
		{
			name: "null.string",
			args: args{
				_f: null.String{
					Valid:  true,
					String: "%yonko",
				},
			},
			wantPrimitiveType: String,
			wantIsMany:        false,
			wantIsLike:        true,
			wantIsSqlNull:     false,
			wantArgs:          string("%yonko"),
		},
		{
			name: "null.string: SET to NULL",
			args: args{
				_f: null.String{
					SqlNull: true,
				},
			},
			wantPrimitiveType: String,
			wantIsMany:        false,
			wantIsLike:        false,
			wantIsSqlNull:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrimitiveType, gotIsMany, gotIsLike, gotIsSqlNull, gotArgs := convertStringArgs(tt.args._f)
			if gotPrimitiveType != tt.wantPrimitiveType {
				t.Errorf("convertStringArgs() gotPrimitiveType = %v, want %v", gotPrimitiveType, tt.wantPrimitiveType)
			}
			if gotIsMany != tt.wantIsMany {
				t.Errorf("convertStringArgs() gotIsMany = %v, want %v", gotIsMany, tt.wantIsMany)
			}
			if gotIsLike != tt.wantIsLike {
				t.Errorf("convertStringArgs() gotIsLike = %v, want %v", gotIsLike, tt.wantIsLike)
			}
			if gotIsSqlNull != tt.wantIsSqlNull {
				t.Errorf("convertTimeArgs() gotIsSqlNull = %v, want %v", gotIsSqlNull, tt.wantIsSqlNull)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("convertStringArgs() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}
