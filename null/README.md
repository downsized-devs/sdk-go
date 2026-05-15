# `null` — SQL-nullable, JSON-friendly types

`import "github.com/downsized-devs/sdk-go/null"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Nullable wrappers for `bool`, `int64`, `float64`, `string`, `time.Time` that round-trip cleanly between `database/sql` and JSON. Mirrors the shape of `sql.NullX` but adds an explicit `SqlNull` flag for "set to NULL in UPDATE" intent.

## Features

- Types: `Bool`, `Int64`, `Float64`, `String`, `Time`
- `New<T>` and `<T>From(v)` constructors
- `database/sql` `Scan`/`Value` implementations
- `encoding/json` `MarshalJSON`/`UnmarshalJSON` implementations
- `SqlNull` field on `Int64`, `String`, `Time` lets you distinguish "leave unchanged" from "set to NULL" in updates

## Installation

```bash
go get github.com/downsized-devs/sdk-go/null
```

## Quick Start

```go
import "github.com/downsized-devs/sdk-go/null"

u := User{
    Name:  null.StringFrom("Alice"),   // valid
    Email: null.NewString("", false),  // null
}

// JSON: {"Name":"Alice","Email":null}
// SQL:  Name='Alice', Email=NULL
```

## API Reference

| Type | Constructors | Special |
|---|---|---|
| `Bool` | `NewBool(b, valid)`, `BoolFrom(b)` | |
| `Int64` | `NewInt64(i, valid)`, `Int64From(i)` | `SqlNull bool` |
| `Float64` | `NewFloat64(f, valid)`, `Float64From(f)` | |
| `String` | `NewString(s, valid)`, `StringFrom(s)` | `SqlNull bool` |
| `Time` | `NewTime(t, valid)`, `TimeFrom(t)` | `SqlNull bool` |

All types implement `driver.Valuer`, `sql.Scanner`, `json.Marshaler`, `json.Unmarshaler`.

## Examples

### SQL UPDATE that nulls a column explicitly

```go
u.Email = null.String{SqlNull: true}
_, _ = db.Exec(ctx, "UPDATE users SET email = ? WHERE id = ?", u.Email, u.ID)
```

## Dependencies

stdlib only.

## Testing

```bash
go test ./null/...
```

Five test files (one per type).

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md).

## Related Packages

- [`sql`](../sql), [`query`](../query) — primary consumers.
- [`parser`](../parser) — JSON unmarshalling round-trips through these types.
