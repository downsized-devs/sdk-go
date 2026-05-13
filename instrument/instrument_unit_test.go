package instrument

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

// Unit-test counterpart to instrument_test.go (which is build-tagged). The
// integration tests already exercise the metrics machinery; these mirror them
// for the default unit-test run and additionally cover the Metrics.Enabled=false
// short-circuit branches.

func Test_instrument_Init_Disabled(t *testing.T) {
	i := Init(Config{}).(*instrument)
	assert.False(t, i.IsEnabled())
	assert.Nil(t, i.requestTotal)
}

func Test_instrument_Init_Enabled(t *testing.T) {
	i := Init(Config{Metrics: MetricsConfig{Enabled: true}}).(*instrument)
	assert.True(t, i.IsEnabled())
	assert.NotNil(t, i.requestTotal)
	assert.NotNil(t, i.responseStatus)
	assert.NotNil(t, i.requestDuration)
	assert.NotNil(t, i.dbQueryDuration)
	assert.NotNil(t, i.schedulerTotal)
	assert.NotNil(t, i.schedulerDuration)
}

func Test_instrument_MetricsHandler(t *testing.T) {
	t.Run("disabled returns 404 handler", func(t *testing.T) {
		i := Init(Config{})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		i.MetricsHandler().ServeHTTP(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("enabled returns 200", func(t *testing.T) {
		i := Init(Config{Metrics: MetricsConfig{Enabled: true}})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		i.MetricsHandler().ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func Test_instrument_HTTPRequestTimer_Unit(t *testing.T) {
	t.Run("disabled returns no-op timer", func(t *testing.T) {
		i := Init(Config{})
		timer := i.HTTPRequestTimer("/x", "GET")
		assert.NotNil(t, timer)
		timer.ObserveDuration()
	})
	t.Run("enabled records observation", func(t *testing.T) {
		i := Init(Config{Metrics: MetricsConfig{Enabled: true}}).(*instrument)
		timer := i.HTTPRequestTimer("/test", "GET")
		timer.ObserveDuration()
		assert.Equal(t, 1, testutil.CollectAndCount(i.requestDuration))
	})
}

func Test_instrument_HTTPRequestCounter_Unit(t *testing.T) {
	t.Run("disabled is a no-op", func(t *testing.T) {
		i := Init(Config{})
		// Just verify no panic.
		i.HTTPRequestCounter("/x", "GET")
	})
	t.Run("enabled increments", func(t *testing.T) {
		i := Init(Config{Metrics: MetricsConfig{Enabled: true}}).(*instrument)
		i.HTTPRequestCounter("/test", "GET")
		i.HTTPRequestCounter("/test", "GET")
		assert.Equal(t, float64(2), testutil.ToFloat64(i.requestTotal.WithLabelValues("/test", "GET")))
	})
}

func Test_instrument_HTTPResponseStatusCounter_Unit(t *testing.T) {
	t.Run("disabled is a no-op", func(t *testing.T) {
		i := Init(Config{})
		i.HTTPResponseStatusCounter(200)
	})
	t.Run("enabled increments", func(t *testing.T) {
		i := Init(Config{Metrics: MetricsConfig{Enabled: true}}).(*instrument)
		i.HTTPResponseStatusCounter(404)
		assert.Equal(t, float64(1), testutil.ToFloat64(i.responseStatus.WithLabelValues(strconv.Itoa(404))))
	})
}

func Test_instrument_DatabaseQueryTimer_Unit(t *testing.T) {
	t.Run("disabled returns no-op", func(t *testing.T) {
		i := Init(Config{})
		timer := i.DatabaseQueryTimer("db", "leader", "q")
		assert.NotNil(t, timer)
		timer.ObserveDuration()
	})
	t.Run("enabled records observation", func(t *testing.T) {
		i := Init(Config{Metrics: MetricsConfig{Enabled: true}}).(*instrument)
		timer := i.DatabaseQueryTimer("testdb", "leader", "rGet")
		timer.ObserveDuration()
		assert.Equal(t, 1, testutil.CollectAndCount(i.dbQueryDuration))
	})
}

func Test_instrument_RegisterDBStats_Unit(t *testing.T) {
	t.Run("disabled is a no-op", func(t *testing.T) {
		// With metrics disabled, no DB stats are registered and the empty *sql.DB
		// pointer is never dereferenced.
		i := Init(Config{})
		i.RegisterDBStats(&sql.DB{}, "ignored")
	})
}

func Test_instrument_SchedulerRunningCounter_Unit(t *testing.T) {
	t.Run("disabled is a no-op", func(t *testing.T) {
		Init(Config{}).SchedulerRunningCounter("s1")
	})
	t.Run("enabled increments", func(t *testing.T) {
		i := Init(Config{Metrics: MetricsConfig{Enabled: true}}).(*instrument)
		i.SchedulerRunningCounter("s1")
		i.SchedulerRunningCounter("s1")
		assert.Equal(t, float64(2), testutil.ToFloat64(i.schedulerTotal.WithLabelValues("s1")))
	})
}

func Test_instrument_SchedulerRunningTimer_Unit(t *testing.T) {
	t.Run("disabled returns no-op timer", func(t *testing.T) {
		timer := Init(Config{}).SchedulerRunningTimer("s1")
		assert.NotNil(t, timer)
		timer.ObserveDuration()
	})
	t.Run("enabled records observation", func(t *testing.T) {
		i := Init(Config{Metrics: MetricsConfig{Enabled: true}}).(*instrument)
		timer := i.SchedulerRunningTimer("s1")
		timer.ObserveDuration()
		assert.Equal(t, 1, testutil.CollectAndCount(i.schedulerDuration))
	})
}
