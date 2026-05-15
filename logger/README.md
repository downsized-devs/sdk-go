# `logger` — structured, context-aware logging

`import "github.com/downsized-devs/sdk-go/logger"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Thin wrapper over [`rs/zerolog`](https://github.com/rs/zerolog) that pulls request metadata from [`appcontext`](../appcontext) and enriches error logs via [`errors`](../errors).

## Features

- Eight log levels: `Trace`, `Debug`, `Debugf`, `Info`, `Warn`, `Error`, `Fatal`, `Panic`
- Automatic context-field extraction (request ID, user ID, service version, etc.)
- Stack-trace + caller info enrichment on `Panic`
- `DefaultLogger()` bootstrap for tests and pre-config code paths
- Single `Interface` — trivial to mock with `gomock`

## Installation

```bash
go get github.com/downsized-devs/sdk-go/logger
```

## Quick Start

```go
package main

import (
    "context"

    "github.com/downsized-devs/sdk-go/logger"
)

func main() {
    log := logger.Init(logger.Config{Level: "info"})
    ctx := context.Background()

    log.Info(ctx, "service started")
    log.Debugf(ctx, "loaded %d feature flags", 42)
}
```

## API Reference

### Construction

| Symbol | Signature | Notes |
|---|---|---|
| `Init` | `func Init(cfg Config) Interface` | Primary constructor. |
| `DefaultLogger` | `func DefaultLogger() Interface` | Pre-config / test bootstrap (stderr, info level). |

### `Interface`

| Method | Signature |
|---|---|
| `Trace` | `(ctx, obj any)` |
| `Debug` | `(ctx, obj any)` |
| `Debugf` | `(ctx, format string, args ...any)` |
| `Info` | `(ctx, obj any)` |
| `Warn` | `(ctx, obj any)` |
| `Error` | `(ctx, obj any)` |
| `Fatal` | `(ctx, obj any)` — logs then `os.Exit(1)` |
| `Panic` | `(obj any)` — logs with stack trace then panics |

## Configuration

| Field | Type | Default | Description |
|---|---|---|---|
| `Level` | `string` | `info` | One of `trace`, `debug`, `info`, `warn`, `error`, `fatal`. Invalid values fall back to `info`. |

## Examples

### HTTP middleware that propagates request ID

```go
log := logger.Init(logger.Config{Level: "info"})

func handler(w http.ResponseWriter, r *http.Request) {
    ctx := appcontext.SetRequestId(r.Context(), uuid.NewString())
    log.Info(ctx, "received request") // request_id field auto-attached
}
```

### Log a wrapped error with its code

```go
if err := repo.Save(ctx, user); err != nil {
    err = errors.WrapWithCode(err, codes.CodeSQLTxRollback, "save user")
    log.Error(ctx, err) // includes stack trace + code
}
```

### Use the default logger in a `TestMain`

```go
func TestMain(m *testing.M) {
    log := logger.DefaultLogger()
    log.Info(context.Background(), "starting test suite")
    os.Exit(m.Run())
}
```

## Error Handling

Methods do not return errors — `logger` writes and continues. `Panic` deliberately re-panics so the caller's recover policy applies.

## Dependencies

- **Internal:** [`appcontext`](../appcontext), [`errors`](../errors)
- **External:** `github.com/rs/zerolog`

## Testing

```bash
go test ./logger/...
```

Two test files cover every level plus context-field extraction.

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). Adding a method to `Interface` is a breaking change — coordinate before merging.

## Related Packages

- [`appcontext`](../appcontext) — populates the fields this logger reads.
- [`errors`](../errors) — supplies the codes/stack traces shown in log lines.
- [`instrument`](../instrument) — Prometheus metrics; logs say *what*, metrics say *how often*.
- [`audit`](../audit) — durable business events vs ephemeral log lines.
