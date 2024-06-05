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
	"github.com/downsized-devs/sdk-go/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var ErrNotFound = sql.ErrNoRows

const (
	connTypeLeader   = "leader"
	connTypeFollower = "follower"
)

type Config struct {
	UseInstrument bool
	LogQuery      bool
	Driver        string
	Name          string
	Leader        ConnConfig
	Follower      ConnConfig
}

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

type ConnOptions struct {
	MaxLifeTime time.Duration
	MaxIdle     int
	MaxOpen     int
}

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
	log        log.Interface
	instrument instrument.Interface
}

func Init(cfg Config, log log.Interface, instr instrument.Interface) Interface {
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
		if s.cfg.Leader.MockDB != nil {
			return sqlx.NewDb(s.cfg.Leader.MockDB, s.cfg.Driver), nil
		}
	} else {
		if s.cfg.Follower.MockDB != nil {
			return sqlx.NewDb(s.cfg.Follower.MockDB, s.cfg.Driver), nil
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
		return nil, errors.NewWithCode(codes.CodeSQLInit, err.Error())
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
	isHostNotEmpty := s.cfg.Follower.Host != ""
	isHostDifferent := (s.cfg.Follower.Host != s.cfg.Leader.Host && s.cfg.Follower.Port == s.cfg.Leader.Port)
	isPortDifferent := (s.cfg.Follower.Host == s.cfg.Leader.Host && s.cfg.Follower.Port != s.cfg.Leader.Port)
	return isHostNotEmpty && (isHostDifferent || isPortDifferent)
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
	default:
		return "", fmt.Errorf(`DB Driver [%s] is not supported`, s.cfg.Driver)
	}
}
