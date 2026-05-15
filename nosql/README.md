# `nosql` — MongoDB CRUD wrapper

`import "github.com/downsized-devs/sdk-go/nosql"`

**Stability:** Beta — no tests in package. See [STABILITY.md](../STABILITY.md).

Thin wrapper around the official MongoDB Go driver. Provides typed CRUD without exposing the driver's full surface.

## Features

- `Find`, `FindOne`, `InsertOne`, `UpdateOne`, `UpdateMany`
- `Close` for clean shutdown

## Installation

```bash
go get github.com/downsized-devs/sdk-go/nosql
```

## Quick Start

```go
db := nosql.Init(nosql.Config{
    URI:      "mongodb://localhost:27017",
    Database: "app",
}, log)
defer db.Close(context.Background())

_ = db.InsertOne(ctx, "users", bson.M{"_id": "alice", "name": "Alice"})
var u User
_ = db.FindOne(ctx, "users", bson.M{"_id": "alice"}, &u)
```

## API Reference

| Symbol | Signature |
|---|---|
| `Init` | `func Init(cfg Config, log logger.Interface) Interface` |
| `Interface.Close` | `(ctx) error` |
| `Interface.Find` | `(ctx, collection string, filter, dest any) error` |
| `Interface.FindOne` | `(ctx, collection string, filter, dest any) error` |
| `Interface.InsertOne` | `(ctx, collection string, doc any) error` |
| `Interface.UpdateOne` | `(ctx, collection string, filter, update any) error` |
| `Interface.UpdateMany` | `(ctx, collection string, filter, update any) error` |

## Configuration

| Field | Description |
|---|---|
| `URI` | Mongo connection string. |
| `Database` | Default database name. |

## Error Handling

Wraps with [`codes`](../codes) nosql-range codes (1400–1499).

## Dependencies

- **Internal:** [`codes`](../codes), [`errors`](../errors), [`logger`](../logger)
- **External:** `go.mongodb.org/mongo-driver/mongo`, `.../mongo/options`

## Testing

No tests yet. Use `mongomock` or run against a local Mongo container.

```bash
go test ./nosql/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). Test coverage gates promotion to Stable.

## Related Packages

- [`sql`](../sql) — relational alternative.
- [`redis`](../redis) — cache layer in front of either.
