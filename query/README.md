# `query` — SQL query/clause builder

`import "github.com/downsized-devs/sdk-go/query"`

**Stability:** Beta — bulk-insert helper carries a "find a better way" marker. See [STABILITY.md](../STABILITY.md).

Dynamic SQL clause builder driven by struct tags. Builds WHERE, ORDER BY, LIMIT, and OFFSET from typed parameter structs. Includes a bulk-insert helper (`FormatQueryForRows`).

## Features

- Struct-tag-driven WHERE clause builder (matches db/param/field tags)
- Cursor-based pagination (`Cursor` interface)
- Sort param normalisation (`sort_by`, `sort-by`, `sortBy`, `sortby` all accepted)
- `FormatQueryForRows` — turn `INSERT ... VALUES` + `[][]any` into a single parameterised query
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

```go
q := "INSERT INTO events(name, payload) VALUES"
full, args, err := query.FormatQueryForRows(ctx, q, [][]any{
    {"login", "{}"},
    {"logout", "{}"},
})
// full == "INSERT INTO events(name, payload) VALUES (?, ?), (?, ?)"
// args == ["login", "{}", "logout", "{}"]
```

## API Reference

| Symbol | Signature |
|---|---|
| `FormatQueryForRows` | `(ctx, q string, inputs [][]any) (string, []any, error)` |
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

Eleven test files cover the clause-builder side. Bulk-insert helper has minimal tests.

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). The "I hate this" comment in `query/query.go:12` is a real signal — proposals for a better bulk-insert API are welcome.

## Related Packages

- [`sql`](../sql) — provides the connection/transaction the builder ultimately runs against.
- [`null`](../null) — nullable types in parameter structs.
