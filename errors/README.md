# `errors` — error wrapping with codes and stack traces

`import "github.com/downsized-devs/sdk-go/errors"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Drop-in replacement for the stdlib `errors` package. Attaches numeric [`codes`](../codes), localised messages from [`language`](../language), and caller information so log lines and HTTP responses can be rendered uniformly.

## Features

- `NewWithCode` / `WrapWithCode` — build a coded error or wrap an existing one
- `GetCode(err) codes.Code` — extract the code from anywhere in the chain
- `GetCaller(err)` — file:line of the original wrap site
- `Compile(err, language string)` — render a `DisplayMessage` for HTTP response bodies
- `Is` / `As` — fully compatible with `errors.Is` / `errors.As` semantics
- Stack frame capture (see `errors_stacktrace.go`)

## Installation

```bash
go get github.com/downsized-devs/sdk-go/errors
```

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/downsized-devs/sdk-go/codes"
    "github.com/downsized-devs/sdk-go/errors"
)

func loadUser(id string) error {
    if id == "" {
        return errors.NewWithCode(codes.CodeBadRequest, "user id is empty")
    }
    return nil
}

func main() {
    if err := loadUser(""); err != nil {
        fmt.Println("code:", errors.GetCode(err))
        fmt.Println("caller:", errors.GetCaller(err))
    }
}
```

## API Reference

| Symbol | Purpose |
|---|---|
| `NewWithCode(code, format, args...)` | Construct a new coded error. |
| `WrapWithCode(err, code, format, args...)` | Wrap a lower-level error with a code. |
| `GetCode(err) codes.Code` | First code found walking the chain. |
| `GetCaller(err) string` | File:line where the error was first wrapped. |
| `Compile(err, lang) codes.DisplayMessage` | Build a HTTP-ready message in the chosen language. |
| `Is(err, target)` / `As(err, target)` | Standard chain inspection. |
| `App` | Concrete error type carrying code + message + caller + cause. |

## Examples

### Wrap a SQL error and dispatch by code

```go
row, err := db.Leader(ctx).Get(ctx, &u, "SELECT ... WHERE id = ?", id)
if err != nil {
    return errors.WrapWithCode(err, codes.CodeSQLRecordDoesNotExist, "load user %s", id)
}

// Later, at the HTTP layer:
switch errors.GetCode(err) {
case codes.CodeSQLRecordDoesNotExist:
    http.Error(w, "not found", http.StatusNotFound)
default:
    http.Error(w, "internal error", http.StatusInternalServerError)
}
```

### Render a localised response body

```go
msg := errors.Compile(err, "en")
// msg.StatusCode, msg.Title, msg.Body — ready for JSON encoding
```

## Error Handling

The package *produces* errors; it does not return any of its own from public functions. `GetCode` returns `codes.NoCode` if the chain has none.

## Dependencies

- **Internal:** [`codes`](../codes), [`language`](../language), [`operator`](../operator)
- **External:** stdlib only

## Testing

```bash
go test ./errors/...
```

Two test files cover the App type, chain walking, and stack-trace capture.

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). This package is critical-path — 14 sibling packages depend on it.

## Related Packages

- [`codes`](../codes) — the code registry.
- [`language`](../language) — locale lookup used by `Compile`.
- [`logger`](../logger) — auto-extracts code + caller when logging an `App` error.
