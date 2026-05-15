# `audit` — audit-trail event capture

`import "github.com/downsized-devs/sdk-go/audit"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Captures domain events with the request/user context already on `context.Context`. Persistent log stream separate from operational [`logger`](../logger) lines.

## Features

- `Capture` / `Record` for fire-and-forget event capture
- Auto-attaches `requestId`, `userId` from [`appcontext`](../appcontext) and auth info from [`auth`](../auth)
- Single `Interface` for easy mocking

## Installation

```bash
go get github.com/downsized-devs/sdk-go/audit
```

## Quick Start

```go
a := audit.Init(audit.Config{})
a.Capture(ctx, audit.Collection{
    Resource: "user",
    Action:   "update",
    Before:   oldUser,
    After:    newUser,
})
```

## API Reference

| Symbol | Signature |
|---|---|
| `Init` | `func Init(cfg Config) Interface` |
| `Interface.Capture` | `(ctx, Collection)` |
| `Interface.Record` | `(ctx, Collection)` |
| `Collection` | Resource/action/before/after payload (see `audit/entity.go`). |

## Error Handling

`Capture`/`Record` do not return errors — failures are logged but never block business logic.

## Dependencies

- **Internal:** [`appcontext`](../appcontext), [`auth`](../auth), [`operator`](../operator)
- **External:** `github.com/rs/zerolog`

## Testing

```bash
go test ./audit/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md).

## Related Packages

- [`logger`](../logger) — operational logs vs durable audit events.
- [`appcontext`](../appcontext) — provides request/user metadata.
- [`auth`](../auth) — provides current-user info.
