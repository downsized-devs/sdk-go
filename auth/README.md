# `auth` — Firebase authentication client

`import "github.com/downsized-devs/sdk-go/auth"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Firebase Auth wrapper providing token verification, refresh-token exchange, user CRUD, and password sign-in. The HTTP exchange for refresh tokens goes through Google's Identity Toolkit endpoint.

## Features

- Token verification via Firebase Admin SDK
- Password sign-in & verification
- Refresh-token exchange and revocation
- User CRUD: register, get (single + batch), update, delete
- Context helpers: `SetUserAuthInfo`, `GetUserAuthInfo`

## Installation

```bash
go get github.com/downsized-devs/sdk-go/auth
```

You'll need a Firebase service-account key file (JSON) and a Firebase Web API key.

## Quick Start

```go
import (
    "context"
    "net/http"

    "github.com/downsized-devs/sdk-go/auth"
    "github.com/downsized-devs/sdk-go/logger"
    "github.com/downsized-devs/sdk-go/parser"
)

log := logger.Init(logger.Config{Level: "info"})
json := parser.InitParser(log, parser.Options{}).JSONParser()

a := auth.Init(auth.Config{
    Firebase: auth.FirebaseConf{
        ApiKey:     "<FIREBASE_WEB_API_KEY>",
        AccountKey: auth.FirebaseAccountKey{ /* ... */ },
    },
}, log, json, http.DefaultClient)

tok, err := a.VerifyToken(context.Background(), "Bearer eyJhbGciOi...")
```

## API Reference

### Construction

```go
func Init(cfg Config, log logger.Interface, json parser.JsonInterface, httpClient *http.Client) Interface
```

### `Interface`

| Method | Signature |
|---|---|
| `VerifyToken` | `(ctx, bearer string) (*firebase_auth.Token, error)` |
| `SignInWithPassword` | `(ctx, UserLogin) (UserLoginResponse, error)` |
| `VerifyPassword` | `(ctx, email, password string) (bool, error)` |
| `RefreshToken` | `(ctx, refreshToken string) (RefreshTokenResponse, error)` |
| `RevokeUserRefreshToken` | `(ctx, uid string) error` |
| `GetUser` | `(ctx, FirebaseUserParam) ([]FirebaseUser, error)` |
| `GetUsers` | `(ctx, []FirebaseUserParam) ([]FirebaseUser, error)` |
| `RegisterUser` | `(ctx, FirebaseUser) (FirebaseUser, error)` |
| `UpdateUser` | `(ctx, FirebaseUser) (FirebaseUser, error)` |
| `DeleteUser` | `(ctx, uid string) error` |
| `SetUserAuthInfo` | `(ctx, UserAuthParam) context.Context` |
| `GetUserAuthInfo` | `(ctx) (UserAuthInfo, error)` |

Key types live in `auth/entity.go`: `FirebaseUser`, `UserLogin`, `UserLoginResponse`, `RefreshTokenRequest/Response`, `UserAuthInfo`, `UserAuthParam`.

## Configuration

| Field | Required | Description |
|---|---|---|
| `Firebase.ApiKey` | yes | Firebase Web API key (used at the Identity Toolkit refresh endpoint). |
| `Firebase.AccountKey` | yes | Service-account credentials. |

Load credentials from a secrets manager — never hard-code.

## Examples

### Gin middleware that verifies `Authorization`

```go
func AuthMiddleware(a auth.Interface) gin.HandlerFunc {
    return func(c *gin.Context) {
        tok, err := a.VerifyToken(c.Request.Context(), c.GetHeader("Authorization"))
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
            return
        }
        ctx := a.SetUserAuthInfo(c.Request.Context(), auth.UserAuthParam{UID: tok.UID})
        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}
```

### Rotate a leaked refresh token

```go
if err := a.RevokeUserRefreshToken(ctx, uid); err != nil {
    return errors.WrapWithCode(err, codes.CodeFirebaseRevokeToken, "rotate session")
}
```

## Error Handling

Errors are wrapped with [`codes`](../codes) auth-range codes (1700–1799). Dispatch with `errors.GetCode(err)`.

## Dependencies

- **Internal:** [`codes`](../codes), [`errors`](../errors), [`logger`](../logger), [`null`](../null), [`parser`](../parser)
- **External:** `firebase.google.com/go`, `firebase.google.com/go/auth`, `google.golang.org/api/identitytoolkit/v1`, `google.golang.org/api/identitytoolkit/v3`, `google.golang.org/api/option`

## Testing

```bash
go test ./auth/...
```

Uses fixture-based mocks — no live Firebase needed.

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). Changing `Interface` is breaking — coordinate with downstream consumers first.

## Related Packages

- [`security`](../security) — for password hashing if you keep your own user store.
- [`ratelimiter`](../ratelimiter) — pair with `auth` to rate-limit token verification.
- [`appcontext`](../appcontext) — for request/user metadata that flows alongside auth state.
