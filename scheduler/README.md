# `scheduler` — gocron v2 wrapper for periodic jobs

`import "github.com/downsized-devs/sdk-go/scheduler"`

**Stability:** Beta — no tests in package. See [STABILITY.md](../STABILITY.md).

Thin facade over [`gocron/v2`](https://github.com/go-co-op/gocron) that lets you register duration-based, daily, weekly, or monthly jobs from configuration.

## Features

- `Register(JobOption)` — add a job
- `Start()` / `Shutdown()` lifecycle
- Job-type constants for duration, daily, weekly, monthly
- Pairs naturally with [`redis`](../redis) `Lock` for multi-replica coordination

## Installation

```bash
go get github.com/downsized-devs/sdk-go/scheduler
```

## Quick Start

```go
import (
    "context"
    "time"

    "github.com/downsized-devs/sdk-go/logger"
    "github.com/downsized-devs/sdk-go/scheduler"
)

log := logger.Init(logger.Config{Level: "info"})
sch := scheduler.New(scheduler.Config{}, log)

sch.Register(scheduler.JobOption{
    Name:      "heartbeat",
    Type:      scheduler.JobTypeDuration,
    Duration:  10 * time.Second,
    Task: func(ctx context.Context) {
        log.Info(ctx, "tick")
    },
})

sch.Start()
defer sch.Shutdown()
```

## API Reference

| Symbol | Signature |
|---|---|
| `New` | `func New(cfg Config, log logger.Interface) Interface` |
| `Interface.Start` | `() error` |
| `Interface.Shutdown` | `() error` |
| `Interface.Register` | `(JobOption) error` |
| `JobOption` | `Name string; Type string; Duration time.Duration; Daily/Weekly/Monthly time fields; Task func(ctx)` |

Job-type constants are exported from `scheduler.go`; check the source for the current set (duration, daily, weekly, monthly at minimum).

## Configuration

| Field | Purpose |
|---|---|
| (See `Config` in `scheduler.go`) | Timezone, monitoring hooks. |

## Examples

### Coordinate cron across replicas

```go
sch.Register(scheduler.JobOption{
    Name:     "nightly-rollup",
    Type:     scheduler.JobTypeDaily,
    Daily:    "02:00",
    Task: func(ctx context.Context) {
        lock, err := rdb.Lock(ctx, "cron:rollup", 30*time.Minute)
        if errors.Is(err, redis.ErrNotObtained) { return }
        if err != nil { log.Error(ctx, err); return }
        defer rdb.LockRelease(ctx, lock)
        doRollup(ctx)
    },
})
```

## Error Handling

`Register` returns an error for unknown job types or malformed time strings. Always check it.

## Dependencies

- **Internal:** [`logger`](../logger)
- **External:** `github.com/go-co-op/gocron/v2`

## Testing

The package has no tests of its own at the time of writing. When adding tests, use `gocron`'s built-in `Run` semantics rather than wall-clock sleeps.

```bash
go test ./scheduler/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). Add tests as part of any non-trivial change — this package is Beta until coverage lands.

## Related Packages

- [`redis`](../redis) — locks for cron coordination.
- [`logger`](../logger) — required at `New` time.
- [`instrument`](../instrument) — already exposes `SchedulerRunningTimer` / `SchedulerRunningCounter` for Prometheus.
