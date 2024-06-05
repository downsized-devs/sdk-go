package audit

import (
	"context"
	"os"
	"sync"

	"github.com/downsized-devs/sdk-go/appcontext"
	"github.com/downsized-devs/sdk-go/auth"
	"github.com/downsized-devs/sdk-go/operator"
	"github.com/rs/zerolog"
)

var once = sync.Once{}

type Interface interface {
	Capture(ctx context.Context)
	Record(ctx context.Context, log Collection)
}

type audit struct {
	log  zerolog.Logger
	auth auth.Interface
}

func Init(auth auth.Interface) Interface {
	var zeroLogging zerolog.Logger

	// Initialize a new Zerolog object to handle the split log file
	once.Do(func() {
		zeroLogging = zerolog.New(os.Stdout).
			With().
			Timestamp().
			Logger()
	})

	return &audit{log: zeroLogging, auth: auth}
}

func (a *audit) Capture(ctx context.Context) {
	// Only audit if contains the event name
	if appcontext.GetEventName(ctx) != "" {
		a.log.Log().Fields(a.getHttpFields(ctx)).Msg("")
	}
}

func (a *audit) Record(ctx context.Context, log Collection) {
	a.log.Log().Fields(a.getDomainFields(ctx, log)).Msg("")
}

func (a *audit) getHttpFields(ctx context.Context) map[string]interface{} {
	user, _ := a.auth.GetUserAuthInfo(ctx)

	return map[string]interface{}{
		"log_type":          logType,
		"event_name":        appcontext.GetEventName(ctx),
		"event_description": appcontext.GetEventDescription(ctx),
		"event_type":        eventTypeHttp,
		"user_id":           user.User.ID,
		"user_company_id":   user.User.CompanyID,
		"user_role_id":      user.User.RoleID,
		"user_name":         user.User.Name,
		"request_id":        appcontext.GetRequestId(ctx),
		"request_device":    appcontext.GetDeviceType(ctx),
		"request_ip":        appcontext.GetRequestIP(ctx),
		"request_agent":     appcontext.GetUserAgent(ctx),
		"request_uri":       appcontext.GetRequestURI(ctx),
		"request_method":    appcontext.GetRequestMethod(ctx),
		"request_query":     appcontext.GetRequestQuery(ctx),
		"request_body":      appcontext.GetRequestBody(ctx),
		"response_code":     appcontext.GetResponseHttpCode(ctx),
	}
}

func (a *audit) getDomainFields(ctx context.Context, log Collection) map[string]interface{} {
	user, _ := a.auth.GetUserAuthInfo(ctx)
	status := operator.Ternary(log.Error != nil, false, true)

	fields := map[string]interface{}{
		"log_type":          logType,
		"event_name":        log.EventName,
		"event_description": log.EventDescription,
		"event_type":        eventTypeDomain,
		"user_id":           user.User.ID,
		"user_company_id":   user.User.CompanyID,
		"user_role_id":      user.User.RoleID,
		"user_name":         user.User.Name,
		"request_id":        appcontext.GetRequestId(ctx),
		"request_device":    appcontext.GetDeviceType(ctx),
		"status":            status,
	}

	fields["data"] = inputs{
		Insert: log.InsertParam,
		Select: log.SelectParam,
		Update: log.UpdateParam,
	}

	if !status {
		fields["msg_error"] = log.Error
	}

	return fields
}
