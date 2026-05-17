# `messaging` — Firebase Cloud Messaging client

`import "github.com/downsized-devs/sdk-go/messaging"`

**Stability:** Stable. See [STABILITY.md](../STABILITY.md).

Push notifications via Firebase Cloud Messaging (FCM). Manages device-token topic subscriptions and broadcasts.

## Features

- `SubscribeToTopic`, `UnsubscribeFromTopic`
- `BroadcastToTopic` — send a data payload to every device on a topic.
- `BatchSendDryRun` — validate token batches without delivering; chunked to `MaximumTokensPerBatch` (500) per call.

## Installation

```bash
go get github.com/downsized-devs/sdk-go/messaging
```

## Quick Start

```go
m := messaging.Init(messaging.Config{
    Firebase: messaging.FirebaseConf{
        AccountKey: messaging.FirebaseAccountKey{ /* ... */ },
    },
}, log, jsonParser, nil) // httpClient is deprecated; pass nil

_ = m.SubscribeToTopic(ctx, "deviceToken1", "news")
_ = m.BroadcastToTopic(ctx, "news", map[string]string{
    "title": "Daily update",
    "body":  "...body...",
})
```

## API Reference

| Symbol | Signature |
|---|---|
| `Init` | `func Init(cfg Config, log logger.Interface, json parser.JsonInterface, _ *http.Client) Interface` |
| `Interface.SubscribeToTopic` | `(ctx, deviceToken string, topic string) error` |
| `Interface.UnsubscribeFromTopic` | `(ctx, deviceToken string, topic string) error` |
| `Interface.BroadcastToTopic` | `(ctx, topic string, payload map[string]string) error` |
| `Interface.BatchSendDryRun` | `(ctx, tokens []string) ([]string, error)` |

The `httpClient` argument to `Init` is accepted for backwards-compatibility but ignored — see the GoDoc for the `Deprecated:` notice.

## Configuration

`Firebase.AccountKey` — full Firebase service-account JSON, loaded from secrets.

## Dependencies

- **Internal:** [`logger`](../logger), [`parser`](../parser)
- **External:** `firebase.google.com/go`, `firebase.google.com/go/messaging`, `google.golang.org/api/option`

## Testing

Unit tests cover the public surface using a `firebaseMessenger` seam; no network calls are made.

```bash
go test ./messaging/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md).

## Related Packages

- [`auth`](../auth) — also uses Firebase; share service-account loading if practical.
- [`parser`](../parser) — for marshalling FCM payloads.
