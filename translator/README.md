# `translator` — i18n via universal-translator

`import "github.com/downsized-devs/sdk-go/translator"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

Wraps [`go-playground/universal-translator`](https://github.com/go-playground/universal-translator) with English and Indonesian locales pre-registered. Reads the user's locale from [`appcontext`](../appcontext).

## Features

- `Translate(ctx, key, params)` — locale-aware translation.
- English + Indonesian locales registered out of the box.

## Installation

```bash
go get github.com/downsized-devs/sdk-go/translator
```

## Quick Start

```go
tr := translator.Init(translator.Config{
    Translations: []translator.Translation{
        { Key: "welcome", En: "Welcome, {0}", Id: "Selamat datang, {0}" },
    },
}, log)

ctx := appcontext.SetAcceptLanguage(context.Background(), language.Indonesian)
msg, _ := tr.Translate(ctx, "welcome", "Alice")
// msg == "Selamat datang, Alice"
```

## API Reference

| Symbol | Signature |
|---|---|
| `Init` | `func Init(cfg Config, log logger.Interface) Interface` |
| `Interface.Translate` | `(ctx, key string, params ...any) (string, error)` |
| `Config` | `{ Translations []Translation }` |
| `Translation` | `{ Key, En, Id string }` |

## Error Handling

Unknown keys return a coded error from [`codes`](../codes). Decide per-call whether to fall back to a literal.

## Dependencies

- **Internal:** [`appcontext`](../appcontext), [`codes`](../codes), [`errors`](../errors), [`language`](../language), [`logger`](../logger)
- **External:** `github.com/go-playground/locales`, `.../locales/en`, `.../locales/id`, `github.com/go-playground/universal-translator`

## Testing

```bash
go test ./translator/...
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). Adding a locale: register it in `translator.go` and extend `Translation` to include the new field.

## Related Packages

- [`language`](../language) — locale constants.
- [`appcontext`](../appcontext) — provides `AcceptLanguage`.
- [`codes`](../codes), [`errors`](../errors) — bilingual error messages.
