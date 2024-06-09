//go:build integration
// +build integration

package query

import (
	"bytes"
	"testing"
	"time"

	"github.com/downsized-devs/sdk-go/instrument"
	"github.com/downsized-devs/sdk-go/null"
	"github.com/downsized-devs/sdk-go/sql"
	"github.com/stretchr/testify/assert"
)

type TestParamNoTag struct {
	ID      int
	Name    string
	Weight  float64
	Height  float64
	Details interface{}
}

type TestParam struct {
	ID      int         `db:"id" param:"id" cursorField:"id"`
	Name    string      `db:"name" param:"name" cursorField:"name"`
	Weight  float64     `db:"weight" param:"weight" cursorField:"weight"`
	Height  float64     `db:"height" param:"height" cursorField:"height"`
	Details interface{} `db:"details" param:"details" cursorField:"details"`
}

type TestParamWithArray struct {
	ID     []int64   `param:"id" db:"id"`
	Name   []string  `param:"name" db:"name"`
	Length []float64 `param:"length" db:"length"`
}

type TestParamWithStringWildcard struct {
	Name    string `param:"name" db:"name"`
	NameOpt string `param:"name__opt" db:"name"`
}

type TestParamLimitAndPage struct {
	ID    int64    `param:"id" db:"id"`
	Names []string `param:"name" db:"name"`
	Limit int64    `param:"limit" db:"limit"`
	Page  int64    `param:"page" db:"page"`
}

type TestParamSortBy struct {
	ID     int64    `param:"id" db:"id"`
	SortBy []string `param:"sort_by" db:"sort_by"`
	Limit  int64    `param:"limit" db:"limit"`
	Page   int64    `param:"page" db:"page"`
}

type TestNullType struct {
	ID        null.Int64   `param:"id" db:"id"`
	Name      null.String  `param:"name" db:"name"`
	Length    null.Float64 `param:"length" db:"length"`
	Active    null.Bool    `param:"active" db:"active"`
	Birthday  null.Date    `param:"birthday" db:"birthday"`
	CreatedAt null.Time    `param:"created_at" db:"created_at"`
}

type TestComplexOperation struct {
	NameNIN []string `param:"name__nin" db:"name"`
	AgeGTE  int64    `param:"age__gte" db:"age"`
	AgeGT   int64    `param:"age__gt" db:"age"`
	AgeLTE  int64    `param:"age__lte" db:"age"`
	AgeLT   int64    `param:"age__lt" db:"age"`
	AgeNE   int64    `param:"age__ne" db:"age"`
	JobOPT  string   `param:"job__opt" db:"job"`
}

type TestAllKindsOfBooleans struct {
	Boolean1 bool         `param:"boolean_1" db:"boolean_1"`
	Boolean2 []bool       `param:"boolean_2" db:"boolean_2"`
	Boolean3 []*bool      `param:"boolean_3" db:"boolean_3"`
	Boolean4 null.Bool    `param:"boolean_4" db:"boolean_4"`
	Boolean5 []null.Bool  `param:"boolean_5" db:"boolean_5"`
	Boolean6 []*null.Bool `param:"boolean_6" db:"boolean_6"`
}

func initTestDatabase() sql.Interface {
	return sql.Init(sql.Config{
		Driver: "mysql",
		Leader: sql.ConnConfig{
			Host:     "127.0.0.1",
			Port:     3306,
			DB:       "",
			User:     "root",
			Password: "password",
		},
		Follower: sql.ConnConfig{
			Host:     "127.0.0.1",
			Port:     3306,
			DB:       "",
			User:     "root",
			Password: "password",
		},
	}, logger.Init(logger.Config{Level: "debug"}),
		instrument.Init(instrument.Config{}))
}

func TestNewSQLQueryBuilder(t *testing.T) {
	type args struct {
		db       sql.Interface
		paramTag string
		dbTag    string
		options  *Option
	}
	tests := []struct {
		name    string
		args    args
		want    *sqlClausebuilder
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				db:       nil,
				paramTag: "param",
				dbTag:    "db",
			},
			want: &sqlClausebuilder{
				fieldTag:        "cursorField",
				paramTag:        "param",
				dbTag:           "db",
				rawQuery:        bytes.NewBufferString(" WHERE 1=1"),
				rawUpdate:       bytes.NewBufferString(" SET"),
				paramToDBMap:    make(map[string]string),
				paramToFieldMap: make(map[string]string),
				sortableParams:  make(map[string]bool),
				aliasMap:        make(map[string]string),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSQLQueryBuilder(tt.args.db, tt.args.paramTag, tt.args.dbTag, tt.args.options)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_sqlClausebuilder_Build(t *testing.T) {
	BooleanTrue := true

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
			args: args{param: &TestParam{
				ID:      1,
				Name:    "test",
				Weight:  float64(123.45),
				Height:  float64(67.89),
				Details: "simple string",
			}},
			want:    " WHERE 1=1 AND id=? AND name=? AND weight=? AND height=? AND details=?;",
			want1:   []interface{}{1, "test", float64(123.45), float64(67.89), "simple string"},
			want2:   " WHERE 1=1 AND id=? AND name=? AND weight=? AND height=? AND details=?;",
			want3:   []interface{}{1, "test", float64(123.45), float64(67.89), "simple string"},
			wantErr: false,
		},
		{
			name: "all primitive param used with no tags",
			args: args{param: &TestParamNoTag{
				ID:      1,
				Name:    "test",
				Weight:  float64(123.45),
				Height:  float64(67.89),
				Details: "simple string",
			}},
			want:    " WHERE 1=1;",
			want1:   nil,
			want2:   " WHERE 1=1;",
			want3:   nil,
			wantErr: false,
		},
		{
			name: "all primitive param zero value",
			args: args{param: &TestParam{
				ID:      0,
				Name:    "",
				Weight:  float64(0),
				Details: "",
			}},
			want:    " WHERE 1=1;",
			want1:   nil,
			want2:   " WHERE 1=1;",
			want3:   nil,
			wantErr: false,
		},
		{
			name: "test array with IN",
			args: args{param: &TestParamWithArray{
				ID:     []int64{1, 2},
				Name:   []string{"jack", "garland"},
				Length: []float64{1.23},
			}},
			want:    " WHERE 1=1 AND id IN (?, ?) AND name IN (?, ?) AND length IN (?);",
			want1:   []interface{}{int64(1), int64(2), "jack", "garland", float64(1.23)},
			want2:   " WHERE 1=1 AND id IN (?, ?) AND name IN (?, ?) AND length IN (?);",
			want3:   []interface{}{int64(1), int64(2), "jack", "garland", float64(1.23)},
			wantErr: false,
		},
		{
			name: "test empty array with IN",
			args: args{param: &TestParamWithArray{
				ID:     []int64{},
				Name:   []string{},
				Length: []float64{},
			}},
			want:    " WHERE 1=1;",
			want1:   nil,
			want2:   " WHERE 1=1;",
			want3:   nil,
			wantErr: false,
		},
		{
			name: "test string LIKE",
			args: args{param: &TestParamWithStringWildcard{
				Name: "garland%",
			}},
			want:    " WHERE 1=1 AND name LIKE ?;",
			want1:   []interface{}{"garland%"},
			want2:   " WHERE 1=1 AND name LIKE ?;",
			want3:   []interface{}{"garland%"},
			wantErr: false,
		},
		{
			name: "test string LIKE",
			args: args{param: &TestParamWithStringWildcard{
				NameOpt: "garland%",
			}},
			want:    " WHERE 1=1 OR name LIKE ?;",
			want1:   []interface{}{"garland%"},
			want2:   " WHERE 1=1 OR name LIKE ?;",
			want3:   []interface{}{"garland%"},
			wantErr: false,
		},
		{
			name: "test limit and page",
			args: args{param: &TestParamLimitAndPage{
				ID:    1,
				Names: []string{"jack", "garland"},
				Limit: 10,
				Page:  1,
			}},
			want:    " WHERE 1=1 AND id=? AND name IN (?, ?) LIMIT 0, 10;",
			want1:   []interface{}{int64(1), "jack", "garland"},
			want2:   " WHERE 1=1 AND id=? AND name IN (?, ?);",
			want3:   []interface{}{int64(1), "jack", "garland"},
			wantErr: false,
		},
		{
			name: "test sort by",
			args: args{param: &TestParamSortBy{
				ID:     1,
				SortBy: []string{"id"},
				Limit:  10,
				Page:   1,
			}},
			want:    " WHERE 1=1 AND id=? ORDER BY id ASC LIMIT 0, 10;",
			want1:   []interface{}{int64(1)},
			want2:   " WHERE 1=1 AND id=?;",
			want3:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name: "test null type",
			args: args{param: &TestNullType{
				ID:        null.Int64{Valid: true},
				Name:      null.String{Valid: true},
				Length:    null.Float64{Valid: true},
				Active:    null.Bool{Valid: true},
				Birthday:  null.Date{Valid: true},
				CreatedAt: null.Time{Valid: true},
			}},
			want:    " WHERE 1=1 AND id=? AND name=? AND length=? AND active=? AND birthday=? AND created_at=?;",
			want1:   []interface{}{int64(0), "", float64(0), false, time.Time{}, time.Time{}},
			want2:   " WHERE 1=1 AND id=? AND name=? AND length=? AND active=? AND birthday=? AND created_at=?;",
			want3:   []interface{}{int64(0), "", float64(0), false, time.Time{}, time.Time{}},
			wantErr: false,
		},
		{
			name: "test null type invalid",
			args: args{param: &TestNullType{
				ID:        null.Int64{},
				Name:      null.String{},
				Length:    null.Float64{},
				Active:    null.Bool{},
				Birthday:  null.Date{},
				CreatedAt: null.Time{},
			}},
			want:    " WHERE 1=1;",
			want1:   nil,
			want2:   " WHERE 1=1;",
			want3:   nil,
			wantErr: false,
		},
		{
			name: "test complex operation",
			args: args{param: &TestComplexOperation{
				NameNIN: []string{"jack", "garland"},
				AgeGTE:  11,
				AgeGT:   10,
				AgeLTE:  9,
				AgeLT:   10,
				AgeNE:   100,
				JobOPT:  "knight",
			}},
			want:    " WHERE 1=1 AND name NOT IN (?, ?) AND age>=? AND age>? AND age<=? AND age<? AND age<>? OR job=?;",
			want1:   []interface{}{"jack", "garland", int64(11), int64(10), int64(9), int64(10), int64(100), "knight"},
			want2:   " WHERE 1=1 AND name NOT IN (?, ?) AND age>=? AND age>? AND age<=? AND age<? AND age<>? OR job=?;",
			want3:   []interface{}{"jack", "garland", int64(11), int64(10), int64(9), int64(10), int64(100), "knight"},
			wantErr: false,
		},
		{
			name: "test all kinds of valid booleans",
			args: args{param: &TestAllKindsOfBooleans{
				Boolean1: true,
				Boolean2: []bool{true},
				Boolean3: []*bool{&BooleanTrue},
				Boolean4: null.BoolFrom(true),
				Boolean5: []null.Bool{null.BoolFrom(true)},
				Boolean6: []*null.Bool{{Valid: true, Bool: true}},
			}},
			want:    " WHERE 1=1 AND boolean_1=? AND boolean_2 IN (?) AND boolean_3 IN (?) AND boolean_4=? AND boolean_5 IN (?) AND boolean_6 IN (?);",
			want1:   []interface{}{true, true, &BooleanTrue, true, true, true},
			want2:   " WHERE 1=1 AND boolean_1=? AND boolean_2 IN (?) AND boolean_3 IN (?) AND boolean_4=? AND boolean_5 IN (?) AND boolean_6 IN (?);",
			want3:   []interface{}{true, true, &BooleanTrue, true, true, true},
			wantErr: false,
		},
		{
			name: "test option disable pagination limit",
			args: args{param: &TestParamLimitAndPage{
				ID:    1,
				Names: []string{"jack", "garland"},
				Limit: 10,
				Page:  1,
			}, opt: &Option{
				DisableLimit: true,
			}},
			want:    " WHERE 1=1 AND id=? AND name IN (?, ?);",
			want1:   []interface{}{int64(1), "jack", "garland"},
			want2:   " WHERE 1=1 AND id=? AND name IN (?, ?);",
			want3:   []interface{}{int64(1), "jack", "garland"},
			wantErr: false,
		},
		{
			name: "test option is active",
			args: args{param: &TestParamLimitAndPage{
				ID:    1,
				Names: []string{"jack", "garland"},
				Limit: 10,
				Page:  1,
			}, opt: &Option{
				IsActive: true,
			}},
			want:    " WHERE 1=1 AND status=1 AND id=? AND name IN (?, ?) LIMIT 0, 10;",
			want1:   []interface{}{int64(1), "jack", "garland"},
			want2:   " WHERE 1=1 AND status=1 AND id=? AND name IN (?, ?);",
			want3:   []interface{}{int64(1), "jack", "garland"},
			wantErr: false,
		},
		{
			name: "test option is inactive",
			args: args{param: &TestParamLimitAndPage{
				ID:    1,
				Names: []string{"jack", "garland"},
				Limit: 10,
				Page:  1,
			}, opt: &Option{
				IsInactive: true,
			}},
			want:    " WHERE 1=1 AND status=0 AND id=? AND name IN (?, ?) LIMIT 0, 10;",
			want1:   []interface{}{int64(1), "jack", "garland"},
			want2:   " WHERE 1=1 AND status=0 AND id=? AND name IN (?, ?);",
			want3:   []interface{}{int64(1), "jack", "garland"},
			wantErr: false,
		},
		{
			name:    "test param is nil",
			args:    args{param: nil},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := initTestDatabase()

			qBuilder := NewSQLQueryBuilder(db, "param", "db", tt.args.opt)

			got, got1, got2, got3, err := qBuilder.Build(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("sqlClausebuilder.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)
			assert.Equal(t, tt.want3, got3)
		})
	}
}

func Test_sqlClausebuilder_AddGroupByQuery(t *testing.T) {
	type args struct {
		param   interface{}
		columns []string
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
			name: "test group by query",
			args: args{param: &TestParamLimitAndPage{
				ID:    1,
				Names: []string{"jack", "garland"},
				Limit: 10,
				Page:  1,
			}, columns: []string{"id"}},
			want:    " WHERE 1=1 AND id=? AND name IN (?, ?) GROUP BY id LIMIT 0, 10;",
			want1:   []interface{}{int64(1), "jack", "garland"},
			want2:   " WHERE 1=1 AND id=? AND name IN (?, ?);",
			want3:   []interface{}{int64(1), "jack", "garland"},
			wantErr: false,
		},
		{
			name: "test group by query multiple columns",
			args: args{param: &TestParamLimitAndPage{
				ID:    1,
				Names: []string{"jack", "garland"},
				Limit: 10,
				Page:  1,
			}, columns: []string{"id", "name"}},
			want:    " WHERE 1=1 AND id=? AND name IN (?, ?) GROUP BY id, name LIMIT 0, 10;",
			want1:   []interface{}{int64(1), "jack", "garland"},
			want2:   " WHERE 1=1 AND id=? AND name IN (?, ?);",
			want3:   []interface{}{int64(1), "jack", "garland"},
			wantErr: false,
		},
		{
			name: "test group by query empty",
			args: args{param: &TestParamLimitAndPage{
				ID:    1,
				Names: []string{"jack", "garland"},
				Limit: 10,
				Page:  1,
			}, columns: []string{}},
			want:    " WHERE 1=1 AND id=? AND name IN (?, ?) LIMIT 0, 10;",
			want1:   []interface{}{int64(1), "jack", "garland"},
			want2:   " WHERE 1=1 AND id=? AND name IN (?, ?);",
			want3:   []interface{}{int64(1), "jack", "garland"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := initTestDatabase()

			qBuilder := NewSQLQueryBuilder(db, "param", "db", nil)

			qBuilder.AddGroupByQuery(tt.args.columns...)

			got, got1, got2, got3, err := qBuilder.Build(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("sqlClausebuilder.AddGroupByQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)
			assert.Equal(t, tt.want3, got3)
		})
	}
}

func Test_sqlClausebuilder_AddPrefixQuery(t *testing.T) {
	type args struct {
		param       interface{}
		prefixquery string
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
			name: "test prefix query",
			args: args{param: &TestParamLimitAndPage{
				ID:    1,
				Names: []string{"jack", "garland"},
				Limit: 10,
				Page:  1,
			}, prefixquery: "active = 1"},
			want:    " WHERE 1=1 AND active = 1 AND id=? AND name IN (?, ?) LIMIT 0, 10;",
			want1:   []interface{}{int64(1), "jack", "garland"},
			want2:   " WHERE 1=1 AND active = 1 AND id=? AND name IN (?, ?);",
			want3:   []interface{}{int64(1), "jack", "garland"},
			wantErr: false,
		},
		{
			name: "test prefix query empty",
			args: args{param: &TestParamLimitAndPage{
				ID:    1,
				Names: []string{"jack", "garland"},
				Limit: 10,
				Page:  1,
			}, prefixquery: ""},
			want:    " WHERE 1=1 AND id=? AND name IN (?, ?) LIMIT 0, 10;",
			want1:   []interface{}{int64(1), "jack", "garland"},
			want2:   " WHERE 1=1 AND id=? AND name IN (?, ?);",
			want3:   []interface{}{int64(1), "jack", "garland"},
			wantErr: false,
		},
		{
			name: "test prefix with bindvars",
			args: args{param: &TestParamLimitAndPage{
				ID:    1,
				Names: []string{"jack", "garland"},
				Limit: 10,
				Page:  1,
			}, prefixquery: "status IN (?)"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := initTestDatabase()

			qBuilder := NewSQLQueryBuilder(db, "param", "db", nil)

			qBuilder.AddPrefixQuery(tt.args.prefixquery)

			got, got1, got2, got3, err := qBuilder.Build(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("sqlClausebuilder.AddPrefixQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)
			assert.Equal(t, tt.want3, got3)
		})
	}
}

func Test_sqlClausebuilder_AddAliasPrefix(t *testing.T) {
	type args struct {
		param       interface{}
		aliasprefix string
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
			name: "test alias prefix",
			args: args{param: &TestParamLimitAndPage{
				ID:    1,
				Names: []string{"jack", "garland"},
				Limit: 10,
				Page:  1,
			}, aliasprefix: "idk"},
			want:    " WHERE 1=1 AND idk.id=? AND idk.name IN (?, ?) LIMIT 0, 10;",
			want1:   []interface{}{int64(1), "jack", "garland"},
			want2:   " WHERE 1=1 AND idk.id=? AND idk.name IN (?, ?);",
			want3:   []interface{}{int64(1), "jack", "garland"},
			wantErr: false,
		},
		{
			name: "test alias prefix no tag",
			args: args{param: &TestParamNoTag{
				ID:   1,
				Name: "jack",
			}, aliasprefix: "idk"},
			want:    " WHERE 1=1;",
			want1:   nil,
			want2:   " WHERE 1=1;",
			want3:   nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := initTestDatabase()

			qBuilder := NewSQLQueryBuilder(db, "param", "db", nil)

			qBuilder.AddAliasPrefix(tt.args.aliasprefix, tt.args.param)

			got, got1, got2, got3, err := qBuilder.Build(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("sqlClausebuilder.AddAliasPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)
			assert.Equal(t, tt.want3, got3)
		})
	}
}

func Test_sqlClausebuilder_AddSuffixQuery(t *testing.T) {
	type args struct {
		param       interface{}
		suffixquery string
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
			name: "test suffix query",
			args: args{param: &TestParamLimitAndPage{
				ID:    1,
				Names: []string{"jack", "garland"},
				Limit: 10,
				Page:  1,
			}, suffixquery: "GROUP BY name"},
			want:    " WHERE 1=1 AND id=? AND name IN (?, ?) LIMIT 0, 10 GROUP BY name;",
			want1:   []interface{}{int64(1), "jack", "garland"},
			want2:   " WHERE 1=1 AND id=? AND name IN (?, ?) GROUP BY name;",
			want3:   []interface{}{int64(1), "jack", "garland"},
			wantErr: false,
		},
		{
			name: "test suffix query empty",
			args: args{param: &TestParamLimitAndPage{
				ID:    1,
				Names: []string{"jack", "garland"},
				Limit: 10,
				Page:  1,
			}, suffixquery: ""},
			want:    " WHERE 1=1 AND id=? AND name IN (?, ?) LIMIT 0, 10;",
			want1:   []interface{}{int64(1), "jack", "garland"},
			want2:   " WHERE 1=1 AND id=? AND name IN (?, ?);",
			want3:   []interface{}{int64(1), "jack", "garland"},
			wantErr: false,
		},
		{
			name: "test suffix with bindvars",
			args: args{param: &TestParamLimitAndPage{
				ID:    1,
				Names: []string{"jack", "garland"},
				Limit: 10,
				Page:  1,
			}, suffixquery: "status IN (?)"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := initTestDatabase()

			qBuilder := NewSQLQueryBuilder(db, "param", "db", nil)

			qBuilder.AddSuffixQuery(tt.args.suffixquery)

			got, got1, got2, got3, err := qBuilder.Build(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("sqlClausebuilder.AddSuffixQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)
			assert.Equal(t, tt.want3, got3)
		})
	}
}

func Test_sqlClausebuilder_BuildUpdate(t *testing.T) {
	mockTime := time.Now()

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
		opt    *Option
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []interface{}
		wantErr bool
	}{
		{
			name:    "error update nil",
			args:    args{update: nil, where: &testWhereStruct{ID: 1}},
			wantErr: true,
		},
		{
			name:    "error where nil",
			args:    args{update: &testUpdateStruct{FirstName: "jack"}, where: nil},
			wantErr: true,
		},
		{
			name:    "error update empty",
			args:    args{update: &testUpdateStruct{}, where: &testWhereStruct{ID: 1}},
			wantErr: true,
		},
		{
			name:    "error where empty",
			args:    args{update: &testUpdateStruct{FirstName: "jack"}, where: &testWhereStruct{}},
			wantErr: true,
		},
		{
			name: "test empty null data",
			args: args{
				update: &testUpdateStruct{
					FirstName: "jack",
					Age:       21,
				},
				where: &testWhereStruct{
					ID: 1,
				}},
			want:    " SET first_name=?, age=? WHERE 1=1 AND id=?;",
			want1:   []interface{}{"jack", int64(21), int64(1)},
			wantErr: false,
		},
		{
			name: "test empty primitive data",
			args: args{
				update: &testUpdateStruct{
					LastName: null.StringFrom(""),
					Children: null.Int64From(0),
				},
				where: &testWhereStruct{
					ID: 1,
				}},
			want:    " SET last_name=?, children=? WHERE 1=1 AND id=?;",
			want1:   []interface{}{"", int64(0), int64(1)},
			wantErr: false,
		},
		{
			name: "test fill all data",
			args: args{
				update: &testUpdateStruct{
					FirstName: "jack",
					LastName:  null.StringFrom("garland"),
					Age:       21,
					Children:  null.Int64From(1),
				},
				where: &testWhereStruct{
					ID: 1,
				}},
			want:    " SET first_name=?, last_name=?, age=?, children=? WHERE 1=1 AND id=?;",
			want1:   []interface{}{"jack", "garland", int64(21), int64(1), int64(1)},
			wantErr: false,
		},
		{
			name: "test fill pagination data",
			args: args{
				update: &testUpdateStruct{
					FirstName: "jack",
					LastName:  null.StringFrom("garland"),
					Age:       21,
					Children:  null.Int64From(1),
				},
				where: &testWhereStruct{
					ID:     1,
					Limit:  5,
					Page:   1,
					SortBy: "name",
				}},
			want:    " SET first_name=?, last_name=?, age=?, children=? WHERE 1=1 AND id=?;",
			want1:   []interface{}{"jack", "garland", int64(21), int64(1), int64(1)},
			wantErr: false,
		},
		{
			name: "test nulltype data",
			args: args{
				update: &TestNullType{
					ID:        null.Int64From(1),
					Name:      null.StringFrom("jack"),
					Length:    null.Float64From(1.23),
					Active:    null.BoolFrom(true),
					Birthday:  null.DateFrom(mockTime),
					CreatedAt: null.TimeFrom(mockTime),
				},
				where: &testWhereStruct{
					ID: 1,
				},
			},
			want:    " SET id=?, name=?, length=?, active=?, birthday=?, created_at=? WHERE 1=1 AND id=?;",
			want1:   []interface{}{int64(1), "jack", float64(1.23), true, mockTime, mockTime, int64(1)},
			wantErr: false,
		},
		{
			name: "test nulltype data zero valid",
			args: args{
				update: &TestNullType{
					ID:        null.Int64From(0),
					Name:      null.StringFrom(""),
					Length:    null.Float64From(0),
					Active:    null.BoolFrom(false),
					Birthday:  null.DateFrom(time.Time{}),
					CreatedAt: null.TimeFrom(time.Time{}),
				},
				where: &testWhereStruct{
					ID: 1,
				},
			},
			want:    " SET id=?, name=?, length=?, active=?, birthday=?, created_at=? WHERE 1=1 AND id=?;",
			want1:   []interface{}{int64(0), "", float64(0), false, time.Time{}, time.Time{}, int64(1)},
			wantErr: false,
		},
		{
			name: "test nulltype date and time update to null",
			args: args{
				update: &TestNullType{
					Birthday:  null.Date{SqlNull: true},
					CreatedAt: null.Time{SqlNull: true},
				},
				where: &testWhereStruct{
					ID: 1,
				},
			},
			want:    " SET birthday=NULL, created_at=NULL WHERE 1=1 AND id=?;",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
		{
			name: "test nulltype string update to null",
			args: args{
				update: &TestNullType{
					Birthday:  null.Date{SqlNull: true},
					CreatedAt: null.Time{SqlNull: true},
					Name:      null.String{SqlNull: true},
				},
				where: &testWhereStruct{
					ID: 1,
				},
			},
			want:    " SET birthday=NULL, created_at=NULL, name=NULL WHERE 1=1 AND id=?;",
			want1:   []interface{}{int64(1)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := initTestDatabase()

			qBuilder := NewSQLQueryBuilder(db, "param", "db", tt.args.opt)

			got, got1, err := qBuilder.BuildUpdate(tt.args.update, tt.args.where)
			if (err != nil) != tt.wantErr {
				t.Errorf("sqlClausebuilder.BuildUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
