# `gqlclient` — low-level GraphQL HTTP client

`import "github.com/downsized-devs/sdk-go/gqlclient"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Minimal GraphQL client with JSON and multipart-form transport. Forked from [`machinebox/graphql`](https://github.com/machinebox/graphql); no internal sdk-go dependencies.

## Features

- Simple, familiar API
- Respects `context.Context` timeouts and cancellation
- Build and execute any kind of GraphQL request
- Use strong Go types for response data
- Use variables and upload files (multipart)
- Pluggable HTTP client via `WithHTTPClient`

## Installation

```bash
go get github.com/downsized-devs/sdk-go/gqlclient
```

## Quick Start

```go
client := gqlclient.NewClient("https://api.example.com/graphql")

req := gqlclient.NewRequest(`
    query($key: String!) {
        items(id: $key) { field1 field2 field3 }
    }
`)
req.Var("key", "value")
req.Header.Set("Cache-Control", "no-cache")

var resp struct{ Items []struct{ Field1, Field2, Field3 string } }
if err := client.Run(ctx, req, &resp); err != nil {
    return err
}
```

## API Reference

| Symbol | Signature |
|---|---|
| `NewClient` | `(endpoint string, opts ...ClientOption) *Client` |
| `NewRequest` | `(query string) *Request` |
| `Client.Run` | `(ctx, *Request, dest any) error` |
| `Request.Var` | `(key string, value any)` |
| `Request.Header` | `http.Header` field |
| `WithHTTPClient(*http.Client)` | option |
| `UseMultipartForm()` | option |
| `ImmediatelyCloseReqBody()` | option |
| `File` | upload struct (`Field`, `Name`, `R io.Reader`). |

## Examples

### File upload via multipart

```go
client := gqlclient.NewClient(url, gqlclient.UseMultipartForm())
req := gqlclient.NewRequest(`mutation($file: Upload!){ upload(file:$file) }`)
req.File("file", "photo.jpg", file)
_ = client.Run(ctx, req, nil)
```

## Dependencies

stdlib only.

## Testing

```bash
go test ./gqlclient/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). Upstream is the machinebox fork — preserve the original behaviour unless intentionally diverging.

## Related Packages

- [`parser`](../parser) — for handling GraphQL responses outside the client decoder.
