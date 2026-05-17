# Stability Policy

> How we version `sdk-go`, what each stability level promises, and how breaking changes are introduced.

## Versioning

`sdk-go` ships as a **single Go module** at the repository root. Every package shares the same module version; you cannot pin packages independently.

We follow [semantic versioning](https://semver.org):

- **MAJOR** — breaking changes to any Stable package's public API. Allowed.
- **MINOR** — additive changes only to Stable packages. Old code keeps compiling.
- **PATCH** — bug fixes only. No surface changes.

Starting with **v1.0.0** this contract is **binding**: no breaking change ships outside a major version. Pre-v1 tags (`v0.y.z`) made no such promise.

## Stability levels

| Level | What it means | Breaking-change cadence |
|---|---|---|
| **Stable** | Public `Interface` is frozen; covered by tests; relied on internally and externally. | Only at major version bumps. |
| **Beta** | API is mostly settled but may shift. Suitable for production with care. | May break in minor versions. |
| **Experimental** | Early-stage or single-use. API will change. | May break or be removed without notice. |
| **Deprecated** | Replaced by another package or symbol. Kept for one major version then removed. | Receives only critical bug fixes. |

## Per-package stability (v1.0.0)

### Stable

`appcontext`, `audit`, `auth`, `character`, `checker`, `clock`, `codes`, `configbuilder`, `configreader`, `convert`, `dates`, `email`, `errors`, `featureflag`, `files`, `gqlclient`, `header`, `instrument`, `language`, `localstorage`, `logger`, `messaging`, `nosql`, `null`, `num`, `operator`, `parser`, `pdf`, `query`, `ratelimiter`, `redis`, `scheduler`, `security`, `slack`, `sql`, `storage`, `stringlib`, `tracker`, `translator`

All Stable packages share these traits:

- Exported `Interface` (or, where one package houses several seams: `JsonInterface`/`CsvInterface`/`TemplateInterface`) with a frozen signature.
- ≥1 test file with meaningful coverage. Packages whose tests require an external service (sql, slack, tracker) ship integration tests under `//go:build integration`.
- No `// TODO` / `// FIXME` markers on hot paths in the public API.

Adding a new public symbol is a minor bump; renaming or removing one is a major bump.

### Beta

_None._ The `pdf`, `query`, `featureflag`, `messaging`, `nosql`, and `scheduler` packages were promoted to Stable in v1.0 after their gaps (missing tests, in-flight rewrites) were closed.

### Experimental

_None._ The code-scaffolding CLI that previously lived under `generator/` was extracted to a sibling repo, [`scaffolder-go`](https://github.com/downsized-devs/scaffolder-go), before v1.

### Deprecated

_None._

## Breaking-change policy

A "breaking change" to a Stable package means any of:

- Removing or renaming an exported symbol.
- Changing a function signature (parameters or return values).
- Adding a method to an `Interface` (downstream mocks break).
- Changing semantics in a way that compiles but behaves differently.

For Stable packages, breaking changes follow this dance:

1. **Deprecate.** Add `// Deprecated: use NewThing instead.` on the symbol. Add a row to the [Current deprecations](#current-deprecations) table below in the same PR.
2. **Provide the replacement.** Ship the new symbol alongside the old one for at least one minor version.
3. **Remove.** Drop the deprecated symbol at the next major version.

Adding a new method to an existing `Interface` is the most common foot-gun — it silently breaks every mock. Prefer either (a) shipping the method on a separate `InterfaceV2` and embedding it, or (b) coordinating the change with all known downstream consumers before merging.

## Current deprecations

| Symbol | Deprecated in | Replacement | Removal target |
|---|---|---|---|
| _none_ | — | — | — |

## Support timeline

- The **latest minor** of the current major receives bug fixes and security patches.
- The **previous major** receives security patches only for 6 months after the next major ships.
- Older majors are end-of-life and will not receive any updates.

## How to propose a stability change

Open a PR that:

1. Edits this file to move the package between sections.
2. Includes the signals justifying the move (coverage %, downstream usage, time-since-last-breaking-change).
3. Updates [PACKAGE_REGISTRY.md](./docs/PACKAGE_REGISTRY.md) so the registry table agrees.

## See also

- [PACKAGE_REGISTRY.md](./docs/PACKAGE_REGISTRY.md) — short-form stability column.
- [DEPENDENCY_GRAPH.md](./docs/DEPENDENCY_GRAPH.md) — critical-path packages where breakage cascades.
- [CONTRIBUTING.md](./CONTRIBUTING.md) — release-process details.
