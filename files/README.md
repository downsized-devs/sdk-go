# `files` — filesystem helpers

`import "github.com/downsized-devs/sdk-go/files"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

## Features

- `GetExtension(filename string) string` — extension after the last dot, lowercased.
- `IsExist(filename string) bool` — file existence check.

## Installation

```bash
go get github.com/downsized-devs/sdk-go/files
```

## Quick Start

```go
import "github.com/downsized-devs/sdk-go/files"

_ = files.GetExtension("/tmp/report.PDF") // "pdf"
_ = files.IsExist("/tmp/report.pdf")       // true/false
```

## API Reference

| Symbol | Signature |
|---|---|
| `GetExtension` | `(filename string) string` |
| `IsExist` | `(filename string) bool` |

## Dependencies

stdlib only.

## Testing

```bash
go test ./files/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md).

## Related Packages

- [`configbuilder`](../configbuilder), [`configreader`](../configreader) — uses `files.IsExist` to validate paths.
- [`pdf`](../pdf) — uses `files.GetExtension` for input validation.
