# `checker` — generic validators

`import "github.com/downsized-devs/sdk-go/checker"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Tiny library of validation helpers. Uses generics; no external deps.

## Features

- `ArrayInt64Contains`, `ArrayContains` — membership tests.
- `ArrayDeduplicate` — remove duplicates while preserving order.
- `IsEmail`, `IsPhoneNumber` — regex-based format checks.

## Installation

```bash
go get github.com/downsized-devs/sdk-go/checker
```

## Quick Start

```go
import "github.com/downsized-devs/sdk-go/checker"

ok := checker.ArrayContains([]string{"a", "b"}, "a") // true
clean := checker.ArrayDeduplicate([]int{1, 1, 2})    // [1, 2]
checker.IsEmail("alice@example.com")                 // true
```

## API Reference

| Symbol | Signature |
|---|---|
| `ArrayContains` | `[T comparable](arr []T, target T) bool` |
| `ArrayInt64Contains` | `(arr []int64, target int64) bool` |
| `ArrayDeduplicate` | `[T AllowedDataType](arr []T) []T` |
| `IsEmail` | `(s string) bool` |
| `IsPhoneNumber` | `(s string) bool` |
| `AllowedDataType` | constraint interface (comparable types). |

## Dependencies

stdlib only.

## Testing

```bash
go test ./checker/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md).

## Related Packages

- [`ratelimiter`](../ratelimiter) — uses `checker` to validate paths.
