# `codes` — central registry of error & success codes

`import "github.com/downsized-devs/sdk-go/codes"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Single source of truth for the numeric codes that flow through [`errors`](../errors) and into HTTP response bodies. Codes are uint32s grouped into reserved ranges so SDK packages can add codes without colliding.

## Features

- `Code` type (alias for `uint32`)
- `AppMessage` map type for bilingual messages
- `DisplayMessage` type used by HTTP responses
- `Compile(c Code, lang string) DisplayMessage` — look up the human-readable text
- Reserved code ranges enforced by convention (see `codes/codes.go`)

## Code ranges

```
10   – 99    Success codes
1000 – 1299  Common / generic errors
1300 – 1399  SQL errors
1400 – 1499  NoSQL errors
1500 – 1599  Third-party / client errors
1600 – 1699  File I/O errors
1700 – 1799  Auth errors
1800 – 1899  (reserved)
```

When adding a new SDK group, append the next free 100-block at the bottom of `codes.go` and update this table.

## Installation

```bash
go get github.com/downsized-devs/sdk-go/codes
```

## Quick Start

```go
import (
    "github.com/downsized-devs/sdk-go/codes"
    "github.com/downsized-devs/sdk-go/errors"
)

err := errors.NewWithCode(codes.CodeBadRequest, "missing field: %s", "email")
msg := errors.Compile(err, "en")
// msg.StatusCode == 400, msg.Body == "missing field: email"
```

## API Reference

| Symbol | Purpose |
|---|---|
| `Code` | `uint32` alias. |
| `AppMessage` | `map[Code]Message` — registry of messages per code. |
| `DisplayMessage` | `{StatusCode int, Title, Body string}` used in HTTP responses. |
| `Compile(c, lang) DisplayMessage` | Build a display message in the chosen language. |
| `NoCode` | Sentinel meaning "no code attached". |
| `ErrorMessages`, `ApplicationMessages` | Pre-populated maps for SDK codes. |

## Examples

### Adding a new code

```go
// codes/codes.go
const (
    CodeMyNewThing Code = 2000 // first code in a new reserved block
)

// codes/code_messages.go — add the human messages
var newThingMessages = AppMessage{
    CodeMyNewThing: {
        StatusCode: http.StatusBadRequest,
        TitleEn:    "My New Thing",
        BodyEn:     "...",
        TitleId:    "...",
        BodyId:     "...",
    },
}
```

Register the map by adding it to the package's initial set so `Compile` can find it.

## Error Handling

`Compile` falls back to a generic "internal error" `DisplayMessage` when the code is unknown. Never panic on lookup.

## Dependencies

- **Internal:** [`language`](../language), [`operator`](../operator)
- **External:** stdlib only

## Testing

```bash
go test ./codes/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). **Never re-number an existing code** — downstream services and clients persist them.

## Related Packages

- [`errors`](../errors) — primary consumer.
- [`language`](../language) — locale lookups.
