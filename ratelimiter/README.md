# `ratelimiter` — Gin rate-limiting middleware

`import "github.com/downsized-devs/sdk-go/ratelimiter"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Per-path rate limiting middleware for [Gin](https://github.com/gin-gonic/gin), backed by [`ulule/limiter`](https://github.com/ulule/limiter)'s in-memory store.

## Features

- `Init` returns a `GinMiddleware` ready to attach.
- Per-path overrides via `ConfigPath`.
- In-memory store (replace with Redis store if you need cross-replica limits — not done here).

## Installation

```bash
go get github.com/downsized-devs/sdk-go/ratelimiter
```

## Quick Start

```go
rl := ratelimiter.Init(ratelimiter.Config{
    Enabled: true,
    Default: "100-S", // 100 requests per second
    Paths: []ratelimiter.ConfigPath{
        { Path: "/api/login", Rate: "5-M" }, // 5/minute on login
    },
}, log, appcontext, checker)

r := gin.New()
r.Use(rl.Limiter())
```

## API Reference

| Symbol | Signature |
|---|---|
| `Init` | `func Init(cfg Config, ...deps) Interface` |
| `Interface.Limiter` | `() gin.HandlerFunc` |
| `Config` | `{ Enabled bool; Default string; Paths []ConfigPath }` |
| `ConfigPath` | `{ Path string; Rate string }` |

Rate strings follow ulule format: `<count>-<unit>` where unit is `S`, `M`, `H`, or `D` (e.g. `10-S`, `1000-H`).

## Configuration

| Field | Description |
|---|---|
| `Enabled` | Master toggle. |
| `Default` | Fallback rate for paths not explicitly listed. |
| `Paths` | Overrides per request path. |

## Error Handling

Limited requests get `429 Too Many Requests`. The middleware does not return errors to callers.

## Dependencies

- **Internal:** [`appcontext`](../appcontext), [`checker`](../checker), [`codes`](../codes), [`errors`](../errors), [`logger`](../logger)
- **External:** `github.com/gin-gonic/gin`, `github.com/ulule/limiter/v3`, `.../drivers/middleware/gin`, `.../drivers/store/memory`

## Testing

```bash
go test ./ratelimiter/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). For multi-replica deployments, contribute a Redis-store variant.

## Related Packages

- [`auth`](../auth) — typically paired so unauthenticated traffic is rate-limited.
- [`logger`](../logger) — required at `Init` time.
