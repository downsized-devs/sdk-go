//go:build integration
// +build integration

package instrument

import (
	"database/sql"
	"strconv"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_instrument_HTTPRequestTimer(t *testing.T) {
	type args struct {
		path   string
		method string
	}
	type want struct {
		labelCount int
	}
	tests := []struct {
		name     string
		mockFunc func(c Config, a args)
		config   Config
		args     args
		want     want
	}{
		{
			name: "ok",
			mockFunc: func(c Config, a args) {
				i := Init(c)
				timer := i.HTTPRequestTimer(a.path, a.method)
				timer.ObserveDuration()
			},
			config: Config{
				Metrics: MetricsConfig{
					Enabled: true,
				},
			},
			args: args{
				path:   "/test",
				method: "GET",
			},
			want: want{
				labelCount: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.config, tt.args)
			assert.Equal(t, tt.want.labelCount, testutil.CollectAndCount(requestDuration))
		})
	}
}

func Test_instrument_HTTPRequestCounter(t *testing.T) {
	type args struct {
		path   string
		method string
	}

	type want struct {
		labelCount int
		labelValue float64
	}

	tests := []struct {
		name     string
		mockFunc func(c Config, a args)
		config   Config
		args     args
		want     want
	}{
		{
			name: "ok",
			mockFunc: func(c Config, a args) {
				i := Init(c)
				i.HTTPRequestCounter(a.path, a.method)
				i.HTTPRequestCounter(a.path, a.method)
			},
			config: Config{
				Metrics: MetricsConfig{
					Enabled: true,
				},
			},
			args: args{
				path:   "/test",
				method: "GET",
			},
			want: want{
				labelCount: 1,
				labelValue: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.config, tt.args)
			assert.Equal(t, tt.want.labelCount, testutil.CollectAndCount(requestTotal))
			assert.Equal(t, tt.want.labelValue, testutil.ToFloat64(requestTotal.WithLabelValues(tt.args.path, tt.args.method)))
		})
	}
}

func Test_instrument_HTTPResponseStatusCounter(t *testing.T) {
	type args struct {
		code int
	}
	type want struct {
		labelCount int
		labelValue float64
	}
	tests := []struct {
		name     string
		mockFunc func(c Config, a args)
		config   Config
		args     args
		want     want
	}{
		{
			name: "ok",
			mockFunc: func(c Config, a args) {
				i := Init(c)
				i.HTTPResponseStatusCounter(404)
			},
			config: Config{
				Metrics: MetricsConfig{
					Enabled: true,
				},
			},
			args: args{
				code: 404,
			},
			want: want{
				labelCount: 1,
				labelValue: 1,
			},
		},
	}
	for _, tt := range tests {
		tt.mockFunc(tt.config, tt.args)
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want.labelCount, testutil.CollectAndCount(responseStatus))
			assert.Equal(t, tt.want.labelValue, testutil.ToFloat64(responseStatus.WithLabelValues(strconv.Itoa(tt.args.code))))
		})
	}
}

// TODO: nedd to execute real SQL
func Test_instrument_RegisterDBStats(t *testing.T) {
	type fields struct {
		cfg   Config
		prome promeRegistry
	}
	type args struct {
		db     *sql.DB
		dbname string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "ok",
			fields: fields{
				cfg: Config{
					Metrics: MetricsConfig{
						Enabled: true,
					},
				},
				prome: promeRegistry{
					registerer: prometheus.DefaultRegisterer,
					gatherer:   prometheus.DefaultGatherer,
				},
			},
			args: args{
				db:     &sql.DB{},
				dbname: "aa",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			i := &instrument{
				cfg:   tt.fields.cfg,
				prome: tt.fields.prome,
			}
			i.RegisterDBStats(tt.args.db, tt.args.dbname)
		})
	}
}

func Test_instrument_MySQLQueryTimer(t *testing.T) {
	type args struct {
		queryname string
	}
	type want struct {
		labelCount int
	}
	tests := []struct {
		name     string
		mockFunc func(c Config, a args)
		config   Config
		args     args
		want     want
	}{
		{
			name: "ok",
			mockFunc: func(c Config, a args) {
				i := Init(c)
				timer := i.MySQLQueryTimer(a.queryname)
				timer.ObserveDuration()
			},
			config: Config{
				Metrics: MetricsConfig{
					Enabled: true,
				},
			},
			args: args{
				queryname: "rGetlist",
			},
			want: want{
				labelCount: 1,
			},
		},
	}
	for _, tt := range tests {
		tt.mockFunc(tt.config, tt.args)
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want.labelCount, testutil.CollectAndCount(dbQueryDuration))
		})
	}
}
