//go:build integration
// +build integration

package sql

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestClose(t *testing.T) {
	mockDb := initTestDatabase()

	tests := []struct {
		name string
	}{
		{
			name: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb.Leader().Close()
		})
	}
}

func TestPing(t *testing.T) {
	mockDb := initTestDatabase()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb.Leader().Ping(tt.args.ctx)
		})
	}
}

func TestIn(t *testing.T) {
	mockDb := initTestDatabase()

	type args struct {
		query string
		args  interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ok",
			args: args{
				query: "SELECT * FROM farm WHERE ID IN (?)",
				args:  int64(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb.Leader().In(tt.args.query, tt.args.args)
		})
	}
}

func TestRebindCmd(t *testing.T) {
	mockDb := initTestDatabase()

	type args struct {
		query string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				query: "SELECT * FROM farm WHERE ID IN (?)",
			},
			want: "SELECT * FROM farm WHERE ID IN (?)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mockDb.Leader().Rebind(tt.args.query)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestQueryIn(t *testing.T) {
	mockDb := initTestDatabase()

	type args struct {
		ctx   context.Context
		name  string
		query string
		args  interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *sqlx.Rows
		wantErr bool
	}{
		{
			name: "error",
			args: args{
				ctx:   context.Background(),
				name:  "testing in",
				query: "SELECT * FROM farm WHERE ID IN (?, ? ?)",
				args:  []interface{}{},
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				ctx:   context.Background(),
				name:  "testing in",
				query: "SELECT * FROM farm WHERE ID IN (?)",
				args:  int64(1),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := mockDb.Leader().QueryIn(tt.args.ctx, tt.args.name, tt.args.query, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestQueryRowCmd(t *testing.T) {
	mockDb := initTestDatabase()

	type args struct {
		ctx   context.Context
		name  string
		query string
		args  interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *sqlx.Row
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ctx:   context.Background(),
				name:  "testing in",
				query: "SELECT * FROM farm WHERE ID IN (?)",
				args:  int64(1),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := mockDb.Leader().QueryRow(tt.args.ctx, tt.args.name, tt.args.query, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryRow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestQueryCmd(t *testing.T) {
	mockDb := initTestDatabase()

	type args struct {
		ctx   context.Context
		name  string
		query string
		args  interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *sqlx.Rows
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ctx:   context.Background(),
				name:  "testing in",
				query: "SELECT * FROM farm WHERE ID IN (?)",
				args:  int64(1),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := mockDb.Leader().Query(tt.args.ctx, tt.args.name, tt.args.query, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNamedQuery(t *testing.T) {
	mockDb := initTestDatabase()

	type args struct {
		ctx   context.Context
		name  string
		query string
		arg   interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *sqlx.Rows
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ctx:   context.Background(),
				name:  "testing in",
				query: "SELECT * FROM farm WHERE id=:id",
				arg:   map[string]interface{}{"id": 1},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := mockDb.Leader().NamedQuery(tt.args.ctx, tt.args.name, tt.args.query, tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NamedQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPrepareCmd(t *testing.T) {
	mockDb := initTestDatabase()
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()

	type args struct {
		ctx   context.Context
		name  string
		query string
	}
	tests := []struct {
		name    string
		args    args
		want    CommandStmt
		wantErr bool
	}{
		{
			name: "error",
			args: args{
				ctx:   ctxTimeout,
				name:  "testing in",
				query: "SELECT * FROM farm WHERE id=?",
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				ctx:   context.Background(),
				name:  "testing in",
				query: "SELECT * FROM farm WHERE id=?",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := mockDb.Leader().Prepare(tt.args.ctx, tt.args.name, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Prepare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNamedExec(t *testing.T) {
	mockDb := initTestDatabase()

	type args struct {
		ctx   context.Context
		name  string
		query string
		args  interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    sql.Result
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ctx:   context.Background(),
				name:  "testing in",
				query: "UPDATE farm SET name = :name WHERE id = :id",
				args:  map[string]interface{}{"name": "testing farm", "id": 1},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := mockDb.Leader().NamedExec(tt.args.ctx, tt.args.name, tt.args.query, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NamedExec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestExecCmd(t *testing.T) {
	mockDb := initTestDatabase()

	type args struct {
		ctx   context.Context
		name  string
		query string
		args  interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    sql.Result
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ctx:   context.Background(),
				name:  "testing in",
				query: "UPDATE farm SET name = ? WHERE id = 1",
				args:  "new farm",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := mockDb.Leader().Exec(tt.args.ctx, tt.args.name, tt.args.query, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestBeginTx(t *testing.T) {
	mockDb := initTestDatabase()
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()

	type args struct {
		ctx  context.Context
		name string
		opt  TxOptions
	}
	tests := []struct {
		name    string
		args    args
		want    CommandTx
		wantErr bool
	}{
		{
			name: "error",
			args: args{
				ctx:  ctxTimeout,
				name: "testing in",
				opt:  TxOptions{},
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				ctx:  context.Background(),
				name: "testing in",
				opt:  TxOptions{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := mockDb.Leader().BeginTx(tt.args.ctx, tt.args.name, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("BeginTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetCmd(t *testing.T) {
	mockDb := initTestDatabase()
	var resCount int64

	type args struct {
		ctx   context.Context
		name  string
		query string
		dest  interface{}
		args  interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ctx:   context.Background(),
				name:  "testing in",
				query: "SELECT COUNT(*) FROM action WHERE id = ?",
				dest:  &resCount,
				args:  int64(1),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mockDb.Leader().Get(tt.args.ctx, tt.args.name, tt.args.query, tt.args.dest, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
