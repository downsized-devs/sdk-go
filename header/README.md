# `header` — HTTP header & MIME constants

`import "github.com/downsized-devs/sdk-go/header"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Plain string constants for header keys, content types, and cache-control values. No tests; nothing to test.

## Features

- ~18 string constants (e.g. `RequestID`, `AcceptLanguage`, `ContentType`, `CacheControl`, `ApplicationJSON`).

## Installation

```bash
go get github.com/downsized-devs/sdk-go/header
```

## Quick Start

```go
import "github.com/downsized-devs/sdk-go/header"

w.Header().Set(header.ContentType, header.ApplicationJSON)
lang := r.Header.Get(header.AcceptLanguage)
```

## API Reference

See `header/header.go` for the full constant list. Constants only — no types or functions.

## Dependencies

stdlib only.

## Testing

No tests; constants are checked at compile time.

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). Add new constants alphabetically and update [`appcontext`](../appcontext) if the new header has a typed getter/setter.

## Related Packages

- [`appcontext`](../appcontext) — uses these constants when reading inbound headers.
