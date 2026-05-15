# `parser` — JSON & CSV parsing with schema validation

`import "github.com/downsized-devs/sdk-go/parser"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

JSON parsing via `json-iterator/go` (configurable presets), CSV via `gocsv`, and optional JSON-Schema validation via `gojsonschema`.

## Features

- `JsonInterface` with `Marshal`/`Unmarshal` variants (5 modes: standard, default, fastest, custom, etc.)
- JSON Schema validation against a registered schema set
- `CsvInterface` for typed CSV reads
- Single `InitParser` factory returns both

## Installation

```bash
go get github.com/downsized-devs/sdk-go/parser
```

## Quick Start

```go
p := parser.InitParser(log, parser.Options{
    Json: parser.JsonOptions{ Config: parser.DefaultConfig },
})
jp := p.JSONParser()

raw, _ := jp.Marshal(map[string]any{"hello": "world"})
var out map[string]any
_ = jp.Unmarshal(raw, &out)
```

## API Reference

| Symbol | Signature |
|---|---|
| `InitParser` | `func InitParser(log logger.Interface, opts Options) Parser` |
| `Parser.JSONParser` | `() JsonInterface` |
| `Parser.CSVParser` | `() CsvInterface` |
| `JsonInterface.Marshal` | `(any) ([]byte, error)` |
| `JsonInterface.Unmarshal` | `([]byte, any) error` |
| `JsonInterface.ValidateAgainstSchema` | `(name string, data []byte) error` |
| `CsvInterface.Marshal` / `Unmarshal` | typed conversions |

## Configuration

| Option | Purpose |
|---|---|
| `Json.Config` | Preset: `vanillaCompatible`, `defaultConfig`, `fastestConfig`, `customConfig`. |
| `Json.Schema` | Map of schema name → file or URL path. |
| `Csv` | CSV-specific options. |

## Examples

### JSON Schema validation

```go
p := parser.InitParser(log, parser.Options{
    Json: parser.JsonOptions{
        Schema: map[string]string{ "user": "file://schemas/user.json" },
    },
})
if err := p.JSONParser().ValidateAgainstSchema("user", raw); err != nil {
    return errors.WrapWithCode(err, codes.CodeJSONSchemaInvalid, "user payload")
}
```

## Error Handling

Validation/parse errors are wrapped with [`codes`](../codes) parser codes.

## Dependencies

- **Internal:** [`codes`](../codes), [`errors`](../errors), [`logger`](../logger)
- **External:** `github.com/json-iterator/go`, `github.com/xeipuuv/gojsonschema`, `github.com/gocarina/gocsv`

## Testing

```bash
go test ./parser/...
```

JSON path is covered; CSV path currently has no tests.

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md).

## Related Packages

- [`auth`](../auth), [`messaging`](../messaging) — primary consumers.
- [`convert`](../convert) — primitive type conversion.
