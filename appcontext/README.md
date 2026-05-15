# `appcontext` — request-scoped context value helpers

`import "github.com/downsized-devs/sdk-go/appcontext"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Typed get/set helpers for the values [`logger`](../logger), [`audit`](../audit), and others expect to find on `context.Context` — request ID, user ID, accept-language, service version, device type, app response code.

## Features

- Setter/getter pairs with private context keys (no collisions with other packages)
- Read by [`logger`](../logger) for automatic field enrichment
- Zero external dependencies

## Installation

```bash
go get github.com/downsized-devs/sdk-go/appcontext
```

## Quick Start

```go
package main

import (
    "context"

    "github.com/downsized-devs/sdk-go/appcontext"
)

func main() {
    ctx := context.Background()
    ctx = appcontext.SetRequestId(ctx, "req-abc")
    ctx = appcontext.SetUserId(ctx, 42)

    _ = appcontext.GetRequestId(ctx) // "req-abc"
    _ = appcontext.GetUserId(ctx)    // 42
}
```

## API Reference

| Method | Purpose |
|---|---|
| `SetRequestId` / `GetRequestId` | Request correlation ID. |
| `SetUserId` / `GetUserId` | Authenticated user identifier. |
| `SetAcceptLanguage` / `GetAcceptLanguage` | Negotiated locale for response rendering. |
| `SetServiceVersion` / `GetServiceVersion` | Service version stamp (build tag, git sha). |
| `SetDeviceType` / `GetDeviceType` | Client device class. |
| `SetAppResponseCode` / `GetAppResponseCode` | Pinned response code for late middleware. |

All `Set*` return a new `context.Context`; the original is not mutated.

## Examples

### Gin middleware that populates everything from headers

```go
func ContextEnricher() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := c.Request.Context()
        ctx = appcontext.SetRequestId(ctx, c.GetHeader(header.RequestID))
        ctx = appcontext.SetAcceptLanguage(ctx, c.GetHeader(header.AcceptLanguage))
        ctx = appcontext.SetServiceVersion(ctx, buildSHA)
        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}
```

## Error Handling

Getters return the zero value when the key is unset; no error is returned. Callers that *require* a value should validate at their own boundary.

## Dependencies

- **Internal:** [`codes`](../codes), [`header`](../header), [`language`](../language)
- **External:** stdlib only

## Testing

```bash
go test ./appcontext/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). Adding a new key is additive — pair every `Set*` with a `Get*` and update [`logger`](../logger) if the field should appear in log lines automatically.

## Related Packages

- [`logger`](../logger) — primary consumer.
- [`header`](../header) — string constants for HTTP header names.
- [`translator`](../translator) — reads `AcceptLanguage`.
- [`audit`](../audit) — reads `RequestId` and `UserId`.
