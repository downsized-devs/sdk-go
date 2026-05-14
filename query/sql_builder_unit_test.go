package query

import (
	"bytes"
	"testing"
	"time"

	"github.com/downsized-devs/sdk-go/null"
	mock_sql "github.com/downsized-devs/sdk-go/tests/mock/sql"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// Unit-test counterpart to sql_builder_test.go (which is integration-tagged
// and requires a live MySQL). These use a gomock-backed sql.Interface whose
// Leader().Rebind passes the query string through unchanged, exercising the
// same builder behavior without a database connection.

type unitTestParamNoTag struct {
	ID      int
	Name    string
	Weight  float64
	Height  float64
	Details interface{}
}

type unitTestParam struct {
	ID      int         `db:"id" param:"id" cursorField:"id"`
	Name    string      `db:"name" param:"name" cursorField:"name"`
	Weight  float64     `db:"weight" param:"weight" cursorField:"weight"`
	Height  float64     `db:"height" param:"height" cursorField:"height"`
	Details interface{} `db:"details" param:"details" cursorField:"details"`
}

type unitTestParamWithArray struct {
	ID     []int64   `param:"id" db:"id"`
	Name   []string  `param:"name" db:"name"`
	Length []float64 `param:"length" db:"length"`
}

type unitTestParamWithStringWildcard struct {
	Name    string `param:"name" db:"name"`
	NameOpt string `param:"name__opt" db:"name"`
}

type unitTestParamLimitAndPage struct {
	ID    int64    `param:"id" db:"id"`
	Names []string `param:"name" db:"name"`
	Limit int64    `param:"limit" db:"limit"`
	Page  int64    `param:"page" db:"page"`
}

type unitTestParamSortBy struct {
	ID     int64    `param:"id" db:"id"`
	SortBy []string `param:"sort_by" db:"sort_by"`
	Limit  int64    `param:"limit" db:"limit"`
	Page   int64    `param:"page" db:"page"`
}

type unitTestNullType struct {
	ID        null.Int64   `param:"id" db:"id"`
	Name      null.String  `param:"name" db:"name"`
	Length    null.Float64 `param:"length" db:"length"`
	Active    null.Bool    `param:"active" db:"active"`
	Birthday  null.Date    `param:"birthday" db:"birthday"`
	CreatedAt null.Time    `param:"created_at" db:"created_at"`
}

type unitTestComplexOperation struct {
	NameNIN []string `param:"name__nin" db:"name"`
	AgeGTE  int64    `param:"age__gte" db:"age"`
	AgeGT   int64    `param:"age__gt" db:"age"`
	AgeLTE  int64    `param:"age__lte" db:"age"`
	AgeLT   int64    `param:"age__lt" db:"age"`
	AgeNE   int64    `param:"age__ne" db:"age"`
	JobOPT  string   `param:"job__opt" db:"job"`
}

type unitTestAllKindsOfBooleans struct {
	Boolean1 bool         `param:"boolean_1" db:"boolean_1"`
	Boolean2 []bool       `param:"boolean_2" db:"boolean_2"`
	Boolean3 []*bool      `param:"boolean_3" db:"boolean_3"`
	Boolean4 null.Bool    `param:"boolean_4" db:"boolean_4"`
	Boolean5 []null.Bool  `param:"boolean_5" db:"boolean_5"`
	Boolean6 []*null.Bool `param:"boolean_6" db:"boolean_6"`
}

// newMockSQLInterface returns a mock sql.Interface whose Leader().Rebind is
// a pass-through (matches the default mysql bindvar style "?", so query text
// produced by the builder is preserved verbatim).
func newMockSQLInterface(t *testing.T) (*mock_sql.MockInterface, *mock_sql.MockCommand) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	db := mock_sql.NewMockInterface(ctrl)
	cmd := mock_sql.NewMockCommand(ctrl)
	db.EXPECT().Leader().Return(cmd).AnyTimes()
	cmd.EXPECT().Rebind(gomock.Any()).DoAndReturn(func(q string) string { return q }).AnyTimes()
	return db, cmd
}

func TestNewSQLQueryBuilder_Unit(t *testing.T) {
	db, _ := newMockSQLInterface(t)
	got := NewSQLQueryBuilder(db, "param", "db", nil)
	expected := &sqlClausebuilder{
		db:              db,
		fieldTag:        "cursorField",
		paramTag:        "param",
		dbTag:           "db",
		rawQuery:        bytes.NewBufferString(" WHERE 1=1"),
		rawUpdate:       bytes.NewBufferString(" SET"),
		paramToDBMap:    make(map[string]string),
		paramToFieldMap: make(map[string]string),
		sortableParams:  make(map[string]bool),
		aliasMap:        make(map[string]string),
	}
	assert.Equal(t, expected, got)
}

func Test_sqlClausebuilder_Build_Unit(t *testing.T) {
	boolTrue := true

	type args struct {
		param interface{}
		opt   *Option
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []interface{}
		want2   string
		want3   []interface{}
		wantErr bool
	}{
		{
			name: "all primitive param used",
			args: args{param: &unitTestParam{
				ID: 1, Name: "test", Weight: 123.45, Height: 67.89, Details: "simple string",
			}},
			want:  " WHERE 1=1 AND id=? AND name=? AND weight=? AND height=? AND details=?;",
			want1: []interface{}{1, "test", 123.45, 67.89, "simple string"},
			want2: " WHERE 1=1 AND id=? AND name=? AND weight=? AND height=? AND details=?;",
			want3: []interface{}{1, "test", 123.45, 67.89, "simple string"},
		},
		{
			name:  "all primitive param used with no tags",
			args:  args{param: &unitTestParamNoTag{ID: 1, Name: "test"}},
			want:  " WHERE 1=1;",
			want2: " WHERE 1=1;",
		},
		{
			name:  "all primitive param zero value",
			args:  args{param: &unitTestParam{}},
			want:  " WHERE 1=1;",
			want2: " WHERE 1=1;",
		},
		{
			name: "array with IN",
			args: args{param: &unitTestParamWithArray{
				ID:     []int64{1, 2},
				Name:   []string{"jack", "garland"},
				Length: []float64{1.23},
			}},
			want:  " WHERE 1=1 AND id IN (?, ?) AND name IN (?, ?) AND length IN (?);",
			want1: []interface{}{int64(1), int64(2), "jack", "garland", float64(1.23)},
			want2: " WHERE 1=1 AND id IN (?, ?) AND name IN (?, ?) AND length IN (?);",
			want3: []interface{}{int64(1), int64(2), "jack", "garland", float64(1.23)},
		},
		{
			name:  "empty array with IN",
			args:  args{param: &unitTestParamWithArray{}},
			want:  " WHERE 1=1;",
			want2: " WHERE 1=1;",
		},
		{
			name:  "string LIKE",
			args:  args{param: &unitTestParamWithStringWildcard{Name: "garland%"}},
			want:  " WHERE 1=1 AND name LIKE ?;",
			want1: []interface{}{"garland%"},
			want2: " WHERE 1=1 AND name LIKE ?;",
			want3: []interface{}{"garland%"},
		},
		{
			name:  "string LIKE opt",
			args:  args{param: &unitTestParamWithStringWildcard{NameOpt: "garland%"}},
			want:  " WHERE 1=1 OR name LIKE ?;",
			want1: []interface{}{"garland%"},
			want2: " WHERE 1=1 OR name LIKE ?;",
			want3: []interface{}{"garland%"},
		},
		{
			name: "limit and page",
			args: args{param: &unitTestParamLimitAndPage{
				ID: 1, Names: []string{"jack", "garland"}, Limit: 10, Page: 1,
			}},
			want:  " WHERE 1=1 AND id=? AND name IN (?, ?) LIMIT 0, 10;",
			want1: []interface{}{int64(1), "jack", "garland"},
			want2: " WHERE 1=1 AND id=? AND name IN (?, ?);",
			want3: []interface{}{int64(1), "jack", "garland"},
		},
		{
			name: "sort by",
			args: args{param: &unitTestParamSortBy{
				ID: 1, SortBy: []string{"id"}, Limit: 10, Page: 1,
			}},
			want:  " WHERE 1=1 AND id=? ORDER BY id ASC LIMIT 0, 10;",
			want1: []interface{}{int64(1)},
			want2: " WHERE 1=1 AND id=?;",
			want3: []interface{}{int64(1)},
		},
		{
			name: "sort by descending",
			args: args{param: &unitTestParamSortBy{
				ID: 1, SortBy: []string{"-id"}, Limit: 10, Page: 1,
			}},
			want:  " WHERE 1=1 AND id=? ORDER BY id DESC LIMIT 0, 10;",
			want1: []interface{}{int64(1)},
			want2: " WHERE 1=1 AND id=?;",
			want3: []interface{}{int64(1)},
		},
		{
			name: "null type valid",
			args: args{param: &unitTestNullType{
				ID:        null.Int64{Valid: true},
				Name:      null.String{Valid: true},
				Length:    null.Float64{Valid: true},
				Active:    null.Bool{Valid: true},
				Birthday:  null.Date{Valid: true},
				CreatedAt: null.Time{Valid: true},
			}},
			want:  " WHERE 1=1 AND id=? AND name=? AND length=? AND active=? AND birthday=? AND created_at=?;",
			want1: []interface{}{int64(0), "", float64(0), false, time.Time{}, time.Time{}},
			want2: " WHERE 1=1 AND id=? AND name=? AND length=? AND active=? AND birthday=? AND created_at=?;",
			want3: []interface{}{int64(0), "", float64(0), false, time.Time{}, time.Time{}},
		},
		{
			name:  "null type invalid",
			args:  args{param: &unitTestNullType{}},
			want:  " WHERE 1=1;",
			want2: " WHERE 1=1;",
		},
		{
			name: "complex operations",
			args: args{param: &unitTestComplexOperation{
				NameNIN: []string{"jack", "garland"},
				AgeGTE:  11, AgeGT: 10, AgeLTE: 9, AgeLT: 10, AgeNE: 100,
				JobOPT: "knight",
			}},
			want:  " WHERE 1=1 AND name NOT IN (?, ?) AND age>=? AND age>? AND age<=? AND age<? AND age<>? OR job=?;",
			want1: []interface{}{"jack", "garland", int64(11), int64(10), int64(9), int64(10), int64(100), "knight"},
			want2: " WHERE 1=1 AND name NOT IN (?, ?) AND age>=? AND age>? AND age<=? AND age<? AND age<>? OR job=?;",
			want3: []interface{}{"jack", "garland", int64(11), int64(10), int64(9), int64(10), int64(100), "knight"},
		},
		{
			name: "all kinds of valid booleans",
			args: args{param: &unitTestAllKindsOfBooleans{
				Boolean1: true,
				Boolean2: []bool{true},
				Boolean3: []*bool{&boolTrue},
				Boolean4: null.BoolFrom(true),
				Boolean5: []null.Bool{null.BoolFrom(true)},
				Boolean6: []*null.Bool{{Valid: true, Bool: true}},
			}},
			want:  " WHERE 1=1 AND boolean_1=? AND boolean_2 IN (?) AND boolean_3 IN (?) AND boolean_4=? AND boolean_5 IN (?) AND boolean_6 IN (?);",
			want1: []interface{}{true, true, &boolTrue, true, true, true},
			want2: " WHERE 1=1 AND boolean_1=? AND boolean_2 IN (?) AND boolean_3 IN (?) AND boolean_4=? AND boolean_5 IN (?) AND boolean_6 IN (?);",
			want3: []interface{}{true, true, &boolTrue, true, true, true},
		},
		{
			name: "option disable limit",
			args: args{
				param: &unitTestParamLimitAndPage{ID: 1, Names: []string{"jack", "garland"}, Limit: 10, Page: 1},
				opt:   &Option{DisableLimit: true},
			},
			want:  " WHERE 1=1 AND id=? AND name IN (?, ?);",
			want1: []interface{}{int64(1), "jack", "garland"},
			want2: " WHERE 1=1 AND id=? AND name IN (?, ?);",
			want3: []interface{}{int64(1), "jack", "garland"},
		},
		{
			name: "option is active",
			args: args{
				param: &unitTestParamLimitAndPage{ID: 1, Names: []string{"jack", "garland"}, Limit: 10, Page: 1},
				opt:   &Option{IsActive: true},
			},
			want:  " WHERE 1=1 AND status=1 AND id=? AND name IN (?, ?) LIMIT 0, 10;",
			want1: []interface{}{int64(1), "jack", "garland"},
			want2: " WHERE 1=1 AND status=1 AND id=? AND name IN (?, ?);",
			want3: []interface{}{int64(1), "jack", "garland"},
		},
		{
			name: "option is inactive",
			args: args{
				param: &unitTestParamLimitAndPage{ID: 1, Names: []string{"jack", "garland"}, Limit: 10, Page: 1},
				opt:   &Option{IsInactive: true},
			},
			want:  " WHERE 1=1 AND status=0 AND id=? AND name IN (?, ?) LIMIT 0, 10;",
			want1: []interface{}{int64(1), "jack", "garland"},
			want2: " WHERE 1=1 AND status=0 AND id=? AND name IN (?, ?);",
			want3: []interface{}{int64(1), "jack", "garland"},
		},
		{
			name:    "param is nil",
			args:    args{param: nil},
			wantErr: true,
		},
		{
			name:    "param is non-pointer",
			args:    args{param: unitTestParam{ID: 1}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, _ := newMockSQLInterface(t)
			qb := NewSQLQueryBuilder(db, "param", "db", tt.args.opt)
			got, got1, got2, got3, err := qb.Build(tt.args.param)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)
			assert.Equal(t, tt.want3, got3)
		})
	}
}

func Test_sqlClausebuilder_AddGroupByQuery_Unit(t *testing.T) {
	param := &unitTestParamLimitAndPage{ID: 1, Names: []string{"jack", "garland"}, Limit: 10, Page: 1}

	tests := []struct {
		name    string
		columns []string
		want    string
	}{
		{"single", []string{"id"}, " WHERE 1=1 AND id=? AND name IN (?, ?) GROUP BY id LIMIT 0, 10;"},
		{"multiple", []string{"id", "name"}, " WHERE 1=1 AND id=? AND name IN (?, ?) GROUP BY id, name LIMIT 0, 10;"},
		{"empty", []string{}, " WHERE 1=1 AND id=? AND name IN (?, ?) LIMIT 0, 10;"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, _ := newMockSQLInterface(t)
			qb := NewSQLQueryBuilder(db, "param", "db", nil)
			qb.AddGroupByQuery(tt.columns...)
			got, _, _, _, err := qb.Build(param)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_sqlClausebuilder_AddPrefixQuery_Unit(t *testing.T) {
	param := &unitTestParamLimitAndPage{ID: 1, Names: []string{"jack", "garland"}, Limit: 10, Page: 1}

	tests := []struct {
		name        string
		prefixQuery string
		want        string
	}{
		{"non-empty", "active = 1", " WHERE 1=1 AND active = 1 AND id=? AND name IN (?, ?) LIMIT 0, 10;"},
		{"empty", "", " WHERE 1=1 AND id=? AND name IN (?, ?) LIMIT 0, 10;"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, _ := newMockSQLInterface(t)
			qb := NewSQLQueryBuilder(db, "param", "db", nil)
			qb.AddPrefixQuery(tt.prefixQuery)
			got, _, _, _, err := qb.Build(param)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}

	t.Run("prefix with bindvars produces sqlx.In error", func(t *testing.T) {
		db, _ := newMockSQLInterface(t)
		qb := NewSQLQueryBuilder(db, "param", "db", nil)
		qb.AddPrefixQuery("status IN (?)")
		_, _, _, _, err := qb.Build(param)
		assert.Error(t, err)
	})
}

func Test_sqlClausebuilder_AddSuffixQuery_Unit(t *testing.T) {
	param := &unitTestParamLimitAndPage{ID: 1, Names: []string{"jack", "garland"}, Limit: 10, Page: 1}

	tests := []struct {
		name        string
		suffixQuery string
		want        string
	}{
		{"non-empty", "GROUP BY name", " WHERE 1=1 AND id=? AND name IN (?, ?) LIMIT 0, 10 GROUP BY name;"},
		{"empty", "", " WHERE 1=1 AND id=? AND name IN (?, ?) LIMIT 0, 10;"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, _ := newMockSQLInterface(t)
			qb := NewSQLQueryBuilder(db, "param", "db", nil)
			qb.AddSuffixQuery(tt.suffixQuery)
			got, _, _, _, err := qb.Build(param)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}

	t.Run("suffix with bindvars produces sqlx.In error", func(t *testing.T) {
		db, _ := newMockSQLInterface(t)
		qb := NewSQLQueryBuilder(db, "param", "db", nil)
		qb.AddSuffixQuery("status IN (?)")
		_, _, _, _, err := qb.Build(param)
		assert.Error(t, err)
	})
}

func Test_sqlClausebuilder_AddAliasPrefix_Unit(t *testing.T) {
	t.Run("with tag", func(t *testing.T) {
		param := &unitTestParamLimitAndPage{ID: 1, Names: []string{"jack", "garland"}, Limit: 10, Page: 1}
		db, _ := newMockSQLInterface(t)
		qb := NewSQLQueryBuilder(db, "param", "db", nil)
		qb.AddAliasPrefix("idk", param)
		got, _, _, _, err := qb.Build(param)
		assert.NoError(t, err)
		assert.Equal(t, " WHERE 1=1 AND idk.id=? AND idk.name IN (?, ?) LIMIT 0, 10;", got)
	})

	t.Run("no tag", func(t *testing.T) {
		param := &unitTestParamNoTag{ID: 1, Name: "jack"}
		db, _ := newMockSQLInterface(t)
		qb := NewSQLQueryBuilder(db, "param", "db", nil)
		qb.AddAliasPrefix("idk", param)
		got, _, _, _, err := qb.Build(param)
		assert.NoError(t, err)
		assert.Equal(t, " WHERE 1=1;", got)
	})

	t.Run("non-pointer panics", func(t *testing.T) {
		db, _ := newMockSQLInterface(t)
		qb := NewSQLQueryBuilder(db, "param", "db", nil)
		assert.Panics(t, func() {
			qb.AddAliasPrefix("idk", unitTestParamNoTag{ID: 1})
		})
	})
}

func Test_sqlClausebuilder_BuildUpdate_Unit(t *testing.T) {
	type testUpdateStruct struct {
		FirstName string      `param:"first_name" db:"first_name"`
		LastName  null.String `param:"last_name" db:"last_name"`
		Age       int64       `param:"age" db:"age"`
		Children  null.Int64  `param:"children" db:"children"`
	}
	type testWhereStruct struct {
		ID        int64  `param:"id" db:"id"`
		FirstName string `param:"first_name" db:"first_name"`
		Limit     int64  `param:"limit" db:"limit"`
		Page      int64  `param:"page" db:"page"`
		SortBy    string `param:"sort_by" db:"sort_by"`
	}

	type args struct {
		update interface{}
		where  interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []interface{}
		wantErr bool
	}{
		{name: "update nil", args: args{update: nil, where: &testWhereStruct{ID: 1}}, wantErr: true},
		{name: "where nil", args: args{update: &testUpdateStruct{FirstName: "jack"}, where: nil}, wantErr: true},
		{name: "update empty", args: args{update: &testUpdateStruct{}, where: &testWhereStruct{ID: 1}}, wantErr: true},
		{name: "where empty", args: args{update: &testUpdateStruct{FirstName: "jack"}, where: &testWhereStruct{}}, wantErr: true},
		{
			name: "update without null fields",
			args: args{
				update: &testUpdateStruct{FirstName: "jack", Age: 21},
				where:  &testWhereStruct{ID: 1},
			},
			want:  " SET first_name=?, age=? WHERE 1=1 AND id=?;",
			want1: []interface{}{"jack", int64(21), int64(1)},
		},
		{
			name: "update with null fields valid zero",
			args: args{
				update: &testUpdateStruct{LastName: null.StringFrom(""), Children: null.Int64From(0)},
				where:  &testWhereStruct{ID: 1},
			},
			want:  " SET last_name=?, children=? WHERE 1=1 AND id=?;",
			want1: []interface{}{"", int64(0), int64(1)},
		},
		{
			name: "update set NULL via SqlNull",
			args: args{
				update: &unitTestNullType{
					Birthday:  null.Date{SqlNull: true},
					CreatedAt: null.Time{SqlNull: true},
					Name:      null.String{SqlNull: true},
				},
				where: &testWhereStruct{ID: 1},
			},
			// Output order follows the struct field declaration order in
			// unitTestNullType: Name, then Birthday, then CreatedAt.
			want:  " SET name=NULL, birthday=NULL, created_at=NULL WHERE 1=1 AND id=?;",
			want1: []interface{}{int64(1)},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, _ := newMockSQLInterface(t)
			qb := NewSQLQueryBuilder(db, "param", "db", nil)
			got, got1, err := qb.BuildUpdate(tt.args.update, tt.args.where)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func Test_sqlClausebuilder_pagePagination(t *testing.T) {
	// Direct test of the helper for the page>0 || limit>0 branch.
	s := &sqlClausebuilder{page: 2, limit: 5}
	s.pagePagination()
	assert.Equal(t, " LIMIT 5, 5", s.paginationClause)

	s2 := &sqlClausebuilder{}
	s2.pagePagination()
	assert.Empty(t, s2.paginationClause)
}

func Test_sqlClausebuilder_getBindVar(t *testing.T) {
	s := &sqlClausebuilder{}
	assert.Equal(t, "?", s.getBindVar())
}
