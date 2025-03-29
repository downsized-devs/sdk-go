package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/downsized-devs/sdk-go/logger"
	"github.com/go-co-op/gocron/v2"
)

const (
	Duration       = "duration"
	RandomDuration = "random-duration"
	Daily          = "daily"
	Weekly         = "weekly"
	Monthly        = "monthly"

	Monday    = "monday"
	Tuesday   = "tuesday"
	Wednesday = "wednesday"
	Thursday  = "thursday"
	Friday    = "friday"
	Saturday  = "saturday"
	Sunday    = "sunday"
)

type Interface interface {
	Start(ctx context.Context)
	Shutdown(ctx context.Context)
	Register(ctx context.Context, opt JobOption, handlerFunc any) error
}

type Config struct{}

type JobOption struct {
	JobType     string
	Duration    time.Duration
	Jitter      time.Duration
	RunningDate int
	RunningDay  time.Weekday
	RunningTime time.Time
}

type scheduler struct {
	log    logger.Interface
	engine gocron.Scheduler
}

func New(cfg Config, log logger.Interface) Interface {
	engine, err := gocron.NewScheduler()
	if err != nil {
		log.Panic(err)
		return nil
	}

	return &scheduler{
		log:    log,
		engine: engine,
	}
}

func (s *scheduler) Start(ctx context.Context) {
	s.engine.Start()
	s.log.Info(ctx, "running all available scheduler")
}

func (s *scheduler) Shutdown(ctx context.Context) {
	if err := s.engine.Shutdown(); err != nil {
		s.log.Error(ctx, err)
	}
}

func (s *scheduler) Register(ctx context.Context, opt JobOption, handlerFunc any) error {
	var jobOption gocron.JobDefinition
	switch opt.JobType {
	case Duration:
		jobOption = gocron.DurationJob(opt.Duration)
	case RandomDuration:
		jobOption = gocron.DurationRandomJob(opt.Duration-opt.Jitter, opt.Duration+opt.Jitter)
	case Daily:
		jobOption = gocron.DailyJob(1, gocron.NewAtTimes(
			gocron.NewAtTime(uint(opt.RunningTime.Hour()), uint(opt.RunningTime.Minute()), uint(opt.RunningTime.Second())), //nolint:gosec
		))
	case Weekly:
		jobOption = gocron.WeeklyJob(1, gocron.NewWeekdays(opt.RunningDay), gocron.NewAtTimes(
			gocron.NewAtTime(uint(opt.RunningTime.Hour()), uint(opt.RunningTime.Minute()), uint(opt.RunningTime.Second())), //nolint:gosec
		))
	case Monthly:
		jobOption = gocron.MonthlyJob(1, gocron.NewDaysOfTheMonth(opt.RunningDate), gocron.NewAtTimes(
			gocron.NewAtTime(uint(opt.RunningTime.Hour()), uint(opt.RunningTime.Minute()), uint(opt.RunningTime.Second())), //nolint:gosec
		))
	}

	job, err := s.engine.NewJob(jobOption, gocron.NewTask(handlerFunc))
	if err != nil {
		return err
	}

	nextRun, err := job.NextRun()
	if err != nil {
		s.log.Error(ctx, err)
	}

	s.log.Debug(ctx, fmt.Sprintf("%s(%s) running the first time at %v", job.Name(), job.ID(), nextRun.Format(time.RFC3339)))

	return nil
}
