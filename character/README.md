# `character` — string casing & password-strength helpers

`import "github.com/downsized-devs/sdk-go/character"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

## Features

- `CapitalizeFirstCharacter(s)` — title-case first rune.
- `IsStrongCharCombination(s)` — regex-based password-strength check.

## Installation

```bash
go get github.com/downsized-devs/sdk-go/character
```

## Quick Start

```go
import "github.com/downsized-devs/sdk-go/character"

_ = character.CapitalizeFirstCharacter("hello")  // "Hello"
_ = character.IsStrongCharCombination("P@ssw0rd") // true/false
```

## API Reference

| Symbol | Signature |
|---|---|
| `CapitalizeFirstCharacter` | `(s string) string` |
| `IsStrongCharCombination` | `(s string) bool` |

## Dependencies

- **External:** `golang.org/x/text/cases`, `golang.org/x/text/language`

## Testing

```bash
go test ./character/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md).

## Related Packages

- [`security`](../security) — hashing once a password passes strength checks.
- [`stringlib`](../stringlib) — random string generation.
