# `clock` — timezone-aware clock with mockable `Now`

`import "github.com/downsized-devs/sdk-go/clock"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Inject this instead of calling `time.Now()` directly — tests can swap in a fixed clock.

## Features

- `GetCurrentTime` — current time in the configured location.
- `AddTime`, `SubstractTime` — duration arithmetic.
- `GetTimeInLocation` — convert across timezones.
- `GetFirstDayOfTheMonth`, `GetLastDayOfTheMonth`.
- `ConvertFromString` / `ConvertToString` — format-aware string ↔ time.
- Package-level `Now` variable for test override.

## Installation

```bash
go get github.com/downsized-devs/sdk-go/clock
```

## Quick Start

```go
c := clock.Init(clock.Location{Name: "Asia/Jakarta"})
now := c.GetCurrentTime()
first := c.GetFirstDayOfTheMonth(now)
```

## API Reference

| Symbol | Signature |
|---|---|
| `Init` | `func Init(loc Location) Interface` |
| `Interface.GetCurrentTime` | `() time.Time` |
| `Interface.AddTime` | `(t time.Time, d time.Duration) time.Time` |
| `Interface.SubstractTime` | `(t time.Time, d time.Duration) time.Time` |
| `Interface.GetTimeInLocation` | `(t time.Time) time.Time` |
| `Interface.GetFirstDayOfTheMonth` | `(t time.Time) time.Time` |
| `Interface.GetLastDayOfTheMonth` | `(t time.Time) time.Time` |
| `Interface.ConvertFromString` | `(s, layout string) (time.Time, error)` |
| `Interface.ConvertToString` | `(t time.Time, layout string) string` |
| `Now` | `var Now = time.Now` — override in tests. |

## Examples

### Test override

```go
clock.Now = func() time.Time { return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) }
defer func() { clock.Now = time.Now }()
```

## Dependencies

stdlib only.

## Testing

```bash
go test ./clock/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md).

## Related Packages

- [`dates`](../dates) — day-level differences between times.
- [`scheduler`](../scheduler) — natural pairing for cron timing.
