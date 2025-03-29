package query

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/downsized-devs/sdk-go/codes"
	"github.com/downsized-devs/sdk-go/errors"
	"github.com/downsized-devs/sdk-go/sql"
	"github.com/jmoiron/sqlx"
)

const cursorField = "cursorField"

type Cursor interface {
	DecodeCursor(v string) error
	EncodeCursor() (string, error)
}

type Option struct {
	DisableLimit bool `form:"disableLimit"`
	IsActive     bool
	IsInactive   bool
}

type sqlClausebuilder struct {
	rawQuery                      *bytes.Buffer
	suffixQuery                   string
	groupBy                       []string
	rawUpdate                     *bytes.Buffer
	param                         reflect.Value
	updateParam                   reflect.Value
	args                          []interface{}
	updateArgs                    []interface{}
	paramToDBMap, paramToFieldMap map[string]string
	sortableParams                map[string]bool
	dbTag                         string
	paramTag                      string
	fieldTag                      string
	sortClause                    string
	paginationClause              string
	paramSortBy                   []string
	dbSortBy                      []string
	limit                         int64
	page                          int64
	db                            sql.Interface
	aliasMap                      map[string]string
	disableLimit                  bool

	// cursors
	useCursor        bool
	rawCursor        string
	cursorArgCounter int
}

// Initiate new query builder object
func NewSQLQueryBuilder(db sql.Interface, paramTag, dbTag string, option *Option) *sqlClausebuilder {
	qb := sqlClausebuilder{
		db:              db,
		rawQuery:        bytes.NewBufferString(" WHERE 1=1"),
		rawUpdate:       bytes.NewBufferString(" SET"),
		args:            nil,
		updateArgs:      nil,
		fieldTag:        cursorField,
		dbTag:           dbTag,
		paramTag:        paramTag,
		paramSortBy:     nil,
		useCursor:       false,
		paramToDBMap:    make(map[string]string),
		paramToFieldMap: make(map[string]string),
		sortableParams:  make(map[string]bool),
		aliasMap:        make(map[string]string),
		limit:           0,
		page:            0,
	}

	if option != nil {
		if option.DisableLimit {
			qb.disableLimit = true
		}
		if option.IsActive {
			_, _ = qb.rawQuery.WriteString(" AND status=1")
		}
		if option.IsInactive {
			_, _ = qb.rawQuery.WriteString(" AND status=0")
		}
	}

	return &qb
}

// append custom query to the beginning of query extension
func (s *sqlClausebuilder) AddPrefixQuery(prefix string) *sqlClausebuilder {
	if len(prefix) > 0 {
		_, _ = s.rawQuery.WriteString(" AND " + prefix)
	}
	return s
}

func (s *sqlClausebuilder) AddGroupByQuery(columns ...string) *sqlClausebuilder {
	if len(columns) > 0 {
		s.groupBy = append(s.groupBy, columns...)
	}
	return s
}

// append custom query to the end of query extension
func (s *sqlClausebuilder) AddSuffixQuery(suffix string) *sqlClausebuilder {
	if len(suffix) > 0 {
		s.suffixQuery = " " + suffix
	}
	return s
}

// add prefix for db alias in query extension clause.
// example: db=id prefix=p result=p.id
func (s *sqlClausebuilder) AddAliasPrefix(alias string, ptr interface{}) *sqlClausebuilder {
	p := reflect.ValueOf(ptr)
	if p.Kind() != reflect.Ptr {
		panic(errors.NewWithCode(codes.CodeInvalidValue, "passed interface{} should be a pointer"))
	}
	v := p.Elem()
	var address string
	if v.CanAddr() {
		address = fmt.Sprint(v.Addr().Pointer())
	}

	s.aliasMap[address] = alias
	return s
}

// Build dynamic extension for sql query
func (s *sqlClausebuilder) Build(param interface{}) (string, []interface{}, string, []interface{}, error) {
	// return error if the param is not a pointer or has nil value
	p := reflect.ValueOf(param)
	if p.Kind() != reflect.Ptr || p.IsNil() {
		return "", nil, "", nil, errors.NewWithCode(codes.CodeInvalidValue, "passed param should be a pointer and cannot be nil")
	}

	// copy param to struct
	s.param = p

	traverseOnParam(s.paramTag, s.dbTag, s.fieldTag, "$", "", "", s.aliasMap, s.param, s.buildSQLQueryString)

	// copy buffer to get count query
	countquery := s.rawQuery.Bytes()

	// group by
	if len(s.groupBy) > 0 {
		s.rawQuery.WriteString(" GROUP BY " + strings.Join(s.groupBy, ", "))
	}

	// page pagination
	// TODO: implement cursor pagination
	if !s.useCursor || len(s.rawCursor) < 1 {
		// sort must be done first before page pagination
		s.sort()
		if len(s.sortClause) > 0 {
			s.rawQuery.WriteString(s.sortClause)
		}

		s.pagePagination()
		if len(s.paginationClause) > 0 && !s.disableLimit {
			s.rawQuery.WriteString(s.paginationClause)
		}
	}

	newQuery, newArgs, err := sqlx.In(s.rawQuery.String()+s.suffixQuery+";", s.args...)
	if err != nil {
		return "", nil, "", nil, err
	}
	newQuery = s.db.Leader().Rebind(newQuery)

	newCountQuery, newCountArgs, err := sqlx.In(string(countquery)+s.suffixQuery+";", s.args[0:len(s.args)-s.cursorArgCounter]...)
	if err != nil {
		return "", nil, "", nil, err
	}
	newCountQuery = s.db.Leader().Rebind(newCountQuery)

	return newQuery, newArgs, newCountQuery, newCountArgs, nil
}

// Build dynamic sql update query
func (s *sqlClausebuilder) BuildUpdate(update interface{}, where interface{}) (string, []interface{}, error) {
	// return error if the update or where param is not a pointer or has nil value
	u := reflect.ValueOf(update)
	w := reflect.ValueOf(where)

	if u.Kind() != reflect.Ptr || u.IsNil() || w.Kind() != reflect.Ptr || w.IsNil() {
		return "", nil, errors.NewWithCode(codes.CodeInvalidValue, "passed update or where param should be a pointer and cannot be nil")
	}

	// copy update and where param to struct
	s.param = w
	s.updateParam = u

	// generate update set query
	traverseOnParam(s.paramTag, s.dbTag, s.fieldTag, "$", "", "", s.aliasMap, s.updateParam, s.buildSQLUpdateString)
	// generate where query
	traverseOnParam(s.paramTag, s.dbTag, s.fieldTag, "$", "", "", s.aliasMap, s.param, s.buildSQLQueryString)

	whereQuery, whereArgs, err := sqlx.In(s.rawQuery.String()+s.suffixQuery+";", s.args...)
	if err != nil {
		return "", nil, err
	}
	whereQuery = s.db.Leader().Rebind(whereQuery)

	if strings.TrimSpace(whereQuery) == "WHERE 1=1;" || strings.TrimSpace(s.rawUpdate.String()) == "SET" {
		return "", nil, errors.NewWithCode(codes.CodeInvalidValue, "generated update or where query clause cannot be ampty")
	}

	// combine all query
	allQuery := s.rawUpdate.String() + whereQuery

	// combine all args
	var allArgs []interface{}
	allArgs = append(allArgs, s.updateArgs...)
	allArgs = append(allArgs, whereArgs...)

	return allQuery, allArgs, nil
}

func (s *sqlClausebuilder) sort() {
	for _, param := range s.paramSortBy {
		reg := regexp.MustCompile(`(?P<sign>-)?(?P<col>[a-zA-Z_\.0-9]+),?`)
		if reg.MatchString(param) {
			for _, _s := range strings.Split(param, ",") {
				direction := "ASC"
				match := reg.FindStringSubmatch(_s)
				for i, name := range reg.SubexpNames() {
					if i == 0 || name == "" {
						continue
					}
					if match != nil {
						if name == "sign" && match[i] == "-" {
							direction = "DESC"
						} else if name == "col" {
							if s.useCursor && len(s.rawCursor) > 0 && !s.sortableParams[match[i]] {
								continue
							}
							if db, ok := s.paramToDBMap[match[i]]; ok {
								db = db + " " + direction
								s.dbSortBy = append(s.dbSortBy, db)
							}
						}
					}
				}
			}
		}
	}
	if len(s.dbSortBy) > 0 {
		s.sortClause = " ORDER BY " + strings.Join(s.dbSortBy, ", ")
	}
}

func (s *sqlClausebuilder) pagePagination() {
	if s.page > 0 || s.limit > 0 {
		offset := getOffset(s.page, s.limit)
		s.paginationClause = fmt.Sprintf(" LIMIT %d, %d", offset, s.limit)
	}
}

func (s *sqlClausebuilder) buildSQLQueryString(primitiveType int8, isLike, isMany, isSqlNull bool, fieldName, paramTag, dbTag string, args interface{}) {
	// map param to field name
	s.paramToFieldMap[paramTag] = fieldName
	// map param to db column name
	s.paramToDBMap[paramTag] = dbTag

	if dbTag == "" {
		return
	}

	if isSortBy(paramTag) {
		v, _ := args.([]string)
		if v != nil {
			s.paramSortBy = normalizeSortBy(v)
		}
		return
	}

	if isPage(paramTag) {
		v, _ := args.(int64)
		s.page = validatePage(v)
		return
	}

	if isLimit(paramTag) {
		v, _ := args.(int64)
		s.limit = validateLimit(v)
		return
	}

	// we only remap if the args is not nil
	if args == nil {
		return
	}

	if !isMany {
		if isLike {
			if strings.Contains(paramTag, "__opt") {
				_, _ = s.rawQuery.WriteString(" OR " + dbTag + " LIKE " + s.getBindVar())
				s.args = append(s.args, args)
				return
			}
			_, _ = s.rawQuery.WriteString(" AND " + dbTag + " LIKE " + s.getBindVar())
			s.args = append(s.args, args)
			return
		}

		switch {
		case strings.Contains(paramTag, "__gte"):
			_, _ = s.rawQuery.WriteString(" AND " + dbTag + ">=" + s.getBindVar())
			s.args = append(s.args, args)
			return
		case strings.Contains(paramTag, "__lte"):
			_, _ = s.rawQuery.WriteString(" AND " + dbTag + "<=" + s.getBindVar())
			s.args = append(s.args, args)
			return
		case strings.Contains(paramTag, "__lt"):
			_, _ = s.rawQuery.WriteString(" AND " + dbTag + "<" + s.getBindVar())
			s.args = append(s.args, args)
			return
		case strings.Contains(paramTag, "__gt"):
			_, _ = s.rawQuery.WriteString(" AND " + dbTag + ">" + s.getBindVar())
			s.args = append(s.args, args)
			return
		case strings.Contains(paramTag, "__ne"):
			_, _ = s.rawQuery.WriteString(" AND " + dbTag + "<>" + s.getBindVar())
			s.args = append(s.args, args)
			return
		case strings.Contains(paramTag, "__opt"):
			_, _ = s.rawQuery.WriteString(" OR " + dbTag + "=" + s.getBindVar())
			s.args = append(s.args, args)
			return
		}

		_, _ = s.rawQuery.WriteString(" AND " + dbTag + "=" + s.getBindVar())
		s.args = append(s.args, args)
		return
	}

	if strings.Contains(paramTag, "__nin") {
		_, _ = s.rawQuery.WriteString(" AND " + dbTag + " NOT IN (" + s.getBindVar() + ")")
		s.args = append(s.args, args)
		return
	}

	// __ in or unstated will result IN
	_, _ = s.rawQuery.WriteString(" AND " + dbTag + " IN (" + s.getBindVar() + ")")
	s.args = append(s.args, args)
}

func (s *sqlClausebuilder) buildSQLUpdateString(primitiveType int8, isLike, isMany, isSqlNull bool, fieldName, paramTag, dbTag string, args interface{}) {
	if dbTag == "" {
		return
	}

	// we only remap if the args is not nil and isSqlNull is false
	if args == nil && !isSqlNull {
		return
	}

	separator := ""
	if strings.TrimSpace(s.rawUpdate.String()) != "SET" {
		separator = ","
	}

	// only append if data is singular
	if !isMany {
		if isSqlNull {
			_, _ = s.rawUpdate.WriteString(separator + " " + dbTag + "=NULL")
			return
		}
		_, _ = s.rawUpdate.WriteString(separator + " " + dbTag + "=" + s.getBindVar())
		s.updateArgs = append(s.updateArgs, args)
		return
	}
}

func (s *sqlClausebuilder) getBindVar() string {
	return "?"
}
