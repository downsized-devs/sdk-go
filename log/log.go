package log

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/downsized-devs/sdk-go/appcontext"
	"github.com/downsized-devs/sdk-go/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var once = sync.Once{}
var now = time.Now

type contextKey string

type Interface interface {
	// TODO add Debugf
	Trace(ctx context.Context, obj any)
	Debug(ctx context.Context, obj any)
	Info(ctx context.Context, obj any)
	Warn(ctx context.Context, obj any)
	Error(ctx context.Context, obj any)
	Fatal(ctx context.Context, obj any)
	Panic(obj any)
}

type Config struct {
	Level string
}

type logger struct {
	log zerolog.Logger
}

func DefaultLogger() Interface {
	return &logger{
		log: zerolog.New(os.Stdout).
			With().
			Timestamp().
			CallerWithSkipFrameCount(3). //Hard code to 3 for now.
			Logger().
			Level(zerolog.DebugLevel),
	}
}

func Init(cfg Config) Interface {
	var zeroLogging zerolog.Logger
	once.Do(func() {
		level, err := zerolog.ParseLevel(cfg.Level)
		if err != nil {
			log.Fatal().Msg(fmt.Sprintf("failed to parse error level with err: %v", err))
		}

		zeroLogging = zerolog.New(os.Stdout).
			With().
			Timestamp().
			CallerWithSkipFrameCount(3). //Hard code to 3 for now.
			Logger().
			Level(level)
	})

	return &logger{log: zeroLogging}
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
		"stracktrace": errStack,
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
