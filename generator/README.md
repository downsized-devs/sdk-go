# `generator` — code scaffolding CLI

`import "github.com/downsized-devs/sdk-go/generator"` (CLI, not a library)

**Stability:** Experimental — CLI tool, not a Go library. See [STABILITY.md](../STABILITY.md).

Generates boilerplate (entity, repository, usecase, HTTP handler) for new services that follow the Downsized Devs layered template.

## Installation

Run directly from the module:

```bash
go run github.com/downsized-devs/sdk-go/generator \
    --name MyEntity \
    --path ./services/my-service \
    --api
```

Or from a local clone:

```bash
cd generator
go run main.go -entity_name=$entity_name -file_location=$file_location -api=$api
```

## Flags

| Flag | Purpose |
|---|---|
| `-entity_name` | Entity name in CamelCase (e.g. `Invoice`). |
| `-file_location` | Output directory — the generic-service path in your project. |
| `-api` | Comma-separated API actions: `create,edit,delete,get,activate`. Empty = all. |

## Layout

```
generator/
    main.go         # CLI entry point
    helper/         # filename + path helpers
    modifiers/      # AST transforms
    services/       # template orchestration
    templates/      # Go file templates + naming map
```

## Testing

```bash
go test ./generator/...
```

Tests live in `helper/`, `modifiers/`, and `services/` subpackages.

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). Template edits require regenerating sample output to verify the result still compiles.

## Related Packages

- [`stringlib`](../stringlib), [`convert`](../convert) — case helpers used in templates.
