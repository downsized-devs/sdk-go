# `sql` — SQL database abstraction with leader/follower routing

`import "github.com/downsized-devs/sdk-go/sql"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Multi-driver (`MySQL`, `Postgres`, `SQLite`) database client built on `sqlx`, with leader/follower split, prepared statements, transactions, and built-in Prometheus instrumentation.

## Features

- Leader/follower routing (`Leader(ctx)`, `Follower(ctx)`)
- Multi-driver: `github.com/go-sql-driver/mysql`, `github.com/lib/pq`, `modernc.org/sqlite`
- Transactions with `BeginTx`
- Prepared statements (`Prepare`)
- Automatic [`instrument`](../instrument) metrics for query timings + connection pool stats
- `ErrNotFound` sentinel for "no rows" lookups

## Installation

```bash
go get github.com/downsized-devs/sdk-go/sql
```

## Quick Start

```go
import (
    "context"

    "github.com/downsized-devs/sdk-go/instrument"
    "github.com/downsized-devs/sdk-go/logger"
    "github.com/downsized-devs/sdk-go/sql"
)

log := logger.Init(logger.Config{Level: "info"})
metrics := instrument.Init(instrument.Config{Enabled: true})

db := sql.Init(sql.Config{
    Driver: "postgres",
    Leader: sql.ConnConfig{
        DSN: "postgres://user:pass@leader:5432/app?sslmode=disable",
    },
    Follower: sql.ConnConfig{
        DSN: "postgres://user:pass@follower:5432/app?sslmode=disable",
    },
}, log, metrics)
defer db.Stop()

var u User
if err := db.Follower(ctx).Get(ctx, &u, "SELECT * FROM users WHERE id=$1", 1); err != nil {
    if errors.Is(err, sql.ErrNotFound) { /* miss */ }
}
```

## API Reference

| Symbol | Signature |
|---|---|
| `Init` | `func Init(cfg Config, log logger.Interface, metrics instrument.Interface) Interface` |
| `Interface.Leader` | `(ctx) Command` — read/write connection. |
| `Interface.Follower` | `(ctx) Command` — read-only connection. |
| `Interface.Stop` | `() error` — close all pools. |
| `Command.Query` | `(ctx, dest, q, args...) error` |
| `Command.Exec` | `(ctx, q, args...) (Result, error)` |
| `Command.Get` | `(ctx, dest, q, args...) error` — single row. |
| `Command.BeginTx` | `(ctx, TxOptions) (Tx, error)` |
| `Command.Prepare` | `(ctx, q) (Stmt, error)` |

`ErrNotFound` is returned by `Get` when the row is missing.

## Configuration

| Field | Required | Description |
|---|---|---|
| `Driver` | yes | `mysql`, `postgres`, or `sqlite`. |
| `Leader.DSN` / `Follower.DSN` | yes | Connection strings. |
| `Leader.MaxOpen`, `MaxIdle`, `MaxLifetime` | no | Pool tuning. Same for follower. |

## Examples

### Run a transaction

```go
tx, err := db.Leader(ctx).BeginTx(ctx, sql.TxOptions{})
if err != nil { return err }
defer tx.Rollback()

if _, err := tx.Exec(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2", 100, 1); err != nil {
    return err
}
return tx.Commit()
```

### Use a prepared statement

```go
stmt, _ := db.Leader(ctx).Prepare(ctx, "INSERT INTO events(name, payload) VALUES (?, ?)")
defer stmt.Close()
for _, e := range events {
    _, _ = stmt.Exec(ctx, e.Name, e.Payload)
}
```

## Error Handling

| Error | Cause | Action |
|---|---|---|
| `sql.ErrNotFound` | `Get` returned no rows. | Treat as miss. |
| Coded errors | Connection, syntax, constraint. | Inspect with `errors.GetCode(err)`. |

## Dependencies

- **Internal:** [`codes`](../codes), [`errors`](../errors), [`instrument`](../instrument), [`logger`](../logger)
- **External:** `github.com/jmoiron/sqlx`, `github.com/go-sql-driver/mysql`, `github.com/lib/pq`, `modernc.org/sqlite`

## Testing

```bash
go test ./sql/...
```

Five test files exercise sqlite in-memory; MySQL/Postgres tests require live databases.

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). Any change to the `Interface` or `Command` interface is breaking — coordinate with downstream services.

## Related Packages

- [`query`](../query) — dynamic WHERE/ORDER builder that consumes `sql.Interface`.
- [`null`](../null) — nullable types matching SQL semantics.
- [`instrument`](../instrument) — receives DB pool stats and query timings.
- [`redis`](../redis) — pair for cache-aside reads.
