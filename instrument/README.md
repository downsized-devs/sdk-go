# `instrument` — Prometheus metrics for HTTP, DB, and scheduler

`import "github.com/downsized-devs/sdk-go/instrument"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Wraps `prometheus/client_golang` with pre-built collectors for the SDK's primary surfaces (HTTP, SQL pools, scheduled jobs).

## Features

- `MetricsHandler()` — drop-in `http.Handler` for `/metrics`.
- `HTTPRequestTimer`, `HTTPRequestCounter`, `HTTPResponseStatusCounter`.
- `RegisterDBStats`, `DatabaseQueryTimer` — used by [`sql`](../sql).
- `SchedulerRunningCounter`, `SchedulerRunningTimer` — used by [`scheduler`](../scheduler).
- `IsEnabled` — quick gate for callers that should no-op when metrics are off.

## Installation

```bash
go get github.com/downsized-devs/sdk-go/instrument
```

## Quick Start

```go
m := instrument.Init(instrument.Config{Enabled: true})
http.Handle("/metrics", m.MetricsHandler())
```

## API Reference

| Symbol | Signature |
|---|---|
| `Init` | `func Init(cfg Config) Interface` |
| `Interface.IsEnabled` | `() bool` |
| `Interface.MetricsHandler` | `() http.Handler` |
| `Interface.HTTPRequestTimer` | `(method, route string) *prometheus.HistogramVec` |
| `Interface.HTTPRequestCounter` | `(method, route string) prometheus.Counter` |
| `Interface.HTTPResponseStatusCounter` | `(method, route, status string) prometheus.Counter` |
| `Interface.RegisterDBStats` | `(name string, db DBStats)` |
| `Interface.DatabaseQueryTimer` | `(name, op string) prometheus.Observer` |
| `Interface.SchedulerRunningCounter` | `(job string) prometheus.Counter` |
| `Interface.SchedulerRunningTimer` | `(job string) prometheus.Observer` |

## Configuration

| Field | Description |
|---|---|
| `Enabled` | Master on/off switch. |
| `Namespace`, `Subsystem` | Optional Prometheus label prefixes. |

## Dependencies

- **External:** `github.com/prometheus/client_golang`

## Testing

```bash
go test ./instrument/...
```

Two test files cover handler, counters, and timer registration.

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). Metric names are public — don't rename existing ones (breaks dashboards and alerts).

## Related Packages

- [`sql`](../sql) — consumer for `RegisterDBStats` and `DatabaseQueryTimer`.
- [`scheduler`](../scheduler) — consumer for scheduler timers/counters.
- [`tracker`](../tracker) — Prometheus *push* gateway (this package exposes *pull*).
