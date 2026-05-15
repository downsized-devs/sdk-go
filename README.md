# sdk-go

> A monorepo of Go libraries by **Downsized Devs** — logging, config, auth, scheduling, storage, and ~35 other small, well-scoped packages designed to be imported individually.

![Go version](https://img.shields.io/github/go-mod/go-version/downsized-devs/sdk-go)
![Build](https://img.shields.io/github/actions/workflow/status/downsized-devs/sdk-go/go.yml)
![Coverage](https://img.shields.io/codecov/c/github/downsized-devs/sdk-go)
![Version](https://img.shields.io/github/v/release/downsized-devs/sdk-go)
![License](https://img.shields.io/github/license/downsized-devs/sdk-go)

## Installation

```bash
go get github.com/downsized-devs/sdk-go
```

Or pull a single package — Go resolves the module graph minimally:

```bash
go get github.com/downsized-devs/sdk-go/logger
```

## Quick Navigation

Pick the entry-point package for what you're trying to build. The full catalogue lives in [docs/PACKAGE_REGISTRY.md](./docs/PACKAGE_REGISTRY.md).

| Category | Packages |
|---|---|
| **Configuration & bootstrap** | [`appcontext`](./appcontext), [`configbuilder`](./configbuilder), [`configreader`](./configreader), [`featureflag`](./featureflag) |
| **Logging, errors, observability** | [`logger`](./logger), [`errors`](./errors), [`codes`](./codes), [`audit`](./audit), [`instrument`](./instrument), [`tracker`](./tracker) |
| **Data & storage** | [`sql`](./sql), [`nosql`](./nosql), [`redis`](./redis), [`storage`](./storage), [`local_storage`](./local_storage), [`query`](./query), [`null`](./null) |
| **Auth & security** | [`auth`](./auth), [`security`](./security), [`ratelimiter`](./ratelimiter) |
| **Messaging & integrations** | [`email`](./email), [`messaging`](./messaging), [`slack`](./slack), [`gqlclient`](./gqlclient) |
| **I18n & locale** | [`language`](./language), [`translator`](./translator) |
| **Time & jobs** | [`clock`](./clock), [`dates`](./dates), [`scheduler`](./scheduler) |
| **Files & documents** | [`files`](./files), [`pdf`](./pdf), [`parser`](./parser) |
| **Primitives & helpers** | [`character`](./character), [`checker`](./checker), [`convert`](./convert), [`num`](./num), [`operator`](./operator), [`stringlib`](./stringlib), [`header`](./header) |
| **Tooling** | [`tests`](./tests) (gomock fixtures). The scaffolding CLI is now [`scaffolder-go`](https://github.com/downsized-devs/scaffolder-go). |

## Common Use Cases

| I need… | Reach for… |
|---|---|
| Structured logging that follows my request context | [`logger`](./logger) + [`appcontext`](./appcontext) |
| Typed errors with codes that survive across HTTP boundaries | [`errors`](./errors) + [`codes`](./codes) |
| A SQL database with leader/follower routing | [`sql`](./sql) (use [`query`](./query) for dynamic clause building) |
| Redis caching with distributed locks | [`redis`](./redis) |
| MongoDB CRUD | [`nosql`](./nosql) |
| S3 upload/download with presigned URLs | [`storage`](./storage) |
| Background cron jobs | [`scheduler`](./scheduler) (lock across replicas with [`redis`](./redis)) |
| Firebase authentication | [`auth`](./auth) |
| AES encryption / password hashing | [`security`](./security) |
| Per-route HTTP rate limiting (Gin) | [`ratelimiter`](./ratelimiter) |
| Prometheus metrics for HTTP/DB/jobs | [`instrument`](./instrument) |
| Send transactional email with MJML templates | [`email`](./email) |
| Send Slack notifications | [`slack`](./slack) |
| Multi-language responses | [`language`](./language) + [`translator`](./translator) |
| Generate API boilerplate | `go run github.com/downsized-devs/scaffolder-go --entity_name X --file_location ./svc --api "create,edit,get,delete"` |

## Minimal example

```go
package main

import (
    "context"

    "github.com/downsized-devs/sdk-go/appcontext"
    "github.com/downsized-devs/sdk-go/logger"
)

func main() {
    log := logger.Init(logger.Config{Level: "info"})
    ctx := appcontext.SetRequestId(context.Background(), "boot")
    log.Info(ctx, "hello, sdk-go")
}
```

## Repository layout

Every top-level directory is a self-contained package. Internal dependencies form a DAG with `logger`, `codes`, and `errors` as the most-depended-on nodes — see [docs/DEPENDENCY_GRAPH.md](./docs/DEPENDENCY_GRAPH.md).

```
/<package>/         # one package per directory (40 in total)
/docs/              # cross-cutting documentation
/tests/mock/        # gomock fixtures, one subdir per mockable package
Makefile, go.mod    # repository root
```

## Documentation

- **[Contributing Guide](./CONTRIBUTING.md)** — code style, test thresholds, release process.
- **[Stability Policy](./STABILITY.md)** — per-package maturity levels, breaking-change policy.
- **[Package Registry](./docs/PACKAGE_REGISTRY.md)** — full catalogue with purpose, features, stability, import paths.
- **[Dependency Graph](./docs/DEPENDENCY_GRAPH.md)** — mermaid + matrices + critical paths.
- **[Package README Template](./docs/PACKAGE_README_TEMPLATE.md)** — for authors adding `<pkg>/README.md`.
- Looking for worked examples? Every package has its own README; start with [`logger`](./logger/README.md), [`auth`](./auth/README.md), or [`redis`](./redis/README.md).

## Migration & breaking changes

From `v1.0.0` onward, every package is **Stable** and follows the binding semver policy in [STABILITY.md](./STABILITY.md). Each major release ships a `MIGRATION-vX.md` in `docs/` describing renamed symbols, removed packages, and field reshapes.

If a breaking change is unavoidable in a minor release, expect:

1. A deprecation notice in code (`// Deprecated: …`) and in [STABILITY.md](./STABILITY.md).
2. The new symbol shipped alongside the old one for one minor cycle.
3. Removal only at the next major version.

## Roadmap

- **v1.0** — every package Stable; semver becomes binding (see [STABILITY.md](./STABILITY.md)).
- **Post-v1** — additive only on Stable interfaces. Deprecations carry one minor cycle before removal at the next major.

## Code scaffolder

The code-scaffolding CLI lives in its own repo: [`scaffolder-go`](https://github.com/downsized-devs/scaffolder-go).

```bash
go run github.com/downsized-devs/scaffolder-go \
    --entity_name <EntityName> \
    --file_location <output-path> \
    --api "create,edit,get,activate,delete"
```

## Tooling

```bash
make build         # compile every package
make run-tests     # full unit-test suite
make mock-all      # regenerate gomock stubs across packages
```

CI is GitHub Actions; coverage is reported to Codecov; CodeQL runs from `.github/workflows/codeql.yml`.

## License

MIT — see [LICENSE](./LICENSE).

## Support

[Open an issue](https://github.com/downsized-devs/sdk-go/issues) for bugs and feature requests. Security-sensitive issues should be reported privately — contact a maintainer before filing publicly.

---

## Quality Metrics

![GitHub Issues](https://img.shields.io/github/issues/downsized-devs/sdk-go)
![GitHub Pull Requests](https://img.shields.io/github/issues-pr/downsized-devs/sdk-go)
![GitHub License](https://img.shields.io/github/license/downsized-devs/sdk-go)
![Code Quality](https://img.shields.io/codefactor/grade/github/downsized-devs/sdk-go)

## Repository Stats

![GitHub Contributors](https://img.shields.io/github/contributors/downsized-devs/sdk-go)
![GitHub Last Commit](https://img.shields.io/github/last-commit/downsized-devs/sdk-go)
![Repo Size](https://img.shields.io/github/repo-size/downsized-devs/sdk-go)
![Language](https://img.shields.io/github/languages/top/downsized-devs/sdk-go)

## Code Health

![Go Report Card](https://goreportcard.com/badge/github.com/downsized-devs/sdk-go)

## Community

![GitHub Stars](https://img.shields.io/github/stars/downsized-devs/sdk-go?style=social)
![GitHub Forks](https://img.shields.io/github/forks/downsized-devs/sdk-go?style=social)
