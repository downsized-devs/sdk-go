# `localstorage` — Bleve-backed local full-text search

`import "github.com/downsized-devs/sdk-go/localstorage"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Local on-disk search index using [Bleve](https://github.com/blevesearch/bleve). For embedded search inside a single service — not a clustered solution.

## Features

- `NewIndex(name)` — create a new index on disk.
- `Index(name, id, doc)` — add or replace a document.
- `Search(name, query)` — text search returning IDs and scores.
- `DeleteIndex(name)` — drop the index.

## Installation

```bash
go get github.com/downsized-devs/sdk-go/localstorage
```

## Quick Start

```go
ls := localstorage.Init(localstorage.Config{
    Path: "./search-index",
}, log)

_ = ls.NewIndex(ctx, "products")
_ = ls.Index(ctx, "products", "p1", map[string]any{"name": "Apple"})
hits, _ := ls.Search(ctx, "products", "apple")
```

## API Reference

| Symbol | Signature |
|---|---|
| `Init` | `func Init(cfg Config, log logger.Interface) Interface` |
| `Interface.NewIndex` | `(ctx, name string) error` |
| `Interface.Index` | `(ctx, name, id string, doc any) error` |
| `Interface.Search` | `(ctx, name, query string) (Results, error)` |
| `Interface.DeleteIndex` | `(ctx, name string) error` |

## Configuration

| Field | Description |
|---|---|
| `Path` | Filesystem path where indexes live. |

## Dependencies

- **Internal:** [`logger`](../logger)
- **External:** `github.com/blevesearch/bleve`

## Testing

```bash
go test ./localstorage/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). The on-disk index format is owned by Bleve — don't make assumptions about layout.

## Related Packages

- [`storage`](../storage) — remote S3-backed object storage.
- [`logger`](../logger) — required at `Init` time.
