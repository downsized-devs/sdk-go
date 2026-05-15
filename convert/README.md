# `convert` — type conversion utilities

`import "github.com/downsized-devs/sdk-go/convert"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Converts between primitive types, string-case styles, and roman numerals.

## Features

- `ToInt64`, `ToFloat64`, `ToString`, `ToArrInt64` — primitive conversion.
- `IntToChar` — integer to Excel-style column letter.
- `ToCamelCase`, `ToPascalCase`, `PascalCaseToCamelCase`.
- `RomanToInt64`, `Int64ToRoman` — roman numerals.

## Installation

```bash
go get github.com/downsized-devs/sdk-go/convert
```

## Quick Start

```go
n, _ := convert.ToInt64("42")        // 42
s := convert.ToCamelCase("hello_world") // "helloWorld"
r := convert.Int64ToRoman(2026)        // "MMXXVI"
```

## API Reference

| Symbol | Signature |
|---|---|
| `ToInt64` | `(any) (int64, error)` |
| `ToFloat64` | `(any) (float64, error)` |
| `ToString` | `(any) (string, error)` |
| `ToArrInt64` | `(any) ([]int64, error)` |
| `IntToChar` | `(int) string` |
| `ToCamelCase` / `ToPascalCase` / `PascalCaseToCamelCase` | `(string) string` |
| `RomanToInt64` | `(string) (int64, error)` |
| `Int64ToRoman` | `(int64) string` |

## Error Handling

Numeric conversions wrap with codes from [`codes`](../codes) (parse-error range) via [`errors`](../errors).

## Dependencies

- **Internal:** [`codes`](../codes), [`errors`](../errors)
- **External:** `github.com/cstockton/go-conv`

## Testing

```bash
go test ./convert/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md).

## Related Packages

- [`null`](../null) — when conversions might produce SQL nulls.
- [`parser`](../parser) — higher-level JSON/CSV parsing.
