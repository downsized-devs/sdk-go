package tracker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/downsized-devs/sdk-go/codes"
	"github.com/downsized-devs/sdk-go/errors"
	"github.com/downsized-devs/sdk-go/logger"
	"github.com/downsized-devs/sdk-go/operator"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/prometheus/common/expfmt"
)

const (
	// DEFAULT VALUE
	defaultJobName        string        = "tracking_metrics"
	defaultURL            string        = "localhost"
	defaultPort           string        = "9091"
	defaultTimeout        time.Duration = 5 * time.Second
	defaultWebhookTimeout time.Duration = 10 * time.Second
)

type Interface interface {
	Push(ctx context.Context, trackingName string, labels map[string]string) error
	PushWebhook(ctx context.Context, payload []byte, headers map[string]string) error
}

type Options struct {
	Enabled bool
	URL     string
	Port    string
	JobName string
	Timeout time.Duration
	Webhook WebhookOptions
}

type WebhookOptions struct {
	Enabled bool
	URL     string
	Timeout time.Duration
}

type tracker struct {
	opt           Options
	log           logger.Interface
	webhookClient *http.Client
}

func Init(opt Options, log logger.Interface) Interface {

	return &tracker{
		opt: opt,
		log: log,
		webhookClient: &http.Client{
			Timeout: operator.Ternary(opt.Webhook.Timeout > 0, opt.Webhook.Timeout, defaultWebhookTimeout),
		},
	}
}

func (t *tracker) Push(ctx context.Context, trackingName string, labels map[string]string) error {
	if !t.opt.Enabled {
		return nil
	}

	opsProcessed := prometheus.NewCounter(prometheus.CounterOpts{
		Name:        fmt.Sprintf("tracking_%s", trackingName),
		Help:        "Tracking Info",
		ConstLabels: labels,
	})

	opsProcessed.Inc()

	conn := fmt.Sprintf("%s:%s", t.opt.URL, t.opt.Port)
	if conn == "" {
		conn = fmt.Sprintf("%s:%s", defaultURL, defaultPort)
	}

	jobName := operator.Ternary(t.opt.JobName != "", t.opt.JobName, defaultJobName)
	timeout := operator.Ternary(t.opt.Timeout > 0, t.opt.Timeout, defaultTimeout)

	if err := push.New(conn, jobName).
		Client(&http.Client{Timeout: timeout}).
		Collector(opsProcessed).
		Format(expfmt.FmtText).
		Push(); err != nil {
		return err
	}

	return nil
}

func (t *tracker) PushWebhook(ctx context.Context, payload []byte, headers map[string]string) error {
	if !t.opt.Webhook.Enabled {
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.opt.Webhook.URL, bytes.NewBuffer(payload))
	if err != nil {
		return errors.NewWithCode(codes.CodeErrorHttpNewRequest, err.Error())
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := t.webhookClient.Do(req)
	if err != nil {
		return errors.NewWithCode(codes.CodeErrorHttpDo, err.Error())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.NewWithCode(codes.CodeErrorIoutilReadAll, err.Error())
	}

	if resp.StatusCode != 200 {
		return errors.NewWithCode(codes.CodeErrorHttpDo, fmt.Sprintf("send event webhook error: %s", string(body)))
	}

	return nil
}
