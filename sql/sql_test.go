//go:build integration
// +build integration

package sql

import (
	"testing"

	"github.com/downsized-devs/sdk-go/instrument"
	"github.com/downsized-devs/sdk-go/log"
	mock_log "github.com/downsized-devs/sdk-go/tests/mock/log"
	"go.uber.org/mock/gomock"
)

func initTestMysqlDatabase() Interface {
	return Init(Config{
		UseInstrument: true,
		LogQuery:      true,
		Driver:        "mysql",
		Leader: ConnConfig{
			Host:     "127.0.0.1",
			Port:     3306,
			DB:       "delos_db",
			User:     "root",
			Password: "password",
		},
		Follower: ConnConfig{
			Host:     "127.0.0.1",
			Port:     3306,
			DB:       "delos_db",
			User:     "root",
			Password: "password",
		},
	}, log.Init(log.Config{Level: "debug"}),
		instrument.Init(instrument.Config{
			Metrics: instrument.MetricsConfig{
				Enabled: true,
			},
		}))
}

func initTestPostgresDatabase() Interface {
	return Init(Config{
		UseInstrument: true,
		LogQuery:      true,
		Driver:        "postgres",
		Leader: ConnConfig{
			Host:     "127.0.0.1",
			Port:     5432,
			DB:       "delos_db",
			User:     "postgres",
			Password: "",
		},
		Follower: ConnConfig{
			Host:     "127.0.0.1",
			Port:     5432,
			DB:       "delos_db",
			User:     "root",
			Password: "",
		},
	}, log.Init(log.Config{Level: "debug"}),
		instrument.Init(instrument.Config{
			Metrics: instrument.MetricsConfig{
				Enabled: true,
			},
		}))
}

func TestInit(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Fatal(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	conConfig := ConnConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		DB:       "delos_db",
		User:     "root",
		Password: "password",
	}
	logInit := log.Init(log.Config{Level: "debug"})
	instr := instrument.Init(instrument.Config{
		Metrics: instrument.MetricsConfig{
			Enabled: true,
		},
	})

	type args struct {
		cfg   Config
		log   log.Interface
		instr instrument.Interface
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "no driver name",
			args: args{
				cfg: Config{
					UseInstrument: true,
					LogQuery:      true,
					Driver:        "",
					Leader:        conConfig,
					Follower:      conConfig,
				},
				log:   logger,
				instr: instr,
			},
		},
		{
			name: "ok",
			args: args{
				cfg: Config{
					UseInstrument: true,
					LogQuery:      true,
					Driver:        "mysql",
					Leader:        conConfig,
					Follower:      conConfig,
				},
				log:   logInit,
				instr: instr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.args.cfg, tt.args.log, tt.args.instr)
		})
	}
}

func TestLeader(t *testing.T) {
	mockMysqlDb := initTestMysqlDatabase()
	mockPostgresDb := initTestPostgresDatabase()

	tests := []struct {
		name string
	}{
		{
			name: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMysqlDb.Leader()
			mockPostgresDb.Leader()
		})
	}
}

func TestFollower(t *testing.T) {
	mockMysqlDb := initTestMysqlDatabase()
	mockPostgresDb := initTestPostgresDatabase()

	tests := []struct {
		name string
	}{
		{
			name: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMysqlDb.Follower()
			mockPostgresDb.Follower()
		})
	}
}

func TestStop(t *testing.T) {
	mockMysqlDb := initTestMysqlDatabase()
	mockPostgresDb := initTestPostgresDatabase()

	tests := []struct {
		name string
	}{
		{
			name: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMysqlDb.Stop()
			mockPostgresDb.Stop()
		})
	}
}
