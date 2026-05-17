# `sql` — SQL database abstraction with leader/follower routing

`import "github.com/downsized-devs/sdk-go/sql"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Multi-driver (`MySQL`, `Postgres`, `SQLite`) database client built on `sqlx`, with leader/follower split, prepared statements, transactions, and built-in Prometheus instrumentation.

## Features

- Leader/follower routing (`Leader()`, `Follower()`)
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
metrics := instrument.Init(instrument.Config{
    Metrics: instrument.MetricsConfig{Enabled: true},
})

db := sql.Init(sql.Config{
    Driver: "postgres",
    Leader: sql.ConnConfig{
        Host:     "leader",
        Port:     5432,
        DB:       "app",
        User:     "user",
        Password: "pass",
    },
    Follower: sql.ConnConfig{
        Host:     "follower",
        Port:     5432,
        DB:       "app",
        User:     "user",
        Password: "pass",
    },
}, log, metrics)
defer db.Stop()

var u User
if err := db.Follower().Get(ctx, &u, "SELECT * FROM users WHERE id=$1", 1); err != nil {
    if errors.Is(err, sql.ErrNotFound) { /* miss */ }
}
```

## API Reference

| Symbol | Signature |
|---|---|
| `Init` | `func Init(cfg Config, log logger.Interface, instr instrument.Interface) Interface` |
| `Interface.Leader` | `() Command` — read/write connection. |
| `Interface.Follower` | `() Command` — read-only connection. |
| `Interface.Stop` | `()` — close all pools. |
| `Command.Query` | `(ctx, dest, q, args...) error` |
| `Command.Exec` | `(ctx, q, args...) (Result, error)` |
| `Command.Get` | `(ctx, dest, q, args...) error` — single row. |
| `Command.BeginTx` | `(ctx, TxOptions) (Tx, error)` |
| `Command.Prepare` | `(ctx, q) (Stmt, error)` |

`ErrNotFound` is returned by `Get` when the row is missing.

## Configuration

| Field | Required | Description |
|---|---|---|
| `Driver` | yes | `mysql`, `postgres`, or `sqlite3`. |
| `Leader` / `Follower` | yes | `ConnConfig{Host, Port, DB, User, Password, SSL, Schema, Options}`. |
| `Leader.Options.MaxOpen`, `MaxIdle`, `MaxLifeTime` | no | Pool tuning. Same for follower. |

## Examples

### Run a transaction

```go
tx, err := db.Leader().BeginTx(ctx, sql.TxOptions{})
if err != nil { return err }
defer tx.Rollback()

if _, err := tx.Exec(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2", 100, 1); err != nil {
    return err
}
return tx.Commit()
```

### Use a prepared statement

```go
stmt, _ := db.Leader().Prepare(ctx, "INSERT INTO events(name, payload) VALUES (?, ?)")
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
