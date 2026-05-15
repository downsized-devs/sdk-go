# `operator` — bitwise & ternary helpers

`import "github.com/downsized-devs/sdk-go/operator"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Tiny package. Used by [`codes`](../codes), [`errors`](../errors), [`audit`](../audit), [`tracker`](../tracker).

## Features

- `Ternary[T comparable](cond, a, b) T` — generic ternary, replaces `if`-statements inline.
- `CheckBitOnPosition(number, position int) bool` — bit-flag test.

## Installation

```bash
go get github.com/downsized-devs/sdk-go/operator
```

## Quick Start

```go
import "github.com/downsized-devs/sdk-go/operator"

label := operator.Ternary(ok, "yes", "no")
flag  := operator.CheckBitOnPosition(0b1010, 2) // true (2nd bit set)
```

## API Reference

| Symbol | Signature |
|---|---|
| `Ternary` | `[T comparable](condition bool, a, b T) T` |
| `CheckBitOnPosition` | `(number, position int) bool` (position is 1-indexed) |

## Dependencies

stdlib only.

## Testing

```bash
go test ./operator/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md).

## Related Packages

- [`codes`](../codes), [`errors`](../errors), [`audit`](../audit), [`tracker`](../tracker) — consumers.
