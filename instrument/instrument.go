package instrument

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	Metrics MetricsConfig
}

type MetricsConfig struct {
	Enabled bool
}

type Interface interface {
	// Other Functions
	IsEnabled() bool
	// HTTP Handler
	MetricsHandler() http.Handler
	// HTTP Metrics
	HTTPRequestTimer(path, method string) *prometheus.Timer
	HTTPRequestCounter(path, method string)
	HTTPResponseStatusCounter(code int)
	// Database Metrics
	RegisterDBStats(db *sql.DB, dbname string)
	DatabaseQueryTimer(dbname, conntype, queryname string) *prometheus.Timer
	// Scheduler Metrics
	SchedulerRunningCounter(schedulername string)
	SchedulerRunningTimer(schedulername string) *prometheus.Timer
}

type instrument struct {
	cfg              Config
	prome            promeRegistry
	requestTotal     *prometheus.CounterVec
	responseStatus   *prometheus.CounterVec
	requestDuration  *prometheus.HistogramVec
	dbQueryDuration  *prometheus.HistogramVec
	schedulerTotal   *prometheus.CounterVec
	schedulerDuration *prometheus.HistogramVec
}

type promeRegistry struct {
	registerer prometheus.Registerer
	gatherer   prometheus.Gatherer
}

func Init(cfg Config) Interface {
	instr := &instrument{cfg: cfg}

	if !cfg.Metrics.Enabled {
		return instr
	}

	registry := prometheus.NewRegistry()
	// Register default Go / process collectors.
	registry.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)

	instr.requestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of API Requests",
		},
		[]string{"path", "method"},
	)
	instr.responseStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_status_total",
			Help: "Number of HTTP Status response",
		},
		[]string{"status"},
	)
	instr.requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests",
		},
		[]string{"path", "method"},
	)
	instr.dbQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "db_query_duration_seconds",
			Help: "Duration of Database query execution",
		},
		[]string{"database", "connection_type", "query_name"},
	)
	instr.schedulerTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "scheduler_running_total",
			Help: "Number of Running Scheduler",
		},
		[]string{"scheduler_name"},
	)
	instr.schedulerDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "scheduler_running_duration_seconds",
			Help: "Duration of Running Scheduler",
		},
		[]string{"scheduler_name"},
	)

	registry.MustRegister(
		instr.requestTotal,
		instr.responseStatus,
		instr.requestDuration,
		instr.dbQueryDuration,
		instr.schedulerTotal,
		instr.schedulerDuration,
	)

	instr.prome = promeRegistry{
		registerer: registry,
		gatherer:   registry,
	}

	return instr
}

func (i *instrument) IsEnabled() bool {
	return i.cfg.Metrics.Enabled
}

func (i *instrument) MetricsHandler() http.Handler {
	if !i.cfg.Metrics.Enabled {
		return http.NotFoundHandler()
	}
	return promhttp.InstrumentMetricHandler(
		i.prome.registerer, promhttp.HandlerFor(i.prome.gatherer, promhttp.HandlerOpts{}),
	)
}

func (i *instrument) HTTPRequestTimer(path, method string) *prometheus.Timer {
	if !i.cfg.Metrics.Enabled {
		return prometheus.NewTimer(prometheus.ObserverFunc(func(float64) {}))
	}
	return prometheus.NewTimer(i.requestDuration.WithLabelValues(path, method))
}

func (i *instrument) HTTPRequestCounter(path, method string) {
	if !i.cfg.Metrics.Enabled {
		return
	}
	i.requestTotal.WithLabelValues(path, method).Inc()
}

func (i *instrument) HTTPResponseStatusCounter(code int) {
	if !i.cfg.Metrics.Enabled {
		return
	}
	i.responseStatus.WithLabelValues(strconv.Itoa(code)).Inc()
}

func (i *instrument) DatabaseQueryTimer(dbname, conntype, queryname string) *prometheus.Timer {
	if !i.cfg.Metrics.Enabled {
		return prometheus.NewTimer(prometheus.ObserverFunc(func(float64) {}))
	}
	return prometheus.NewTimer(i.dbQueryDuration.WithLabelValues(dbname, conntype, queryname))
}

func (i *instrument) RegisterDBStats(db *sql.DB, dbname string) {
	if !i.cfg.Metrics.Enabled {
		return
	}
	i.prome.registerer.MustRegister(collectors.NewDBStatsCollector(db, dbname))
}

// SchedulerRunningCounter increments the running-scheduler counter.
func (i *instrument) SchedulerRunningCounter(schedulername string) {
	if !i.cfg.Metrics.Enabled {
		return
	}
	i.schedulerTotal.WithLabelValues(schedulername).Inc()
}

// SchedulerRunningTimer returns a duration timer for the named scheduler.
func (i *instrument) SchedulerRunningTimer(schedulername string) *prometheus.Timer {
	if !i.cfg.Metrics.Enabled {
		return prometheus.NewTimer(prometheus.ObserverFunc(func(float64) {}))
	}
	return prometheus.NewTimer(i.schedulerDuration.WithLabelValues(schedulername))
}
