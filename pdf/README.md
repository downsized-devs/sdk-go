# `pdf` — PDF password protection

`import "github.com/downsized-devs/sdk-go/pdf"`

**Stability:** Beta — single feature, minimal tests. See [STABILITY.md](../STABILITY.md).

Encrypts an existing PDF with AES-256 password protection using [`pdfcpu`](https://github.com/pdfcpu/pdfcpu).

## Features

- `SetPasswordFile(inPath, outPath, ownerPwd, userPwd)` — write a password-protected copy.

## Installation

```bash
go get github.com/downsized-devs/sdk-go/pdf
```

## Quick Start

```go
p := pdf.Init(log)
err := p.SetPasswordFile(ctx, "/tmp/in.pdf", "/tmp/out.pdf", "ownerPwd", "userPwd")
```

## API Reference

| Symbol | Signature |
|---|---|
| `Init` | `func Init(log logger.Interface) PdfInterface` |
| `PdfInterface.SetPasswordFile` | `(ctx, in, out, ownerPwd, userPwd string) error` |

## Dependencies

- **Internal:** [`logger`](../logger)
- **External:** `github.com/pdfcpu/pdfcpu/pkg/api`, `.../pdfcpu/model`

## Testing

```bash
go test ./pdf/...
```

One test file with basic coverage.

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). Add tests as part of any non-trivial change — Beta classification will hold until coverage improves.

## Related Packages

- [`files`](../files) — extension validation for input paths.
- [`security`](../security) — for non-PDF crypto needs.
