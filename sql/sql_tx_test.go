//go:build integration
// +build integration

package sql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestCommit(t *testing.T) {
	mockMysqlDb := initTestMysqlDatabase()
	mockPostgresDb := initTestPostgresDatabase()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "ok",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, _ := mockMysqlDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
			err := tx.Commit()
			if (err != nil) != tt.wantErr {
				t.Errorf("Commit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			txPg, _ := mockPostgresDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
			errPg := txPg.Commit()
			if (errPg != nil) != tt.wantErr {
				t.Errorf("Commit() error = %v, wantErr %v", errPg, tt.wantErr)
				return
			}
		})
	}
}

func TestRollback(t *testing.T) {
	mockMysqlDb := initTestMysqlDatabase()
	mockPostgresDb := initTestPostgresDatabase()

	tests := []struct {
		name       string
		prepTxMock func() CommandTx
		wantErr    bool
	}{
		{
			name: "ok mysql",
			prepTxMock: func() CommandTx {
				tx, _ := mockMysqlDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
				return tx
			},
			wantErr: false,
		},
		{
			name: "ok postgres",
			prepTxMock: func() CommandTx {
				tx, _ := mockPostgresDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
				return tx
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := tt.prepTxMock()
			tx.Rollback()
		})
	}
}

func TestRebind(t *testing.T) {
	mockMysqlDb := initTestMysqlDatabase()
	mockPostgresDb := initTestPostgresDatabase()

	type args struct {
		query string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				query: "SELECT * FROM test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, _ := mockMysqlDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
			tx.Rebind(tt.args.query)

			txPg, _ := mockPostgresDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
			txPg.Rebind(tt.args.query)
		})
	}
}

func TestSelect(t *testing.T) {
	mockMysqlDb := initTestMysqlDatabase()
	mockPostgresDb := initTestPostgresDatabase()

	type args struct {
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
			name:    "error",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, _ := mockMysqlDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
			err := tx.Select(tt.args.name, tt.args.query, tt.args.dest, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Select() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			txPg, _ := mockPostgresDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
			errPg := txPg.Select(tt.args.name, tt.args.query, tt.args.dest, tt.args.args)
			if (errPg != nil) != tt.wantErr {
				t.Errorf("Select() error = %v, wantErr %v", errPg, tt.wantErr)
				return
			}
		})
	}
}

func TestGet(t *testing.T) {
	mockMysqlDb := initTestMysqlDatabase()
	mockPostgresDb := initTestPostgresDatabase()

	type args struct {
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
			name:    "error",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, _ := mockMysqlDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
			err := tx.Get(tt.args.name, tt.args.query, tt.args.dest, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			txPg, _ := mockPostgresDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
			errPg := txPg.Get(tt.args.name, tt.args.query, tt.args.dest, tt.args.args)
			if (errPg != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", errPg, tt.wantErr)
				return
			}
		})
	}
}

func TestQueryRow(t *testing.T) {
	mockMysqlDb := initTestMysqlDatabase()
	mockPostgresDb := initTestPostgresDatabase()

	type args struct {
		name  string
		query string
		args  interface{}
	}

	tests := []struct {
		name    string
		driver  string
		args    args
		want    *sqlx.Row
		wantErr bool
	}{
		{
			name:   "ok mysql",
			driver: "mysql",
			args: args{
				name:  "testing",
				query: "SELECT * FROM farm WHERE id = ?",
				args:  int64(1),
			},
			wantErr: false,
		},
		{
			name:   "ok postgres",
			driver: "postgres",
			args: args{
				name:  "testing",
				query: "SELECT * FROM farm WHERE id = $1",
				args:  int64(1),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.driver {
			case "mysql":
				tx, _ := mockMysqlDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
				_, err := tx.QueryRow(tt.args.name, tt.args.query, tt.args.args)
				if (err != nil) != tt.wantErr {
					t.Errorf("QueryRow() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			case "postgres":
				txPg, _ := mockPostgresDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
				_, errPg := txPg.QueryRow(tt.args.name, tt.args.query, tt.args.args)
				if (errPg != nil) != tt.wantErr {
					t.Errorf("QueryRow() error = %v, wantErr %v", errPg, tt.wantErr)
					return
				}
			}

		})
	}
}

func TestQuery(t *testing.T) {
	mockMysqlDb := initTestMysqlDatabase()
	mockPostgresDb := initTestPostgresDatabase()

	type args struct {
		name  string
		query string
		args  interface{}
	}

	tests := []struct {
		name    string
		driver  string
		args    args
		want    *sqlx.Row
		wantErr bool
	}{
		{
			name:   "ok",
			driver: "mysql",
			args: args{
				name:  "testing",
				query: "SELECT * FROM farm WHERE id = ?",
				args:  int64(1),
			},
			wantErr: false,
		},
		{
			name:   "ok",
			driver: "postgres",
			args: args{
				name:  "testing",
				query: "SELECT * FROM farm WHERE id = $1",
				args:  int64(1),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.driver {
			case "mysql":
				tx, _ := mockMysqlDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
				_, err := tx.Query(tt.args.name, tt.args.query, tt.args.args)
				if (err != nil) != tt.wantErr {
					t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			case "postgres":
				txPg, _ := mockPostgresDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
				_, errPg := txPg.Query(tt.args.name, tt.args.query, tt.args.args)
				if (errPg != nil) != tt.wantErr {
					t.Errorf("Query() error = %v, wantErr %v", errPg, tt.wantErr)
					return
				}
			}

		})
	}
}

func TestPrepare(t *testing.T) {
	mockMysqlDb := initTestMysqlDatabase()
	mockPostgresDb := initTestPostgresDatabase()

	type args struct {
		name  string
		query string
	}

	tests := []struct {
		name    string
		driver  string
		args    args
		want    CommandStmt
		wantErr bool
	}{
		{
			name:   "ok",
			driver: "mysql",
			args: args{
				name:  "testing",
				query: "SELECT * FROM farm WHERE id = ?",
			},
			wantErr: false,
		},
		{
			name:   "ok",
			driver: "postgres",
			args: args{
				name:  "testing",
				query: "SELECT * FROM farm WHERE id = $1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.driver {
			case "mysql":
				tx, _ := mockMysqlDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
				_, err := tx.Prepare(tt.args.name, tt.args.query)
				if (err != nil) != tt.wantErr {
					t.Errorf("Prepare() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			case "postgress":
				txPg, _ := mockPostgresDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
				_, errPg := txPg.Prepare(tt.args.name, tt.args.query)
				if (errPg != nil) != tt.wantErr {
					t.Errorf("Prepare() error = %v, wantErr %v", errPg, tt.wantErr)
					return
				}
			}
		})
	}
}

func TestNamedExecTx(t *testing.T) {
	mockMysqlDb := initTestMysqlDatabase()
	mockPostgresDb := initTestPostgresDatabase()

	type args struct {
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
				name:  "testing",
				query: "UPDATE farm SET name = :name WHERE id = :id",
				args:  map[string]interface{}{"name": "new name", "id": 1},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, _ := mockMysqlDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
			_, err := tx.NamedExec(tt.args.name, tt.args.query, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NamedExec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tx.Rollback()

			txPg, _ := mockPostgresDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
			_, errPg := txPg.NamedExec(tt.args.name, tt.args.query, tt.args.args)
			if (errPg != nil) != tt.wantErr {
				t.Errorf("NamedExec() error = %v, wantErr %v", errPg, tt.wantErr)
				return
			}
			txPg.Rollback()
		})
	}
}

func TestExec(t *testing.T) {
	mockMysqlDb := initTestMysqlDatabase()
	mockPostgresDb := initTestPostgresDatabase()

	type args struct {
		name  string
		query string
		args  interface{}
	}

	tests := []struct {
		name    string
		driver  string
		args    args
		want    sql.Result
		wantErr bool
	}{
		{
			name:   "ok",
			driver: "mysql",
			args: args{
				name:  "testing",
				query: "INSERT INTO farm (fk_company_id) VALUES (?)",
				args:  int64(1),
			},
			wantErr: false,
		},
		{
			name:   "ok",
			driver: "postgres",
			args: args{
				name:  "testing",
				query: "INSERT INTO farm (id) VALUES ($1)",
				args:  int64(1),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.driver {
			case "mysql":
				tx, _ := mockMysqlDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
				_, err := tx.Exec(tt.args.name, tt.args.query, tt.args.args)
				if (err != nil) != tt.wantErr {
					t.Errorf("Exec() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				tx.Rollback()
			case "postgres":
				txPg, _ := mockPostgresDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
				_, errPg := txPg.Exec(tt.args.name, tt.args.query, tt.args.args)
				if (errPg != nil) != tt.wantErr {
					t.Errorf("Exec() error = %v, wantErr %v", errPg, tt.wantErr)
					return
				}
				txPg.Rollback()
			}
		})
	}
}

func TestStmt(t *testing.T) {
	mockMysqlDb := initTestMysqlDatabase()
	mockPostgresDb := initTestPostgresDatabase()

	type args struct {
		name string
		stmt *sqlx.Stmt
	}

	tests := []struct {
		name string
		args args
		want CommandStmt
	}{
		{
			name: "ok",
			args: args{
				name: "testing",
				stmt: &sqlx.Stmt{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, _ := mockMysqlDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
			tx.Stmt(tt.args.name, tt.args.stmt)
			tx.Rollback()

			txPg, _ := mockPostgresDb.Leader().BeginTx(context.Background(), "test", TxOptions{})
			txPg.Stmt(tt.args.name, tt.args.stmt)
			txPg.Rollback()
		})
	}
}
