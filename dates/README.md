# `dates` — date arithmetic helpers

`import "github.com/downsized-devs/sdk-go/dates"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Tiny package: day-level differences between two `time.Time` values.

## Features

- `Difference(a, b time.Time) int64` — whole-day difference (b − a).

## Installation

```bash
go get github.com/downsized-devs/sdk-go/dates
```

## Quick Start

```go
import "github.com/downsized-devs/sdk-go/dates"

days := dates.Difference(start, end) // e.g. 7
```

## API Reference

| Symbol | Signature |
|---|---|
| `Difference` | `(a, b time.Time) int64` |

## Dependencies

stdlib only.

## Testing

```bash
go test ./dates/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md).

## Related Packages

- [`clock`](../clock) — timezone-aware "now" + month boundaries.
