# `<pkgname>` — <one-line description>

> Template README for packages inside `github.com/downsized-devs/sdk-go`. Copy this file to `<pkgname>/README.md` and replace every `<…>` placeholder. Delete sections that genuinely don't apply rather than leaving "N/A".

`import "github.com/downsized-devs/sdk-go/<pkgname>"`

**Stability:** `Experimental | Beta | Stable | Deprecated` — see [../STABILITY.md](../STABILITY.md)

---

## Features

- <key capability 1>
- <key capability 2>
- <key capability 3>

## Installation

The package is part of the `sdk-go` monorepo. Either install the whole module:

```bash
go get github.com/downsized-devs/sdk-go
```

or import this package directly — Go's module resolution will pull only what's needed:

```go
import "github.com/downsized-devs/sdk-go/<pkgname>"
```

## Quick Start

```go
package main

import (
    "context"

    "github.com/downsized-devs/sdk-go/<pkgname>"
)

func main() {
    ctx := context.Background()
    client := <pkgname>.Init(<pkgname>.Config{
        // minimal config fields
    })
    _ = client
    _ = ctx
}
```

## API Reference

### Construction

| Symbol | Signature | Notes |
|---|---|---|
| `Init` | `func Init(cfg Config, ...) Interface` | Primary entry point. |

### Types

| Type | Purpose |
|---|---|
| `Interface` | Public surface of the package. Mock with `go.uber.org/mock`. |
| `Config` | Construction-time configuration. |

### Methods on `Interface`

| Method | Signature | Notes |
|---|---|---|
| `<MethodA>` | `(ctx, …) (…, error)` | <one line> |
| `<MethodB>` | `(ctx, …) (…, error)` | <one line> |

## Configuration

| Field | Type | Required | Default | Description |
|---|---|---|---|---|
| `<FieldA>` | `string` | yes | — | <description> |
| `<FieldB>` | `int` | no | `0` | <description> |

Environment variables (if any): list them here, e.g. `APP_<PKG>_FOO`.

## Examples

### Example 1 — minimal usage

```go
// ...
```

### Example 2 — typical production wiring

```go
// ...
```

## Error Handling

This package returns errors created via [`errors`](../errors) tagged with codes from [`codes`](../codes). Use `errors.GetCode(err)` to dispatch:

```go
if err != nil {
    switch errors.GetCode(err) {
    case codes.<CodeA>:
        // recoverable
    default:
        return err
    }
}
```

Common codes returned by this package: `<list them here>`.

## Dependencies

**Internal:** `<list of sibling sdk-go packages>` — keep this list in sync with imports.

**External:** `<github.com/foo/bar v1.2.3>` — only list direct imports, not transitive ones.

## Testing

Run the package tests:

```bash
go test ./<pkgname>/...
```

If you need to mock this package's `Interface` in downstream code:

```bash
make mock-all
```

(See the root `Makefile` for the mockgen invocation.)

## Contributing

Bug fixes and enhancements welcome. Before opening a PR:

1. Read [../CONTRIBUTING.md](../CONTRIBUTING.md).
2. Run `go test ./<pkgname>/...` and `golangci-lint run ./<pkgname>/...`.
3. Update this README if you add or change the public API.

## Related Packages

- [`<sibling-1>`](../<sibling-1>) — <why related>
- [`<sibling-2>`](../<sibling-2>) — <why related>

See the [package registry](../docs/PACKAGE_REGISTRY.md) for the full catalogue.
