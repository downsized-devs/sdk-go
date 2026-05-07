package logger

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/downsized-devs/sdk-go/appcontext"
	"github.com/downsized-devs/sdk-go/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var now = time.Now

// defaultCallerSkipFrameCount is the default number of stack frames to skip
// when reporting the caller location in log entries.
const defaultCallerSkipFrameCount = 3

type Interface interface {
	Trace(ctx context.Context, obj any)
	Debug(ctx context.Context, obj any)
	// Debugf logs a formatted debug message.
	Debugf(ctx context.Context, format string, args ...any)
	Info(ctx context.Context, obj any)
	Warn(ctx context.Context, obj any)
	Error(ctx context.Context, obj any)
	Fatal(ctx context.Context, obj any)
	Panic(obj any)
}

type Config struct {
	Level string
	// CallerSkipFrameCount controls the total number of stack frames skipped
	// when reporting the caller. A value of 0 uses the default (3).
	CallerSkipFrameCount int
}

type logger struct {
	log zerolog.Logger
}

func DefaultLogger() Interface {
	return &logger{
		log: zerolog.New(os.Stdout).
			With().
			Timestamp().
			CallerWithSkipFrameCount(defaultCallerSkipFrameCount).
			Logger().
			Level(zerolog.DebugLevel),
	}
}

// Init creates a new logger configured according to cfg.  Each call returns an
// independent logger instance; there is no global singleton.
func Init(cfg Config) Interface {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("failed to parse log level %q: %v", cfg.Level, err))
	}

	skipFrames := defaultCallerSkipFrameCount
	if cfg.CallerSkipFrameCount > 0 {
		skipFrames = cfg.CallerSkipFrameCount
	}

	zl := zerolog.New(os.Stdout).
		With().
		Timestamp().
		CallerWithSkipFrameCount(skipFrames).
		Logger().
		Level(level)

	return &logger{log: zl}
}

func (l *logger) Trace(ctx context.Context, obj any) {
	l.log.Trace().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *logger) Debug(ctx context.Context, obj any) {
	l.log.Debug().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *logger) Debugf(ctx context.Context, format string, args ...any) {
	l.log.Debug().
		Fields(getContextFields(ctx)).
		Msgf(format, args...)
}

func (l *logger) Info(ctx context.Context, obj any) {
	l.log.Info().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *logger) Warn(ctx context.Context, obj any) {
	l.log.Warn().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *logger) Error(ctx context.Context, obj any) {
	l.log.Error().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *logger) Fatal(ctx context.Context, obj any) {
	l.log.Fatal().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *logger) Panic(obj any) {
	defer func() { recover() }()
	l.log.Panic().
		Fields(getPanicStacktrace()).
		Msg(fmt.Sprint(getCaller(obj)))
}

func getPanicStacktrace() map[string]any {
	errStack := strings.Split(strings.ReplaceAll(string(debug.Stack()), "\t", ""), "\n")
	return map[string]any{
		"stacktrace": errStack,
	}
}

func getCaller(obj any) any {
	switch tr := obj.(type) {
	case error:
		file, line, msg, err := errors.GetCaller(tr)
		if err == nil {
			obj = fmt.Sprintf("%s:%#v --- %s", file, line, msg)
		}
	case string:
		obj = tr
	default:
		obj = fmt.Sprintf("%#v", tr)
	}

	return obj
}

func getContextFields(ctx context.Context) map[string]any {
	reqstart := appcontext.GetRequestStartTime(ctx)
	apprespcode := appcontext.GetAppResponseCode(ctx)
	appErrMsg := appcontext.GetAppErrorMessage(ctx)
	timeElapsed := "0ms"
	if !time.Time.IsZero(reqstart) {
		timeElapsed = fmt.Sprintf("%dms", int64(now().Sub(reqstart)/time.Millisecond))
	}

	cf := map[string]interface{}{
		"request_id":      appcontext.GetRequestId(ctx),
		"user_agent":      appcontext.GetUserAgent(ctx),
		"user_id":         appcontext.GetUserId(ctx),
		"service_version": appcontext.GetServiceVersion(ctx),
		"time_elapsed":    timeElapsed,
	}

	if apprespcode > 0 {
		cf["app_resp_code"] = apprespcode
	}

	if appErrMsg != "" {
		cf["app_err_msg"] = appErrMsg
	}

	return cf
}
