# `security` — cryptographic primitives

`import "github.com/downsized-devs/sdk-go/security"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

AES-GCM symmetric encryption, PBKDF2-based password hashing, and Scrypt-based password hashing.

## Features

- `Encrypt` / `Decrypt` — AES-GCM symmetric.
- `HashPassword` / `CompareHashPassword` — PBKDF2.
- `ScryptPassword` / `CompareScryptPassword` — scrypt.

## Installation

```bash
go get github.com/downsized-devs/sdk-go/security
```

## Quick Start

```go
sec := security.Init(security.Config{
    SecretKey: "<32-byte-key>",
}, log)

enc, _ := sec.Encrypt(ctx, []byte("plaintext"))
plain, _ := sec.Decrypt(ctx, enc)

hash, _ := sec.HashPassword("hunter2")
ok, _   := sec.CompareHashPassword("hunter2", hash)
```

## API Reference

| Symbol | Signature |
|---|---|
| `Init` | `func Init(cfg Config, log logger.Interface) Interface` |
| `Interface.Encrypt` | `(ctx, []byte) ([]byte, error)` |
| `Interface.Decrypt` | `(ctx, []byte) ([]byte, error)` |
| `Interface.HashPassword` | `(plain string) (string, error)` |
| `Interface.CompareHashPassword` | `(plain, hash string) (bool, error)` |
| `Interface.ScryptPassword` | `(plain string, cfg ScryptConfig) (string, error)` |
| `Interface.CompareScryptPassword` | `(plain, hash string, cfg ScryptConfig) (bool, error)` |
| `ScryptConfig` | `{ N, R, P, KeyLength int }` |

## Configuration

| Field | Description |
|---|---|
| `SecretKey` | Symmetric encryption key (32 bytes for AES-256-GCM). |
| Scrypt parameters | Passed per-call as `ScryptConfig`. |

## Error Handling

Errors are wrapped with [`codes`](../codes) crypto codes.

## Dependencies

- **Internal:** [`codes`](../codes), [`errors`](../errors), [`logger`](../logger)
- **External:** `golang.org/x/crypto/pbkdf2`, `golang.org/x/crypto/scrypt`

## Testing

```bash
go test ./security/...
```

Two test files. Crypto is high-risk — never lower the assertion bar.

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). Any change to encryption schemes must go through a second reviewer with a written rationale.

## Related Packages

- [`auth`](../auth) — for Firebase-managed identities (passwords stay with Firebase).
- [`character`](../character) — pre-hash strength validation.
