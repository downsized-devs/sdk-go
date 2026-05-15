# `email` — SMTP email sender with MJML templating

`import "github.com/downsized-devs/sdk-go/email"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Sends transactional email over SMTP, with mustache/MJML templating for the body.

## Features

- `SendEmail` — SMTP send.
- `GenerateBody` — render template + data into HTML.
- `FromHTML`, `FromMJML` — template source helpers.

## Installation

```bash
go get github.com/downsized-devs/sdk-go/email
```

## Quick Start

```go
mailer := email.Init(email.Config{
    Host: "smtp.example.com", Port: 587,
    Username: "user", Password: "pass",
    From: "no-reply@example.com",
}, log)

body, _ := mailer.GenerateBody().FromMJML("templates/welcome.mjml", map[string]any{"Name": "Alice"})
_ = mailer.SendEmail(ctx, email.Payload{
    To:      []string{"alice@example.com"},
    Subject: "Welcome",
    Body:    body,
})
```

## API Reference

| Symbol | Signature |
|---|---|
| `Init` | `func Init(cfg Config, log logger.Interface) Interface` |
| `Interface.SendEmail` | `(ctx, Payload) error` |
| `Interface.GenerateBody` | `() TemplateInterface` |
| `TemplateInterface.FromHTML` | `(path string, data any) (string, error)` |
| `TemplateInterface.FromMJML` | `(path string, data any) (string, error)` |

## Configuration

| Field | Required | Description |
|---|---|---|
| `Host`, `Port` | yes | SMTP endpoint. |
| `Username`, `Password` | yes | SMTP auth. |
| `From` | yes | Default sender. |

## Error Handling

Errors are wrapped with [`codes`](../codes) (email/third-party range).

## Dependencies

- **Internal:** [`codes`](../codes), [`errors`](../errors), [`logger`](../logger)
- **External:** `gopkg.in/gomail.v2`, `github.com/Boostport/mjml-go`

## Testing

```bash
go test ./email/...
```

Two test files cover template rendering and SMTP send (mocked).

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md).

## Related Packages

- [`parser`](../parser) — for parsing email payloads from external services.
- [`logger`](../logger) — required at `Init` time.
