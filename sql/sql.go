// Package sql is a sqlx-backed SQL client with leader/follower routing,
// transactions, prepared statements, and optional Prometheus
// instrumentation. It targets MySQL, Postgres, and SQLite drivers.
package sql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/downsized-devs/sdk-go/codes"
	"github.com/downsized-devs/sdk-go/errors"
	"github.com/downsized-devs/sdk-go/instrument"
	"github.com/downsized-devs/sdk-go/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

var ErrNotFound = sql.ErrNoRows

const (
	connTypeLeader   = "leader"
	connTypeFollower = "follower"
)

// Config controls how Init builds the SQL client: which driver to use,
// connection settings for the leader and follower, whether to emit
// instrumentation metrics, and whether to log every query (useful in
// development and noisy in production).
type Config struct {
	UseInstrument bool
	LogQuery      bool
	Driver        string
	Name          string
	Leader        ConnConfig
	Follower      ConnConfig
}

// ConnConfig describes a single connection. MockDB, when non-nil, is
// preferred over the Host/Port/DB/User/Password fields and lets tests
// inject a pre-built *sql.DB (for example a sqlmock or an in-memory
// sqlite database).
type ConnConfig struct {
	Host     string
	Port     int
	DB       string
	User     string
	Password string
	SSL      bool
	Schema   string
	Options  ConnOptions
	MockDB   *sql.DB
}

// ConnOptions tunes the underlying database/sql pool.
type ConnOptions struct {
	MaxLifeTime time.Duration
	MaxIdle     int
	MaxOpen     int
}

// Interface is the public surface of the sql package. Leader returns
// the read/write Command, Follower returns the read-only Command (or
// the leader when no follower is configured), and Stop closes both
// pools at most once.
type Interface interface {
	Leader() Command
	Follower() Command
	Stop()
}

type sqlDB struct {
	endOnce    *sync.Once
	leader     Command
	follower   Command
	cfg        Config
	log        logger.Interface
	instrument instrument.Interface
}

// Init constructs a sql client from cfg. When cfg.Driver is "sqlmock"
// instrumentation is force-disabled. cfg.Name defaults to cfg.Driver
// when empty. Init calls log.Fatal if the leader connection cannot be
// established.
func Init(cfg Config, log logger.Interface, instr instrument.Interface) Interface {
	if cfg.Driver == "sqlmock" {
		cfg.UseInstrument = false
	}

	if cfg.Name == "" {
		cfg.Name = cfg.Driver
	}
	sql := &sqlDB{
		endOnce:    &sync.Once{},
		log:        log,
		cfg:        cfg,
		instrument: instr,
	}

	sql.initDB()
	return sql
}

func (s *sqlDB) Leader() Command {
	return s.leader
}

func (s *sqlDB) Follower() Command {
	return s.follower
}

func (s *sqlDB) Stop() {
	s.endOnce.Do(func() {
		ctx := context.Background()
		if s.leader != nil {
			if err := s.leader.Close(); err != nil {
				s.log.Error(ctx, err)
			}
		}
		if s.follower != nil {
			if err := s.follower.Close(); err != nil {
				s.log.Error(ctx, err)
			}
		}
	})
}

func (s *sqlDB) initDB() {
	ctx := context.Background()
	db, err := s.connect(true)
	if err != nil {
		s.log.Fatal(ctx, fmt.Sprintf("[FATAL] cannot connect to db %s leader: %s on port %d, with error: %s", s.cfg.Leader.DB, s.cfg.Leader.Host, s.cfg.Leader.Port, err))
	}
	s.log.Info(ctx, fmt.Sprintf("SQL: [LEADER] driver=%s db=%s @%s:%v ssl=%v", s.cfg.Driver, s.cfg.Leader.DB, s.cfg.Leader.Host, s.cfg.Leader.Port, s.cfg.Leader.SSL))
	s.leader = initCommand(db, s.cfg.Name, s.instrument, s.log, true, s.cfg.UseInstrument, s.cfg.LogQuery)

	if s.isFollowerEnabled() {
		db, err = s.connect(false)
		if err != nil {
			s.log.Fatal(ctx, fmt.Sprintf("[FATAL] cannot connect to db %s leader: %s on port %d, with error: %s", s.cfg.Leader.DB, s.cfg.Leader.Host, s.cfg.Leader.Port, err))
		}
		s.log.Info(ctx, fmt.Sprintf("SQL: [FOLLOWER] driver=%s db=%s @%s:%v ssl=%v", s.cfg.Driver, s.cfg.Follower.DB, s.cfg.Follower.Host, s.cfg.Follower.Port, s.cfg.Leader.SSL))
		s.follower = initCommand(db, s.cfg.Name, s.instrument, s.log, false, s.cfg.UseInstrument, s.cfg.LogQuery)
	} else {
		s.follower = s.leader
	}
}

func (s *sqlDB) connect(toLeader bool) (*sqlx.DB, error) {
	conf := s.cfg.Leader
	if !toLeader {
		conf = s.cfg.Follower
	}

	if !toLeader {
		if s.cfg.Follower.MockDB != nil {
			return sqlx.NewDb(s.cfg.Follower.MockDB, s.cfg.Driver), nil
		}
	} else {
		if s.cfg.Leader.MockDB != nil {
			return sqlx.NewDb(s.cfg.Leader.MockDB, s.cfg.Driver), nil
		}
	}

	uri, err := s.getURI(conf)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(s.cfg.Driver, uri)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, errors.NewWithCode(codes.CodeSQLInit, "%s", err.Error())
	}

	sqlxDB := sqlx.NewDb(db, s.cfg.Driver)
	sqlxDB.SetMaxOpenConns(conf.Options.MaxOpen)
	sqlxDB.SetMaxIdleConns(conf.Options.MaxIdle)
	sqlxDB.SetConnMaxLifetime(conf.Options.MaxLifeTime)

	if s.cfg.UseInstrument {
		if toLeader {
			s.instrument.RegisterDBStats(sqlxDB.DB, fmt.Sprintf("%s_%s", s.cfg.Name, connTypeLeader))
		} else {
			s.instrument.RegisterDBStats(sqlxDB.DB, fmt.Sprintf("%s_%s", s.cfg.Name, connTypeFollower))
		}
	}

	return sqlxDB, nil
}

func (s *sqlDB) isFollowerEnabled() bool {
	if s.cfg.Follower.Host == "" {
		return false
	}
	leaderAddr := fmt.Sprintf("%s:%d", s.cfg.Leader.Host, s.cfg.Leader.Port)
	followerAddr := fmt.Sprintf("%s:%d", s.cfg.Follower.Host, s.cfg.Follower.Port)
	return leaderAddr != followerAddr
}

func (s *sqlDB) getURI(conf ConnConfig) (string, error) {
	switch s.cfg.Driver {
	case "postgres":
		ssl := `disable`
		if conf.SSL {
			ssl = `require`
		}
		if conf.Schema == "" {
			conf.Schema = "public"
		}
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=%s", conf.Host, conf.Port, conf.User, conf.Password, conf.DB, conf.Schema, ssl), nil
	case "mysql":
		ssl := `false`
		if conf.SSL {
			ssl = `true`
		}
		return fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?tls=%s&parseTime=true", conf.User, conf.Password, conf.Host, conf.Port, conf.DB, ssl), nil
	case "sqlite3":
		return conf.DB, nil
	default:
		return "", fmt.Errorf(`DB Driver [%s] is not supported`, s.cfg.Driver)
	}
}
