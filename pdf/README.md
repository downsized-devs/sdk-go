# `pdf` — PDF manipulation

`import "github.com/downsized-devs/sdk-go/pdf"`

**Stability:** Stable. See [STABILITY.md](../STABILITY.md).

In-memory PDF operations: AES-256 password protection, page-level merge / split, text watermarking, page count, and best-effort text extraction.

## Features

| Method | What it does |
|---|---|
| `Encrypt`           | AES-256 password protect (sets both user and owner password). |
| `RemovePassword`    | Strip password from an encrypted PDF given the password. |
| `Merge`             | Concatenate two or more PDFs into one. |
| `Split`             | Split a PDF into chunks of `span` pages. |
| `AddTextWatermark`  | Stamp text across every page. |
| `ExtractText`       | Best-effort plain text extraction (no layout). |
| `PageCount`         | Number of pages in a PDF. |

All methods take and return `[]byte` so no temp files are visible to callers (Split uses a temp directory internally; it is removed before returning).

## Installation

```bash
go get github.com/downsized-devs/sdk-go/pdf
```

## Quick Start

```go
log := logger.Init(logger.Config{})
p := pdf.Init(log)

encrypted, err := p.Encrypt(ctx, raw, "s3cret")
if err != nil { /* ... */ }

merged, err := p.Merge(ctx, partA, partB, partC)
chunks, err := p.Split(ctx, big, 1)            // one page per chunk
stamped, err := p.AddTextWatermark(ctx, raw, "DRAFT")
text,    err := p.ExtractText(ctx, raw)
pages,   err := p.PageCount(ctx, raw)
```

## API Reference

| Symbol | Signature |
|---|---|
| `Init` | `func Init(log logger.Interface) Interface` |
| `Interface.Encrypt`          | `(ctx, data []byte, password string) ([]byte, error)` |
| `Interface.RemovePassword`   | `(ctx, data []byte, password string) ([]byte, error)` |
| `Interface.Merge`            | `(ctx, parts ...[]byte) ([]byte, error)` |
| `Interface.Split`            | `(ctx, data []byte, span int) ([][]byte, error)` |
| `Interface.AddTextWatermark` | `(ctx, data []byte, text string) ([]byte, error)` |
| `Interface.ExtractText`      | `(ctx, data []byte) (string, error)` |
| `Interface.PageCount`        | `(ctx, data []byte) (int, error)` |

`ctx` is accepted on every method for forward compatibility; cancellation is not yet propagated to the underlying libraries.

## Error Handling

All `[]byte` inputs are rejected with `codes.CodeBadRequest` when empty. Underlying pdfcpu / ledongthuc errors are surfaced unchanged.

## Dependencies

- **Internal:** [`codes`](../codes), [`errors`](../errors), [`logger`](../logger)
- **External:** `github.com/pdfcpu/pdfcpu` (all operations except `ExtractText`), `github.com/ledongthuc/pdf` (`ExtractText` only).

## Testing

```bash
go test ./pdf/...
```

Coverage ≥85% via the bundled `example.pdf` round-tripped through every method.

## Related Packages

- [`files`](../files) — extension validation for input paths.
- [`security`](../security) — for non-PDF crypto needs.
