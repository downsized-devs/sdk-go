# `query` — SQL query/clause builder

`import "github.com/downsized-devs/sdk-go/query"`

**Stability:** Stable. See [STABILITY.md](../STABILITY.md).

Dynamic SQL clause builder driven by struct tags. Builds WHERE, ORDER BY, LIMIT, and OFFSET from typed parameter structs.

## Features

- Struct-tag-driven WHERE clause builder (matches db/param/field tags)
- Cursor-based pagination (`Cursor` interface)
- Sort param normalisation (`sort_by`, `sort-by`, `sortBy`, `sortby` all accepted)
- Typed converters for int/int8/.../uint64, float, string, bool, time, plus their `*Arr` variants

## Installation

```bash
go get github.com/downsized-devs/sdk-go/query
```

## Quick Start

```go
type ProductParam struct {
    Name     string `db:"name"     param:"name"     field:"Name"`
    IsActive bool   `db:"is_active" param:"is_active" field:"IsActive"`
    SortBy   string `param:"sort_by"`
    Page     int64  `param:"page"`
    Limit    int64  `param:"limit"`
}

// (Use the SQL builder via the package's exported helpers — see source for the
// current API shape; the builder is reflected into sql.Interface queries.)
```

### Bulk insert

For bulk inserts, use the `sql` package's `NamedExec` with a slice of structs (sqlx expands the named placeholders automatically):

```go
type Event struct {
    Name    string `db:"name"`
    Payload string `db:"payload"`
}
events := []Event{{"login", "{}"}, {"logout", "{}"}}
_, err := db.Leader().NamedExec(ctx, "insert-events",
    "INSERT INTO events (name, payload) VALUES (:name, :payload)", events)
```

## API Reference

| Symbol | Signature |
|---|---|
| `Cursor` | interface { `DecodeCursor(string) error`; `EncodeCursor() (string, error)` } |
| `Option` | `{ DisableLimit, IsActive, IsInactive bool }` |
| `Int`, `Int64`, `String`, etc. | primitive-type constants used by the clause builder. |

## Error Handling

- Empty rows or columns → coded error from [`codes`](../codes) (`CodeSQLPrepareStmt`).

## Dependencies

- **Internal:** [`codes`](../codes), [`errors`](../errors), [`null`](../null), [`sql`](../sql)
- **External:** `github.com/jmoiron/sqlx`

## Testing

```bash
go test ./query/...
```

Test files cover the clause builder, converters, sort normalisation, and cursor encoding.

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md).

## Related Packages

- [`sql`](../sql) — provides the connection/transaction the builder ultimately runs against.
- [`null`](../null) — nullable types in parameter structs.
