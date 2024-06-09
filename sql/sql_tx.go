package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/downsized-devs/sdk-go/instrument"
	"github.com/downsized-devs/sdk-go/logger"
	"github.com/jmoiron/sqlx"
)

type CommandTx interface {
	Commit() error
	Rollback()
	Rebind(query string) string
	Select(name string, query string, dest interface{}, args ...interface{}) error
	Get(name string, query string, dest interface{}, args ...interface{}) error
	QueryRow(name string, query string, args ...interface{}) (*sqlx.Row, error)
	Query(name string, query string, args ...interface{}) (*sqlx.Rows, error)
	Prepare(name string, query string) (CommandStmt, error)

	NamedExec(name string, query string, args interface{}) (sql.Result, error)
	Exec(name string, query string, args ...interface{}) (sql.Result, error)
	Stmt(name string, stmt *sqlx.Stmt) CommandStmt
}

type commandTx struct {
	ctx           context.Context
	name          string
	connName      string
	connType      string
	tx            *sqlx.Tx
	log           logger.Interface
	instrument    instrument.Interface
	useInstrument bool
	logQuery      bool
}

func initTx(ctx context.Context, name, connName string, tx *sqlx.Tx, opts *sql.TxOptions, log logger.Interface, instr instrument.Interface, isLeader, useInstr, logQuery bool) CommandTx {
	c := &commandTx{
		ctx:           ctx,
		name:          name,
		connName:      connName,
		connType:      connTypeLeader,
		tx:            tx,
		log:           log,
		instrument:    instr,
		useInstrument: useInstr,
		logQuery:      logQuery,
	}

	if !isLeader {
		c.connType = connTypeFollower
	}

	return c
}

func (x *commandTx) Commit() error {
	return x.tx.Commit()
}

// Rollback needs to be called with defer right after calling BeginTx.
// Read here: https://go.dev/doc/database/execute-transactions.
func (x *commandTx) Rollback() {
	if err := x.tx.Rollback(); err != nil && err != sql.ErrTxDone {
		x.log.Error(x.ctx, err)
	}
}

func (x *commandTx) Rebind(query string) string {
	return x.tx.Rebind(query)
}

func (x *commandTx) Select(name string, query string, dest interface{}, args ...interface{}) error {
	if x.useInstrument {
		timer := x.instrument.DatabaseQueryTimer(x.connName, x.connType, name)
		defer timer.ObserveDuration()
	}
	if x.logQuery {
		x.log.Info(x.ctx, fmt.Sprintf(queryLogMessage, name, replaceBindvarsWithArgs(query, args...)))
	}
	return x.tx.SelectContext(x.ctx, dest, query, args...)
}

func (x *commandTx) Get(name string, query string, dest interface{}, args ...interface{}) error {
	if x.useInstrument {
		timer := x.instrument.DatabaseQueryTimer(x.connName, x.connType, name)
		defer timer.ObserveDuration()
	}
	if x.logQuery {
		x.log.Info(x.ctx, fmt.Sprintf(queryLogMessage, name, replaceBindvarsWithArgs(query, args...)))
	}
	return x.tx.GetContext(x.ctx, dest, query, args...)
}

func (x *commandTx) QueryRow(name string, query string, args ...interface{}) (*sqlx.Row, error) {
	if x.useInstrument {
		timer := x.instrument.DatabaseQueryTimer(x.connName, x.connType, name)
		defer timer.ObserveDuration()
	}
	if x.logQuery {
		x.log.Info(x.ctx, fmt.Sprintf(queryLogMessage, name, replaceBindvarsWithArgs(query, args...)))
	}
	row := x.tx.QueryRowxContext(x.ctx, query, args...)
	return row, row.Err()
}

func (x *commandTx) Query(name string, query string, args ...interface{}) (*sqlx.Rows, error) {
	if x.useInstrument {
		timer := x.instrument.DatabaseQueryTimer(x.connName, x.connType, name)
		defer timer.ObserveDuration()
	}
	if x.logQuery {
		x.log.Info(x.ctx, fmt.Sprintf(queryLogMessage, name, replaceBindvarsWithArgs(query, args...)))
	}
	return x.tx.QueryxContext(x.ctx, query, args...)
}

func (x *commandTx) Prepare(name string, query string) (CommandStmt, error) {
	if x.useInstrument {
		timer := x.instrument.DatabaseQueryTimer(x.connName, x.connType, name)
		defer timer.ObserveDuration()
	}
	stmt, err := x.tx.PreparexContext(x.ctx, query)
	if err != nil {
		return nil, err
	}
	return initStmt(x.ctx, name, x.connName, stmt, x.instrument, x.connType == connTypeLeader, x.useInstrument), nil
}

func (x *commandTx) NamedExec(name string, query string, args interface{}) (sql.Result, error) {
	if x.useInstrument {
		timer := x.instrument.DatabaseQueryTimer(x.connName, x.connType, name)
		defer timer.ObserveDuration()
	}
	if x.logQuery {
		x.log.Info(x.ctx, fmt.Sprintf(queryLogMessage, name, replaceBindvarsWithArgs(query)))
	}
	return x.tx.NamedExecContext(x.ctx, query, args)
}

func (x *commandTx) Exec(name string, query string, args ...interface{}) (sql.Result, error) {
	if x.useInstrument {
		timer := x.instrument.DatabaseQueryTimer(x.connName, x.connType, name)
		defer timer.ObserveDuration()
	}
	if x.logQuery {
		x.log.Info(x.ctx, fmt.Sprintf(queryLogMessage, name, replaceBindvarsWithArgs(query, args...)))
	}
	return x.tx.ExecContext(x.ctx, query, args...)
}

func (x *commandTx) Stmt(name string, stmt *sqlx.Stmt) CommandStmt {
	if x.useInstrument {
		timer := x.instrument.DatabaseQueryTimer(x.connName, x.connType, name)
		defer timer.ObserveDuration()
	}
	return initStmt(x.ctx, name, x.connName, stmt, x.instrument, x.connType == connTypeLeader, x.useInstrument)
}
