# Stability Policy

> How we version `sdk-go`, what each stability level promises, and how breaking changes are introduced. Stability levels here were **inferred from code signals** (test presence, exported `Interface`, TODO markers) — maintainers should review and adjust before treating this as policy.

## Versioning

`sdk-go` ships as a **single Go module** at the repository root. Every package shares the same module version; you cannot pin packages independently.

We follow [semantic versioning](https://semver.org):

- **MAJOR** — breaking changes to any **Stable** package's public API.
- **MINOR** — additive changes to Stable packages; any change (including breaking) to **Beta** or **Experimental** packages.
- **PATCH** — bug fixes only.

> No v1.0.0 has been cut yet. Until then, every release is `0.y.z` and the project reserves the right to make breaking changes in minor versions across all packages — though we will avoid this in practice for the packages marked Stable below.

## Stability levels

| Level | What it means | Breaking-change cadence |
|---|---|---|
| **Stable** | Public `Interface` is frozen; covered by tests; relied on internally and externally. | Only at major version bumps after v1.0. |
| **Beta** | API is mostly settled but may shift. Tests partial or absent. Suitable for production with care. | May break in minor versions. |
| **Experimental** | Early-stage or single-use. API will change. | May break or be removed without notice. |
| **Deprecated** | Replaced by another package or symbol. Kept for one major version then removed. | Receives only critical bug fixes. |

## Per-package stability

Stability below is **inferred from code signals only** — presence of tests, presence of an exported `Interface`, and any in-source `TODO`/`FIXME` markers. Maintainers should override any row that doesn't match their actual stance.

### Stable

`appcontext`, `audit`, `auth`, `character`, `checker`, `clock`, `codes`, `configbuilder`, `configreader`, `convert`, `dates`, `email`, `errors`, `files`, `gqlclient`, `header`, `instrument`, `language`, `local_storage`, `logger`, `null`, `num`, `operator`, `parser`, `ratelimiter`, `redis`, `security`, `slack`, `sql`, `storage`, `stringlib`, `tracker`, `translator`

Signals: exported `Interface`, ≥1 test file with meaningful coverage, no in-source deprecation markers, used by ≥2 other packages or by external consumers.

### Beta

| Package | Why Beta (signal) |
|---|---|
| `featureflag` | No test files in package; wraps a fast-moving upstream (`go-feature-flag`). |
| `messaging` | No test files in package; Firebase Cloud Messaging surface may grow. |
| `nosql` | No test files in package; MongoDB API surface is intentionally small but not exercised. |
| `pdf` | Single-feature (`SetPasswordFile`); minimal test; expect more methods. |
| `query` | `// I hate this, find a better way for insert many rows` TODO in `query/query.go`. The clause-builder side is mature but the bulk-insert helper is on notice. |
| `scheduler` | No test files in package; gocron upstream still evolving. |

### Experimental

| Package | Why Experimental (signal) |
|---|---|
| `generator` | CLI tool, not a library. Public surface is the flags, not Go symbols. Templates change with consumer needs. |

### Deprecated

None at this time.

## Breaking-change policy

A "breaking change" to a Stable package means any of:

- Removing or renaming an exported symbol.
- Changing a function signature (parameters or return values).
- Adding a method to an `Interface` (downstream mocks break).
- Changing semantics in a way that compiles but behaves differently.

For Stable packages, breaking changes follow this dance:

1. **Deprecate.** Add `// Deprecated: use NewThing instead.` on the symbol. Add a row to the [Deprecations](#current-deprecations) table below in the same PR.
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

This is the *intended* policy after v1.0. Before v1.0 we make no support guarantees and recommend pinning to a specific tag.

## How to propose a stability change

Open a PR that:

1. Edits this file to move the package between sections.
2. Includes the signals justifying the move (coverage %, downstream usage, time-since-last-breaking-change).
3. Updates [PACKAGE_REGISTRY.md](./docs/PACKAGE_REGISTRY.md) so the registry table agrees.

## See also

- [PACKAGE_REGISTRY.md](./docs/PACKAGE_REGISTRY.md) — short-form stability column.
- [DEPENDENCY_GRAPH.md](./docs/DEPENDENCY_GRAPH.md) — critical-path packages where breakage cascades.
- [CONTRIBUTING.md](./CONTRIBUTING.md) — release-process details.
