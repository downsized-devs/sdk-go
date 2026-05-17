# `storage` — AWS S3 client

`import "github.com/downsized-devs/sdk-go/storage"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Wraps AWS SDK for Go (v1) for the common S3 operations: upload, download, delete, and presigned URLs.

## Features

- `Upload`, `Download`, `Delete`
- `GetPresignedUrl`, `GetPresignedUrlWithDuration`
- `CreateUrlByKey` for static URLs

## Installation

```bash
go get github.com/downsized-devs/sdk-go/storage
```

## Quick Start

```go
s := storage.Init(storage.Config{
    AWSS3: storage.AWSS3Config{
        Region:          "ap-southeast-1",
        BucketName:      "my-bucket",
        AccessKeyID:     "<KEY>",
        SecretAccessKey: "<SECRET>",
        PresignDuration: 15 * time.Minute,
    },
}, log)

url, _ := s.Upload(ctx, "key/path.jpg", "path.jpg", "image/jpeg", fileBytes)
presigned, _ := s.GetPresignedUrl(ctx, "key/path.jpg")
```

## API Reference

| Symbol | Signature |
|---|---|
| `Init` | `func Init(cfg Config, log logger.Interface) Interface` |
| `Interface.Upload` | `(ctx, key string, filename, contentType string, data []byte) (url string, err error)` |
| `Interface.Download` | `(ctx, key string) ([]byte, error)` |
| `Interface.Delete` | `(ctx, key string) error` |
| `Interface.GetPresignedUrl` | `(ctx, key string) (string, error)` |
| `Interface.GetPresignedUrlWithDuration` | `(ctx, key string, d time.Duration) (string, error)` |
| `Interface.CreateUrlByKey` | `(key string) string` |

## Configuration

| Field | Required | Description |
|---|---|---|
| `AWSS3.Region` | yes | S3 region. |
| `AWSS3.BucketName` | yes | Default bucket. |
| `AWSS3.AccessKeyID` / `AWSS3.SecretAccessKey` | yes | IAM credentials. Use IRSA / IAM-role-for-service-account in EKS instead when possible. |
| `AWSS3.PresignDuration` | no | Default expiry used by `GetPresignedUrl`. |

## Error Handling

Errors are wrapped with [`codes`](../codes) third-party codes.

## Dependencies

- **Internal:** [`codes`](../codes), [`errors`](../errors), [`logger`](../logger)
- **External:** `github.com/aws/aws-sdk-go/aws`, `.../credentials`, `.../session`, `.../service/s3`

## Testing

```bash
go test ./storage/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). Migrating to AWS SDK v2 would be a major change — file an issue first.

## Related Packages

- [`localstorage`](../localstorage) — on-disk full-text index (different problem).
- [`files`](../files) — local-file helpers.
