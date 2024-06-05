package instrument

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of API Requests",
		},
		[]string{"path", "method"},
	)
	responseStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_status_total",
			Help: "Number of HTTP Status response",
		},
		[]string{"status"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests",
		},
		[]string{"path", "method"},
	)
	dbQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "db_query_duration_seconds",
			Help: "Duration of Database query execution",
		},
		[]string{"database", "connection_type", "query_name"},
	)
	schedulerTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "scheduler_running_total",
			Help: "Number of Running Scheduler",
		},
		[]string{"scheduler_name"},
	)
	schedulerDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "scheduler_running_duration_seconds",
			Help: "Duration of Running Scheduler",
		},
		[]string{"scheduler_name"},
	)
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
	cfg   Config
	prome promeRegistry
}

type promeRegistry struct {
	registerer prometheus.Registerer
	gatherer   prometheus.Gatherer
}

func Init(cfg Config) Interface {
	if cfg.Metrics.Enabled {
		var registry = prometheus.NewRegistry()
		var registerer prometheus.Registerer = registry
		var gatherer prometheus.Gatherer = registry
		// Register Default Metrics
		registerer.MustRegister(
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
			collectors.NewGoCollector(),
		)
		// Register HTTP Metrics
		registerer.Register(requestTotal)
		registerer.Register(responseStatus)
		registerer.Register(requestDuration)
		// Register Database Metrics
		registerer.Register(dbQueryDuration)
		// Register Scheduler Metrics
		registerer.Register(schedulerTotal)
		registerer.Register(schedulerDuration)

		return &instrument{
			cfg: cfg,
			prome: promeRegistry{
				registerer: registerer,
				gatherer:   gatherer,
			},
		}
	} else {
		return &instrument{cfg: cfg}
	}
}

func (i *instrument) IsEnabled() bool {
	return i.cfg.Metrics.Enabled
}

func (i *instrument) MetricsHandler() http.Handler {
	return promhttp.InstrumentMetricHandler(
		i.prome.registerer, promhttp.HandlerFor(i.prome.gatherer, promhttp.HandlerOpts{}),
	)
}

func (i *instrument) HTTPRequestTimer(path, method string) *prometheus.Timer {
	return prometheus.NewTimer(requestDuration.WithLabelValues(path, method))
}

func (i *instrument) HTTPRequestCounter(path, method string) {
	requestTotal.WithLabelValues(path, method).Inc()
}

func (i *instrument) HTTPResponseStatusCounter(code int) {
	responseStatus.WithLabelValues(strconv.Itoa(code)).Inc()
}

func (i *instrument) DatabaseQueryTimer(dbname, conntype, queryname string) *prometheus.Timer {
	return prometheus.NewTimer(dbQueryDuration.WithLabelValues(dbname, conntype, queryname))
}

func (i *instrument) RegisterDBStats(db *sql.DB, dbname string) {
	i.prome.registerer.MustRegister(collectors.NewDBStatsCollector(db, dbname))
}

// Increase Running Scheduler Counter
func (i *instrument) SchedulerRunningCounter(schedulername string) {
	schedulerTotal.WithLabelValues(schedulername).Inc()
}

// Return Duration Timer for Running Scheduler
func (i *instrument) SchedulerRunningTimer(schedulername string) *prometheus.Timer {
	return prometheus.NewTimer(schedulerDuration.WithLabelValues(schedulername))
}
