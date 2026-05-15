# `redis` — Redis client with distributed locking

`import "github.com/downsized-devs/sdk-go/redis"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Wraps [`go-redis/redis`](https://github.com/go-redis/redis) and [`bsm/redislock`](https://github.com/bsm/redislock). Supports TLS, exposes a CRC16 helper for cluster slot routing.

## Features

- `Get` / `SetEX` with default TTL fallback
- Distributed lock: `Lock` / `LockRelease`
- `Del`, `FlushAll`, `FlushAllAsync`, `FlushDB`, `FlushDBAsync`, `Ping`
- Optional TLS with private CA / mTLS
- `CRC16(s)` for cluster slot hashing

## Installation

```bash
go get github.com/downsized-devs/sdk-go/redis
```

## Quick Start

```go
import (
    "context"
    "time"

    "github.com/downsized-devs/sdk-go/logger"
    "github.com/downsized-devs/sdk-go/redis"
)

log := logger.Init(logger.Config{Level: "info"})
rdb := redis.Init(redis.Config{
    Address:    "localhost:6379",
    DefaultTTL: 5 * time.Minute,
}, log)

ctx := context.Background()
_ = rdb.SetEX(ctx, "hello", "world", time.Minute)
val, _ := rdb.Get(ctx, "hello")
```

## API Reference

### Construction

```go
func Init(cfg Config, log logger.Interface) Interface
```

### `Interface`

| Method | Signature | Notes |
|---|---|---|
| `Get` | `(ctx, key string) (string, error)` | Returns `redis.Nil` on miss. |
| `SetEX` | `(ctx, key, val string, ttl time.Duration) error` | `0` ttl → `Config.DefaultTTL`. |
| `Del` | `(ctx, key string) error` | |
| `Lock` | `(ctx, key string, expTime time.Duration) (*redislock.Lock, error)` | `redis.ErrNotObtained` if contended. |
| `LockRelease` | `(ctx, lock *redislock.Lock) error` | Pair every `Lock`. |
| `FlushAll`/`FlushAllAsync` | `(ctx) error` | Wipes *every* database. |
| `FlushDB`/`FlushDBAsync` | `(ctx) error` | Wipes selected database. |
| `Ping` | `(ctx) error` | Liveness check. |
| `GetDefaultTTL` | `(ctx) time.Duration` | |

### Top-level helpers

| Symbol | Purpose |
|---|---|
| `Nil` | Sentinel for `Get` miss; same as `go-redis/redis.Nil`. |
| `ErrNotObtained` | Returned by `Lock` when contended. |
| `CRC16(s string) uint16` | CRC16-XMODEM, used for cluster slot routing. |

## Configuration

| Field | Type | Required | Default | Description |
|---|---|---|---|---|
| `Address` | `string` | yes | — | `host:port`. |
| `Password` | `string` | no | `""` | |
| `DB` | `int` | no | `0` | DB index. |
| `DefaultTTL` | `time.Duration` | no | `0` | Used when `SetEX` ttl is `0`. |
| `TLS.Enabled` | `bool` | no | `false` | |
| `TLS.CA`/`Cert`/`Key` | `string` | conditional | — | Required for private CA / mTLS. |

## Examples

### Cache-aside

```go
val, err := rdb.Get(ctx, "user:"+id)
if errors.Is(err, redis.Nil) {
    user, err := repo.LoadUser(ctx, id)
    if err != nil { return nil, err }
    _ = rdb.SetEX(ctx, "user:"+id, user.Marshal(), 5*time.Minute)
    return user, nil
}
```

### Distributed lock around a cron task

```go
lock, err := rdb.Lock(ctx, "cron:nightly-rollup", 10*time.Minute)
if errors.Is(err, redis.ErrNotObtained) {
    return nil // another instance has it
}
if err != nil { return err }
defer rdb.LockRelease(ctx, lock)

return doRollup(ctx)
```

## Error Handling

| Error | Action |
|---|---|
| `redis.Nil` | Treat as miss; reload from origin. |
| `redis.ErrNotObtained` | Skip work or back off. |
| Coded errors | Inspect with `errors.GetCode(err)`. |

## Dependencies

- **Internal:** [`codes`](../codes), [`errors`](../errors), [`logger`](../logger)
- **External:** `github.com/go-redis/redis/v8`, `github.com/bsm/redislock`

## Testing

```bash
go test ./redis/...
```

Two test files; unit test runs without a live Redis.

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). New multi-key methods should be verified with `CRC16` to ensure all keys hash to the same cluster slot.

## Related Packages

- [`sql`](../sql) — primary store; pair for cache-aside.
- [`scheduler`](../scheduler) — uses Redis locks to coordinate cron across replicas.
- [`instrument`](../instrument) — register Redis hit/miss metrics.
