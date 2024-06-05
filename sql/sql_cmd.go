package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/downsized-devs/sdk-go/instrument"
	"github.com/downsized-devs/sdk-go/log"
	"github.com/jmoiron/sqlx"
)

type Command interface {
	Close() error
	Ping(ctx context.Context) error
	In(query string, args ...interface{}) (string, []interface{}, error)
	Rebind(query string) string
	QueryIn(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRow(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Row, error)
	Query(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Rows, error)
	NamedQuery(ctx context.Context, name string, query string, arg interface{}) (*sqlx.Rows, error)
	Prepare(ctx context.Context, name string, query string) (CommandStmt, error)

	NamedExec(ctx context.Context, name string, query string, args interface{}) (sql.Result, error)
	Exec(ctx context.Context, name string, query string, args ...interface{}) (sql.Result, error)
	BeginTx(ctx context.Context, name string, opts TxOptions) (CommandTx, error)

	Get(ctx context.Context, name string, query string, dest interface{}, args ...interface{}) error
}

type TxOptions struct {
	Isolation sql.IsolationLevel
	ReadOnly  bool
}

type command struct {
	db            *sqlx.DB
	connName      string
	connType      string
	log           log.Interface
	instrument    instrument.Interface
	useInstrument bool
	logQuery      bool
}

func initCommand(db *sqlx.DB, connName string, instr instrument.Interface, log log.Interface, isLeader, useInstr, logQuery bool) Command {
	c := &command{
		db:            db,
		connName:      connName,
		connType:      connTypeLeader,
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

func (c *command) Close() error {
	return c.db.Close()
}

func (c *command) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

func (c *command) In(query string, args ...interface{}) (string, []interface{}, error) {
	return sqlx.In(query, args...)
}

func (c *command) Rebind(query string) string {
	return c.db.Rebind(query)
}

func (c *command) QueryIn(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Rows, error) {
	q, newArgs, err := sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}
	if c.useInstrument {
		timer := c.instrument.DatabaseQueryTimer(c.connName, c.connType, name)
		defer timer.ObserveDuration()
	}
	if c.logQuery {
		c.log.Info(ctx, fmt.Sprintf(queryLogMessage, name, replaceBindvarsWithArgs(q, newArgs...)))
	}
	return c.Query(ctx, name, c.Rebind(q), newArgs...)
}

// QueryRow should be avoided as it cannot be mocked using ExpectQuery
func (c *command) QueryRow(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Row, error) {
	if c.useInstrument {
		timer := c.instrument.DatabaseQueryTimer(c.connName, c.connType, name)
		defer timer.ObserveDuration()
	}
	if c.logQuery {
		c.log.Info(ctx, fmt.Sprintf(queryLogMessage, name, replaceBindvarsWithArgs(query, args...)))
	}
	row := c.db.QueryRowxContext(ctx, query, args...)
	return row, row.Err()
}

func (c *command) Query(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Rows, error) {
	if c.useInstrument {
		timer := c.instrument.DatabaseQueryTimer(c.connName, c.connType, name)
		defer timer.ObserveDuration()
	}
	if c.logQuery {
		c.log.Info(ctx, fmt.Sprintf(queryLogMessage, name, replaceBindvarsWithArgs(query, args...)))
	}
	return c.db.QueryxContext(ctx, query, args...)
}

func (c *command) NamedQuery(ctx context.Context, name string, query string, arg interface{}) (*sqlx.Rows, error) {
	if c.useInstrument {
		timer := c.instrument.DatabaseQueryTimer(c.connName, c.connType, name)
		defer timer.ObserveDuration()
	}
	if c.logQuery {
		c.log.Info(ctx, fmt.Sprintf(queryLogMessage, name, replaceBindvarsWithArgs(query)))
	}
	return c.db.NamedQueryContext(ctx, query, arg)
}

func (c *command) Prepare(ctx context.Context, name string, query string) (CommandStmt, error) {
	if c.useInstrument {
		timer := c.instrument.DatabaseQueryTimer(c.connName, c.connType, name)
		defer timer.ObserveDuration()
	}
	stmt, err := c.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return initStmt(ctx, name, c.connName, stmt, c.instrument, c.connType == connTypeLeader, c.useInstrument), nil
}

func (c *command) NamedExec(ctx context.Context, name string, query string, args interface{}) (sql.Result, error) {
	if c.useInstrument {
		timer := c.instrument.DatabaseQueryTimer(c.connName, c.connType, name)
		defer timer.ObserveDuration()
	}
	if c.logQuery {
		c.log.Info(ctx, fmt.Sprintf(queryLogMessage, name, replaceBindvarsWithArgs(query)))
	}
	return c.db.NamedExecContext(ctx, query, args)
}

func (c *command) Exec(ctx context.Context, name string, query string, args ...interface{}) (sql.Result, error) {
	if c.useInstrument {
		timer := c.instrument.DatabaseQueryTimer(c.connName, c.connType, name)
		defer timer.ObserveDuration()
	}
	if c.logQuery {
		c.log.Info(ctx, fmt.Sprintf(queryLogMessage, name, replaceBindvarsWithArgs(query, args...)))
	}
	return c.db.ExecContext(ctx, query, args...)
}

func (c *command) BeginTx(ctx context.Context, name string, opt TxOptions) (CommandTx, error) {
	if c.useInstrument {
		timer := c.instrument.DatabaseQueryTimer(c.connName, c.connType, name)
		defer timer.ObserveDuration()
	}
	opts := &sql.TxOptions{
		Isolation: opt.Isolation,
		ReadOnly:  opt.ReadOnly,
	}
	tx, err := c.db.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return initTx(ctx, name, c.connName, tx, opts, c.log, c.instrument, c.connType == connTypeLeader, c.useInstrument, c.logQuery), nil
}

func (c *command) Get(ctx context.Context, name string, query string, dest interface{}, args ...interface{}) error {
	if c.useInstrument {
		timer := c.instrument.DatabaseQueryTimer(c.connName, c.connType, name)
		defer timer.ObserveDuration()
	}
	if c.logQuery {
		c.log.Info(ctx, fmt.Sprintf(queryLogMessage, name, replaceBindvarsWithArgs(query, args...)))
	}
	return c.db.GetContext(ctx, dest, query, args...)
}
