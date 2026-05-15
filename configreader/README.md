# `configreader` — layered configuration parser

`import "github.com/downsized-devs/sdk-go/configreader"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Reads JSON/YAML config files via viper, with JSON-reference resolution and custom `time.Duration` decoding. Pairs with [`configbuilder`](../configbuilder).

## Features

- `ReadConfig(target any)` — unmarshal into your struct
- `AllSettings()` — flat view of every key
- JSON reference resolution (`$ref`) for shared fragments
- Built-in `time.Duration` and `string-to-slice` decode hooks

## Installation

```bash
go get github.com/downsized-devs/sdk-go/configreader
```

## Quick Start

```go
type AppConfig struct {
    Redis struct {
        Address    string
        DefaultTTL time.Duration
    }
    Logger struct{ Level string }
}

r := configreader.Init(configreader.Options{
    ConfigFile: "config/runtime.json",
})
var cfg AppConfig
if err := r.ReadConfig(&cfg); err != nil {
    log.Fatal(err)
}
```

## API Reference

| Symbol | Signature |
|---|---|
| `Init` | `func Init(opts Options) Interface` |
| `Interface.ReadConfig` | `(target any) error` |
| `Interface.AllSettings` | `() map[string]any` |
| `Options.ConfigFile` | `string` — path to file. |
| `Options.Type` | `string` (json, yaml, toml…) — defaults inferred from extension. |
| `Options.AdditionalConfigOptions` | `[]AdditionalConfigOptions` — extra files merged in. |

## Examples

### Layer defaults + environment overrides

```go
r := configreader.Init(configreader.Options{
    ConfigFile: "config/defaults.json",
    AdditionalConfigOptions: []configreader.AdditionalConfigOptions{
        { ConfigFile: fmt.Sprintf("config/%s.json", os.Getenv("ENV")), Type: "json" },
    },
})
```

### Use with `Duration` fields

```go
type Cfg struct {
    Timeout time.Duration // "5s", "1m30s" both decoded
}
```

## Error Handling

`ReadConfig` returns descriptive errors for missing files, type mismatches, and unresolved `$ref` pointers. Don't continue on error — config-load failure should be fatal.

## Dependencies

- **Internal:** [`files`](../files)
- **External:** `github.com/mitchellh/mapstructure`, `github.com/spf13/viper`

## Testing

```bash
go test ./configreader/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md).

## Related Packages

- [`configbuilder`](../configbuilder) — write side; produces the files this package reads.
- [`files`](../files) — path-existence checks.
