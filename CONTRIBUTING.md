# Contributing to `sdk-go`

Welcome. This guide is for anyone adding a new package, extending an existing one, or fixing a bug across the monorepo.

## Before you start

1. Read the [Package Registry](./docs/PACKAGE_REGISTRY.md) — your work may already exist or live next to a similar package.
2. Check the [Stability Policy](./STABILITY.md) — Stable packages have a higher bar for changes.
3. Check the [Dependency Graph](./docs/DEPENDENCY_GRAPH.md) — breaking changes to `logger`, `errors`, or `codes` cascade across most of the SDK.

## Setup

```bash
git clone https://github.com/downsized-devs/sdk-go
cd sdk-go
go mod download
make build
make run-tests
```

Required tooling:

- Go ≥ the version pinned in `go.mod` (currently `1.25`).
- `golangci-lint` — config in `.golangci.yml`.
- `mockgen` — used by `make mock-all`.

## Code style

- **`gofmt -s`** is mandatory. CI fails on unformatted files.
- **`golangci-lint run`** must pass. We use the bundled linters in `.golangci.yml`; do not disable lints in code without a `//nolint:<name> // <reason>` comment.
- **Effective Go + Go Code Review Comments** are the baseline. Notably:
  - Exported symbols need doc comments that start with the symbol name.
  - Errors are values — return them, don't panic, except in `init()` or true logic invariants.
  - Receivers are short (`c *cache`, not `cache *Cache`).

## Package structure conventions

Every package follows the same shape:

```
<pkg>/
    <pkg>.go             # package doc, Interface, Config, Init, struct
    <pkg>_test.go        # table-driven unit tests
    entity.go            # types only (optional)
    <feature>.go         # additional grouped methods (optional)
    README.md            # use ../docs/PACKAGE_README_TEMPLATE.md
```

Conventions:

- One package per top-level directory. Do not nest `internal/` unless you have a real reason.
- Expose an `Interface` and an `Init(cfg Config, deps…) Interface` factory. This makes mocking with `mockgen` trivial.
- `Config` is a plain struct; do not validate inside `Init` — return descriptive errors at the caller's `cfg.Validate()` step or fail at first use.
- Keep external dependencies behind the package boundary; do not re-export third-party types unless that's the package's reason for existing (e.g., `redis.Nil`).

## Adding a new package

1. Create `/<pkgname>/`.
2. Add `<pkgname>.go` with `package <pkgname>`, an `Interface`, a `Config`, and `Init`.
3. Add `<pkgname>_test.go` with table-driven tests covering the happy path and at least one error case per method.
4. Copy [PACKAGE_README_TEMPLATE.md](./docs/PACKAGE_README_TEMPLATE.md) to `<pkgname>/README.md` and fill it in.
5. Add a row to [PACKAGE_REGISTRY.md](./docs/PACKAGE_REGISTRY.md).
6. Add a row to [DEPENDENCY_GRAPH.md](./docs/DEPENDENCY_GRAPH.md) (matrix + mermaid if non-trivial).
7. Mark stability **Experimental** until tests + downstream usage justify Beta or Stable.

## Testing

- **Unit tests** live next to the package as `*_test.go`.
- **Cross-package integration tests** live in `/tests/`.
- Use **table-driven tests** with `name`, `args`, `want`, `wantErr` fields. Match the style already in `num/division_test.go`, `clock/clock_test.go`, etc.
- Use `github.com/stretchr/testify/assert` for assertions; `go.uber.org/mock` for mocks.
- Don't use `time.Sleep` in tests — inject a `clock.Interface` or override the `now` variable like `num/random.go` does.

### Coverage thresholds

- **Stable packages**: ≥ 80% line coverage.
- **Beta packages**: ≥ 60% line coverage.
- **Experimental**: best effort.

Coverage is reported to Codecov; see `.codecov.yml` for thresholds enforced at the PR level.

### Regenerating mocks

```bash
make mock-all
```

If you add or change an `Interface`, regenerate mocks in the same PR so downstream code stays buildable.

## Documentation requirements

Every PR must update documentation if it changes any of:

- A public symbol → update the package `README.md` (API Reference section).
- An `Interface` → also update [STABILITY.md](./STABILITY.md) if the symbol is Stable.
- Internal imports → update [DEPENDENCY_GRAPH.md](./docs/DEPENDENCY_GRAPH.md).
- A package's purpose or scope → update [PACKAGE_REGISTRY.md](./docs/PACKAGE_REGISTRY.md).

## Performance considerations

- **Allocate sparingly** in hot paths (HTTP request handlers, log lines). Prefer pre-sized slices and `strings.Builder` over `+=`.
- **Don't reflect** unless the alternative is significantly worse. `query/sql_builder.go` uses reflection deliberately for ergonomics; weigh the cost.
- **Benchmark** before optimizing — `go test -bench=. -benchmem ./<pkg>/...`. Include a `Benchmark*` function when adding a perf-sensitive path.
- **Avoid global state.** Goroutine-safe `Interface` implementations are the norm; mutable package-level vars are usually a bug waiting to happen.

## Security considerations

- Never log credentials, tokens, or PII. `logger` does not redact for you.
- Crypto code lives in [`security`](../security). Don't roll your own AES/HMAC elsewhere — import it.
- When touching `auth`, `security`, or anything that handles tokens, request a second reviewer.
- The CodeQL workflow at `.github/workflows/codeql.yml` runs on every push; treat its warnings as blocking.

## Pull request checklist

Before requesting review:

- [ ] `gofmt -s -w .` clean
- [ ] `golangci-lint run ./...` passes
- [ ] `go test ./<pkg>/...` passes locally
- [ ] Coverage at or above the threshold for the package's stability level
- [ ] Package README updated if public API changed
- [ ] [PACKAGE_REGISTRY.md](./docs/PACKAGE_REGISTRY.md) updated if scope/stability changed
- [ ] [DEPENDENCY_GRAPH.md](./docs/DEPENDENCY_GRAPH.md) updated if imports changed
- [ ] [STABILITY.md](./STABILITY.md) updated if you deprecated, promoted, or demoted a package
- [ ] PR description explains the *why*, not just the *what*

## Commit messages

Short, imperative, scoped:

```
redis: add MGet with cluster-slot safety
logger: extract caller from wrapped errors
docs: refresh dependency graph after sql split
```

Group related changes per commit; avoid mega-commits that span unrelated packages.

## Release process

1. Open a release PR that bumps the version in any version-bearing files and finalises a `CHANGELOG.md` entry.
2. After merge, tag: `git tag vX.Y.Z && git push origin vX.Y.Z`.
3. GitHub release auto-publishes from the tag.
4. For majors, ship `docs/MIGRATION-vX.md` describing renames and removals before tagging.

## Reporting bugs & security issues

- **Public bugs**: open a GitHub issue with a minimal reproducer.
- **Security vulnerabilities**: do *not* file a public issue. Contact a maintainer directly first.

## See also

- [PACKAGE_REGISTRY.md](./docs/PACKAGE_REGISTRY.md)
- [DEPENDENCY_GRAPH.md](./docs/DEPENDENCY_GRAPH.md)
- [STABILITY.md](./STABILITY.md)
- [PACKAGE_README_TEMPLATE.md](./docs/PACKAGE_README_TEMPLATE.md)
