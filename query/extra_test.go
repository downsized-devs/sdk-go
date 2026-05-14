package query

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Plugs the remaining coverage gaps in query/query.go and query/converter.go
// without depending on a database.

func TestFormatQueryForRows_NoRows(t *testing.T) {
	_, _, err := FormatQueryForRows(context.Background(), "INSERT INTO t (a)", nil)
	assert.Error(t, err)
}

func TestFormatQueryForRows_NoCols(t *testing.T) {
	// Row exists but contains zero columns.
	_, _, err := FormatQueryForRows(context.Background(), "INSERT INTO t", [][]interface{}{{}})
	assert.Error(t, err)
}

// Hits the case in traverseOnParam where a top-level time.Time is passed
// directly (Kind=Struct, isTimeType=true).
func Test_traverseOnParam_TimeAtTopLevel(t *testing.T) {
	var called bool
	now := time.Now()
	traverseOnParam("param", "db", cursorField, "$", "p", "d", map[string]string{}, reflect.ValueOf(now), func(int8, bool, bool, bool, string, string, string, interface{}) {
		called = true
	})
	assert.True(t, called)
}

// Hits the dbTagValue == "-" skip branch in traverseOnParam.
func Test_traverseOnParam_DashTagSkipped(t *testing.T) {
	type skipStruct struct {
		Skipped string `param:"skipped" db:"-"`
		Kept    string `param:"kept" db:"kept"`
	}
	calls := map[string]string{}
	s := skipStruct{Skipped: "no", Kept: "yes"}
	v := reflect.ValueOf(&s).Elem()
	traverseOnParam("param", "db", cursorField, "$", "", "", map[string]string{}, v, func(_ int8, _, _, _ bool, fieldName, paramTag, dbTag string, args interface{}) {
		calls[paramTag] = dbTag
	})
	// The "skipped" field should never have been visited (dash skip).
	_, ok := calls["skipped"]
	assert.False(t, ok)
	assert.Equal(t, "kept", calls["kept"])
}
