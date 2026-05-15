# `tracker` — Prometheus push gateway + webhook

`import "github.com/downsized-devs/sdk-go/tracker"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Pushes metrics to a Prometheus push gateway for short-lived jobs, and posts arbitrary JSON to webhook endpoints.

## Features

- `Push(Options)` — push a metric to the push gateway.
- `PushWebhook(WebhookOptions)` — POST a JSON payload to an HTTP webhook.

## Installation

```bash
go get github.com/downsized-devs/sdk-go/tracker
```

## Quick Start

```go
t := tracker.Init(tracker.Config{
    PushGateway: "http://pgw.local:9091",
}, log)

_ = t.Push(ctx, tracker.Options{
    Job:    "nightly-rollup",
    Metric: "rollup_rows_processed",
    Value:  12345,
    Labels: map[string]string{ "env": "prod" },
})
```

## API Reference

| Symbol | Signature |
|---|---|
| `Init` | `func Init(cfg Config, log logger.Interface) Interface` |
| `Interface.Push` | `(ctx, Options) error` |
| `Interface.PushWebhook` | `(ctx, WebhookOptions) error` |
| `Options` | `{ Job, Metric string; Value float64; Labels map[string]string }` |
| `WebhookOptions` | `{ URL string; Headers http.Header; Body any }` |

## Configuration

| Field | Description |
|---|---|
| `PushGateway` | Push-gateway URL. |
| Webhook defaults | per-call overrides recommended. |

## Error Handling

Errors are wrapped with [`codes`](../codes) third-party codes.

## Dependencies

- **Internal:** [`codes`](../codes), [`errors`](../errors), [`logger`](../logger), [`operator`](../operator)
- **External:** `github.com/prometheus/client_golang/prometheus`, `.../prometheus/push`, `github.com/prometheus/common/expfmt`

## Testing

```bash
go test ./tracker/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md).

## Related Packages

- [`instrument`](../instrument) — Prometheus *pull* (the inverse of this package).
- [`slack`](../slack) — alternative alerting channel.
