//go:build integration
// +build integration

package instrument

import (
	"database/sql"
	"strconv"
	"testing"

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
		mockFunc func(i *instrument, a args)
		config   Config
		args     args
		want     want
	}{
		{
			name: "ok",
			mockFunc: func(i *instrument, a args) {
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
			instr := Init(tt.config).(*instrument)
			tt.mockFunc(instr, tt.args)
			assert.Equal(t, tt.want.labelCount, testutil.CollectAndCount(instr.requestDuration))
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
		mockFunc func(i *instrument, a args)
		config   Config
		args     args
		want     want
	}{
		{
			name: "ok",
			mockFunc: func(i *instrument, a args) {
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
			instr := Init(tt.config).(*instrument)
			tt.mockFunc(instr, tt.args)
			assert.Equal(t, tt.want.labelCount, testutil.CollectAndCount(instr.requestTotal))
			assert.Equal(t, tt.want.labelValue, testutil.ToFloat64(instr.requestTotal.WithLabelValues(tt.args.path, tt.args.method)))
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
		mockFunc func(i *instrument, a args)
		config   Config
		args     args
		want     want
	}{
		{
			name: "ok",
			mockFunc: func(i *instrument, a args) {
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
		t.Run(tt.name, func(t *testing.T) {
			instr := Init(tt.config).(*instrument)
			tt.mockFunc(instr, tt.args)
			assert.Equal(t, tt.want.labelCount, testutil.CollectAndCount(instr.responseStatus))
			assert.Equal(t, tt.want.labelValue, testutil.ToFloat64(instr.responseStatus.WithLabelValues(strconv.Itoa(tt.args.code))))
		})
	}
}

// TODO: need to execute real SQL
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
			},
			args: args{
				db:     &sql.DB{},
				dbname: "aa",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Init(tt.fields.cfg)
			i.RegisterDBStats(tt.args.db, tt.args.dbname)
		})
	}
}

func Test_instrument_DatabaseQueryTimer(t *testing.T) {
	type args struct {
		queryname string
	}
	type want struct {
		labelCount int
	}
	tests := []struct {
		name     string
		mockFunc func(i *instrument, a args)
		config   Config
		args     args
		want     want
	}{
		{
			name: "ok",
			mockFunc: func(i *instrument, a args) {
				timer := i.DatabaseQueryTimer("testdb", "leader", a.queryname)
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
		t.Run(tt.name, func(t *testing.T) {
			instr := Init(tt.config).(*instrument)
			tt.mockFunc(instr, tt.args)
			assert.Equal(t, tt.want.labelCount, testutil.CollectAndCount(instr.dbQueryDuration))
		})
	}
}
