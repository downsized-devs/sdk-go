# `slack` — Slack message sender

`import "github.com/downsized-devs/sdk-go/slack"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Tiny wrapper around `slack-go/slack` for sending messages with attachments.

## Features

- `SendMessage(channel, text, attachments...)`

## Installation

```bash
go get github.com/downsized-devs/sdk-go/slack
```

## Quick Start

```go
s := slack.Init(slack.Config{
    Token: "xoxb-...",
})
_ = s.SendMessage("#alerts", "Build failed", slack.Attachment{
    Title:  "CI",
    Text:   "main is red",
    Fields: []slack.AttachmentField{{Title: "Branch", Value: "main", Short: true}},
})
```

## API Reference

| Symbol | Signature |
|---|---|
| `Init` | `func Init(cfg Config) Interface` |
| `Interface.SendMessage` | `(channel, text string, attachments ...Attachment) error` |
| `Attachment` | `{ Title, Text, Color string; Fields []AttachmentField }` |
| `AttachmentField` | `{ Title, Value string; Short bool }` |

## Configuration

| Field | Description |
|---|---|
| `Token` | Slack bot token (xoxb-…). Store in a secrets manager. |

## Dependencies

- **External:** `github.com/slack-go/slack`

## Testing

```bash
go test ./slack/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md).

## Related Packages

- [`tracker`](../tracker) — for webhook-based alerting alternatives.
