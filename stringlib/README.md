# `stringlib` — random string helper

`import "github.com/downsized-devs/sdk-go/stringlib"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Single-function package.

## Features

- `RandStringBytes(n int) string` — random alphanumeric string of length `n`.

## Installation

```bash
go get github.com/downsized-devs/sdk-go/stringlib
```

## Quick Start

```go
import "github.com/downsized-devs/sdk-go/stringlib"

_ = stringlib.RandStringBytes(16) // 16-char random string
```

## API Reference

| Symbol | Signature |
|---|---|
| `RandStringBytes` | `(n int) string` |

## Dependencies

stdlib only.

## Testing

```bash
go test ./stringlib/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). Prefer `num.RandomString` (which overrides `now` for deterministic tests) for new code.

## Related Packages

- [`num.RandomString`](../num) — newer equivalent with mockable time source.
