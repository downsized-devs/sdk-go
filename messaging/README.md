# `messaging` — Firebase Cloud Messaging client

`import "github.com/downsized-devs/sdk-go/messaging"`

**Stability:** Beta — no tests in package. See [STABILITY.md](../STABILITY.md).

Push notifications via Firebase Cloud Messaging (FCM). Manages device-token topic subscriptions and broadcasts.

## Features

- `SubscribeToTopic`, `UnsubscribeFromTopic`
- `BroadcastToTopic` — send a notification to every device on a topic.
- `BatchSendDryRun` — validate a batch send without delivering.

## Installation

```bash
go get github.com/downsized-devs/sdk-go/messaging
```

## Quick Start

```go
m := messaging.Init(messaging.Config{
    Firebase: messaging.FirebaseConf{ AccountKey: messaging.FirebaseAccountKey{ /* ... */ } },
}, log, jsonParser)

_ = m.SubscribeToTopic(ctx, []string{"deviceToken1"}, "news")
_ = m.BroadcastToTopic(ctx, "news", "Daily update", "...body...")
```

## API Reference

| Symbol | Signature |
|---|---|
| `Init` | `func Init(cfg Config, log logger.Interface, json parser.JsonInterface) Interface` |
| `Interface.SubscribeToTopic` | `(ctx, tokens []string, topic string) error` |
| `Interface.UnsubscribeFromTopic` | `(ctx, tokens []string, topic string) error` |
| `Interface.BroadcastToTopic` | `(ctx, topic, title, body string) error` |
| `Interface.BatchSendDryRun` | `(ctx, [...]) error` |

## Configuration

`Firebase.AccountKey` — full Firebase service-account JSON, loaded from secrets.

## Dependencies

- **Internal:** [`logger`](../logger), [`parser`](../parser)
- **External:** `firebase.google.com/go`, `firebase.google.com/go/messaging`, `google.golang.org/api/option`

## Testing

No tests yet. Use the Firebase Admin SDK's fake-app pattern when adding them.

```bash
go test ./messaging/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). Test coverage gates promotion to Stable.

## Related Packages

- [`auth`](../auth) — also uses Firebase; share service-account loading if practical.
- [`parser`](../parser) — for marshalling FCM payloads.
