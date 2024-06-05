package query

import (
	"reflect"
	"testing"
	"time"

	"github.com/downsized-devs/sdk-go/null"
)

func Test_convertTimeArgs(t *testing.T) {
	loc := time.UTC
	mockTime := time.Date(2022, 9, 7, 20, 42, 0, 0, loc)

	getTimePtr := func(x time.Time) *time.Time {
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
			name: "time.Time",
			args: args{
				_f: mockTime,
			},
			wantPrimitiveType: Time,
			wantIsMany:        false,
			wantIsSqlNull:     false,
			wantArgs:          mockTime,
		},
		{
			name: "[]time.Time",
			args: args{
				_f: []time.Time{mockTime},
			},
			wantPrimitiveType: TimeArr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []time.Time{mockTime},
		},
		{
			name: "[]*time.Time",
			args: args{
				_f: []*time.Time{getTimePtr(mockTime)},
			},
			wantPrimitiveType: TimeArr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []*time.Time{getTimePtr(mockTime)},
		},
		{
			name: "[]null.Time",
			args: args{
				_f: []null.Time{
					{Time: mockTime, Valid: true},
				},
			},
			wantPrimitiveType: TimeArr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []time.Time{mockTime},
		},
		{
			name: "[]*null.Time",
			args: args{
				_f: []*null.Time{
					{Time: *getTimePtr(mockTime), Valid: true},
				},
			},
			wantPrimitiveType: TimeArr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []time.Time{mockTime},
		},
		{
			name: "[]null.Date",
			args: args{
				_f: []null.Date{
					{Time: mockTime, Valid: true},
				},
			},
			wantPrimitiveType: TimeArr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []time.Time{mockTime},
		},
		{
			name: "[]*null.Date",
			args: args{
				_f: []*null.Date{
					{Time: *getTimePtr(mockTime), Valid: true},
				},
			},
			wantPrimitiveType: TimeArr,
			wantIsMany:        true,
			wantIsSqlNull:     false,
			wantArgs:          []time.Time{mockTime},
		},
		{
			name: "null.Time",
			args: args{
				_f: null.Time{Time: mockTime, Valid: true},
			},
			wantPrimitiveType: Time,
			wantIsMany:        false,
			wantIsSqlNull:     false,
			wantArgs:          mockTime,
		},
		{
			name: "null.Date",
			args: args{
				_f: null.Date{Time: mockTime, Valid: true},
			},
			wantPrimitiveType: Time,
			wantIsMany:        false,
			wantIsSqlNull:     false,
			wantArgs:          mockTime,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrimitiveType, gotIsMany, gotIsSqlNull, gotArgs := convertTimeArgs(tt.args._f)
			if gotPrimitiveType != tt.wantPrimitiveType {
				t.Errorf("convertTimeArgs() gotPrimitiveType = %v, want %v", gotPrimitiveType, tt.wantPrimitiveType)
			}
			if gotIsMany != tt.wantIsMany {
				t.Errorf("convertTimeArgs() gotIsMany = %v, want %v", gotIsMany, tt.wantIsMany)
			}
			if gotIsSqlNull != tt.wantIsSqlNull {
				t.Errorf("convertTimeArgs() gotIsSqlNull = %v, want %v", gotIsSqlNull, tt.wantIsSqlNull)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("convertTimeArgs() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}
