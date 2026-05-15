# `num` — numeric & matrix utilities

`import "github.com/downsized-devs/sdk-go/num"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

## Features

- `SafeDivision(numerator, denominator, zeroValue)` — guarded division.
- `RandomString(n)` — alphanumeric random string of length `n`.
- `RoundFloat` — round to N decimal places (see `rounding.go`).
- `EmptyStringSlice(length)` — pre-sized `[]string{"", "", ...}`.
- `ExcelGenerateCoords(ctx, range)` — expand `D10:D13` into `[D10, D11, D12, D13]`.

## Installation

```bash
go get github.com/downsized-devs/sdk-go/num
```

## Quick Start

```go
import "github.com/downsized-devs/sdk-go/num"

result := num.SafeDivision(10, 0, true)        // 0
random := num.RandomString(12)                  // "aB3kX..."
cells, _ := num.ExcelGenerateCoords(ctx, "A1:A3") // ["A1","A2","A3"]
```

## API Reference

| Symbol | Signature |
|---|---|
| `SafeDivision` | `(numerator, denominator float64, zeroValue bool) float64` |
| `RandomString` | `(n int) string` |
| `RoundFloat` | see `rounding.go` |
| `EmptyStringSlice` | `(length int) []string` |
| `ExcelGenerateCoords` | `(ctx, range string) ([]string, error)` |
| `now` | private package var (`time.Now`), overridable in tests. |

## Dependencies

- **External:** `github.com/xuri/excelize/v2`

## Testing

```bash
go test ./num/...
```

Five test files; `RandomString` test overrides the `now` var for determinism.

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md).

## Related Packages

- [`stringlib`](../stringlib) — older random-string helper.
- [`convert`](../convert) — numeric type conversion.
