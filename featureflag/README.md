# `featureflag` — wrapper around `go-feature-flag`

`import "github.com/downsized-devs/sdk-go/featureflag"`

**Stability:** Beta — no tests in package. See [STABILITY.md](../STABILITY.md).

Wraps [`go-feature-flag`](https://github.com/thomaspoignant/go-feature-flag) for user-targeted boolean and JSON flag evaluation.

## Features

- `CheckUserFlags(ctx, user, defaults)` — evaluate flags for a user.
- `GetAllUserFlags(ctx, user)` — full flag snapshot.
- `Refresh()` — pull updated rules from the configured retriever.

## Installation

```bash
go get github.com/downsized-devs/sdk-go/featureflag
```

## Quick Start

```go
ff := featureflag.Init(featureflag.Config{
    Enabled: true,
    PollingInterval: 60,
    Retriever: featureflag.Retriever{ /* file/HTTP/S3 source */ },
}, log)
defer ff.Close()

flags := ff.CheckUserFlags(ctx, ffuser.NewUser("alice"), map[string]any{"newCheckout": false})
```

## API Reference

| Symbol | Signature |
|---|---|
| `Init` | `func Init(cfg Config, log logger.Interface) Interface` |
| `Interface.CheckUserFlags` | `(ctx, user, defaults) map[string]any` |
| `Interface.GetAllUserFlags` | `(ctx, user) map[string]any` |
| `Interface.Refresh` | `() error` |

## Configuration

See `Config` in `featureflag.go` — fields cover enable toggle, polling interval, retriever source (file / HTTP / S3).

## Dependencies

- **Internal:** [`logger`](../logger)
- **External:** `github.com/thomaspoignant/go-feature-flag`, `.../ffuser`, `.../retriever`

## Testing

No tests yet. When adding them, mock the `go-feature-flag` client with `gomock`.

```bash
go test ./featureflag/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). Test coverage gates promotion to Stable.

## Related Packages

- [`logger`](../logger) — required at `Init` time.
