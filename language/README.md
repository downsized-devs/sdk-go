# `language` — locale constants + HTTP status text

`import "github.com/downsized-devs/sdk-go/language"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Locale identifiers (English, Indonesian, Japanese, Deutsch) and translated HTTP status text. Used by [`codes`](../codes) and [`translator`](../translator).

## Features

- Locale constants: `English`, `Indonesian`, `Japanese`, `Deutsch`.
- `HTTPStatusText(lang string, code int) string`.

## Installation

```bash
go get github.com/downsized-devs/sdk-go/language
```

## Quick Start

```go
import "github.com/downsized-devs/sdk-go/language"

txt := language.HTTPStatusText(language.Indonesian, 404) // "Tidak ditemukan"
```

## API Reference

| Symbol | Signature |
|---|---|
| `English`, `Indonesian`, `Japanese`, `Deutsch` | `const string` |
| `HTTPStatusText` | `(lang string, code int) string` |

## Dependencies

stdlib only.

## Testing

```bash
go test ./language/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). Add new locales additively — never re-purpose an existing constant.

## Related Packages

- [`codes`](../codes), [`errors`](../errors) — render messages in the chosen locale.
- [`translator`](../translator) — general-purpose i18n.
