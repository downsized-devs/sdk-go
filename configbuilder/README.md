# `configbuilder` — mustache-template config generator

`import "github.com/downsized-devs/sdk-go/configbuilder"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Renders mustache templates into runtime config files at startup. Useful when you want a single source-of-truth template that fans out into per-environment JSON/YAML files.

## Features

- `BuildConfig(opts Options)` — render and write
- Mustache placeholders, optionally fed by viper or environment

## Installation

```bash
go get github.com/downsized-devs/sdk-go/configbuilder
```

## Quick Start

```go
import "github.com/downsized-devs/sdk-go/configbuilder"

b := configbuilder.Init()
b.BuildConfig(configbuilder.Options{
    TemplatePath: "config/app.json.tmpl",
    OutputPath:   "config/app.json",
    Data:         map[string]any{ "env": "prod", "host": os.Getenv("APP_HOST") },
})
```

## API Reference

| Symbol | Signature |
|---|---|
| `Init` | `func Init() Interface` |
| `Interface.BuildConfig` | `(Options) error` |
| `Options.TemplatePath` | `string` — path to `.tmpl` file. |
| `Options.OutputPath` | `string` — destination file. |
| `Options.Data` | `map[string]any` — values for `{{placeholders}}`. |

## Examples

### Generate per-environment config from CI

```bash
ENV=staging go run ./tools/build_config
```

```go
b.BuildConfig(configbuilder.Options{
    TemplatePath: fmt.Sprintf("config/%s.json.tmpl", os.Getenv("ENV")),
    OutputPath:   "config/runtime.json",
    Data:         map[string]any{
        "redisAddr": os.Getenv("REDIS_ADDR"),
        "dbDSN":     os.Getenv("DB_DSN"),
    },
})
```

## Error Handling

`BuildConfig` returns errors from missing templates and write failures. Surface them at process startup — failing to build config should be fatal.

## Dependencies

- **Internal:** [`files`](../files)
- **External:** `github.com/cbroglie/mustache`, `github.com/spf13/viper`

## Testing

```bash
go test ./configbuilder/...
```

Two test files (unit + integration).

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md).

## Related Packages

- [`configreader`](../configreader) — the read side: parses the file `configbuilder` produces.
- [`files`](../files) — used internally for path checks.
