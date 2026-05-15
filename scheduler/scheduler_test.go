package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/downsized-devs/sdk-go/logger"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newReal() Interface {
	return Init(Config{}, logger.Init(logger.Config{}))
}

func TestNew_ReturnsScheduler(t *testing.T) {
	s := newReal()
	require.NotNil(t, s)
}

func TestStartShutdown_NoJobs(t *testing.T) {
	s := newReal()
	ctx := context.Background()
	s.Start(ctx)
	s.Shutdown(ctx)
}

func TestRegister_AllJobTypes(t *testing.T) {
	cases := []struct {
		name string
		opt  JobOption
	}{
		{"duration", JobOption{JobType: Duration, Duration: time.Hour}},
		{"random-duration", JobOption{JobType: RandomDuration, Duration: time.Hour, Jitter: 10 * time.Minute}},
		{"daily", JobOption{JobType: Daily, RunningTime: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC)}},
		{"weekly", JobOption{JobType: Weekly, RunningDay: time.Monday, RunningTime: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC)}},
		{"monthly", JobOption{JobType: Monthly, RunningDate: 1, RunningTime: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC)}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := newReal()
			ctx := context.Background()
			err := s.Register(ctx, tc.opt, func() {})
			require.NoError(t, err)
			s.Shutdown(ctx)
		})
	}
}

func TestRegister_UnknownTypeReturnsError(t *testing.T) {
	s := newReal()
	ctx := context.Background()
	err := s.Register(ctx, JobOption{JobType: "nope", Duration: time.Hour}, func() {})
	assert.Error(t, err, "gocron rejects a nil JobDefinition with an error")
	s.Shutdown(ctx)
}

// ----------------- Shutdown error path uses a fake engine ----------------- //

type fakeEngine struct {
	gocron.Scheduler // embed to satisfy unused methods
	shutdownErr      error
	started          bool
}

func (f *fakeEngine) Start()             { f.started = true }
func (f *fakeEngine) Shutdown() error    { return f.shutdownErr }
func (f *fakeEngine) Jobs() []gocron.Job { return nil }
func (f *fakeEngine) NewJob(_ gocron.JobDefinition, _ gocron.Task, _ ...gocron.JobOption) (gocron.Job, error) {
	return nil, errors.New("not implemented in fake")
}
func (f *fakeEngine) RemoveByTags(_ ...string) {}
func (f *fakeEngine) RemoveJob(_ uuid.UUID) error {
	return nil
}
func (f *fakeEngine) StopJobs() error { return nil }

func TestShutdown_ErrorIsLogged(t *testing.T) {
	s := &scheduler{
		log:    logger.Init(logger.Config{}),
		engine: &fakeEngine{shutdownErr: errors.New("shutdown boom")},
	}
	// We just want to assert no panic when Shutdown returns an error.
	assert.NotPanics(t, func() { s.Shutdown(context.Background()) })
}
