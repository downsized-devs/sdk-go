# Package Registry

> Comprehensive index of every Go package in this monorepo. Stability matches the binding policy in [STABILITY.md](../STABILITY.md). "Last updated" reflects the most recent `git log` commit touching the package directory as of 2026-05-16.

All packages share the import root `github.com/downsized-devs/sdk-go`. The scaffolding CLI previously at `generator/` was extracted to a sibling repo, [`scaffolder-go`](https://github.com/downsized-devs/scaffolder-go).

## Index by Category

- **Configuration & bootstrap**: [appcontext](#appcontext) · [configbuilder](#configbuilder) · [configreader](#configreader) · [featureflag](#featureflag)
- **Logging, errors, observability**: [logger](#logger) · [errors](#errors) · [codes](#codes) · [audit](#audit) · [instrument](#instrument) · [tracker](#tracker)
- **Data & storage**: [sql](#sql) · [nosql](#nosql) · [redis](#redis) · [storage](#storage) · [localstorage](#localstorage) · [query](#query) · [null](#null)
- **Auth & security**: [auth](#auth) · [security](#security) · [ratelimiter](#ratelimiter)
- **Messaging & integrations**: [email](#email) · [messaging](#messaging) · [slack](#slack) · [gqlclient](#gqlclient)
- **I18n & locale**: [language](#language) · [translator](#translator)
- **Time & jobs**: [clock](#clock) · [dates](#dates) · [scheduler](#scheduler)
- **Files & documents**: [files](#files) · [pdf](#pdf) · [parser](#parser)
- **Utilities & primitives**: [character](#character) · [checker](#checker) · [convert](#convert) · [num](#num) · [operator](#operator) · [stringlib](#stringlib) · [header](#header)
- **Tooling**: [tests](#tests)

## Registry Table

| Package | Purpose | Key Features | Stability | Last Updated |
|---|---|---|---|---|
| <a id="appcontext"></a>**appcontext** | Request-scoped context value helpers | Setter/getter pairs for request ID, user ID, accept-language, service version, device type, response code | Stable | May 2026 |
| <a id="audit"></a>**audit** | Audit trail event capture | `Capture`/`Record` API; pulls request + user context from `appcontext` | Stable | Jun 2024 |
| <a id="auth"></a>**auth** | Firebase authentication client | Token verify/refresh, user CRUD, password sign-in, refresh-token revoke | Stable | May 2026 |
| <a id="character"></a>**character** | String casing & password-strength helpers | `CapitalizeFirstCharacter`, `IsStrongCharCombination` | Stable | Jun 2024 |
| <a id="checker"></a>**checker** | Generic validators | `ArrayContains`, `ArrayDeduplicate`, `IsEmail`, `IsPhoneNumber` (generic, no external deps) | Stable | Mar 2025 |
| <a id="clock"></a>**clock** | Timezone-aware clock with mockable `Now` | `GetCurrentTime`, `AddTime`, `SubstractTime`, `GetTimeInLocation`, first/last day of month | Stable | Feb 2025 |
| <a id="codes"></a>**codes** | Centralised error/success code registry | Reserved code ranges, bilingual `DisplayMessage` map, `Compile()` helper | Stable | May 2026 |
| <a id="configbuilder"></a>**configbuilder** | Mustache-template config file generator | Renders `*.tmpl` to runtime config files; viper-aware | Stable | May 2026 |
| <a id="configreader"></a>**configreader** | Layered configuration reader | JSON-ref resolution, viper-backed, custom duration decode hooks | Stable | May 2026 |
| <a id="convert"></a>**convert** | Type conversion utilities | Int/float/string conversion, camel/pascal case, roman numerals | Stable | Jul 2025 |
| <a id="dates"></a>**dates** | Date arithmetic helpers | `Difference` (day-level diff between two times) | Stable | Jun 2024 |
| <a id="email"></a>**email** | SMTP email sender with MJML templating | `SendEmail`, `GenerateBody`, `FromHTML`, `FromMJML` | Stable | Apr 2026 |
| <a id="errors"></a>**errors** | Error wrapping with codes & stack traces | `NewWithCode`, `WrapWithCode`, `Compile`, `GetCode`, `GetCaller`, `Is`/`As` | Stable | May 2026 |
| <a id="featureflag"></a>**featureflag** | Wrapper around `go-feature-flag` | `CheckUserFlags`, `GetAllUserFlags`, `Refresh` | Stable | May 2026 |
| <a id="files"></a>**files** | Filesystem helpers | `GetExtension`, `IsExist` | Stable | Jun 2024 |
| <a id="gqlclient"></a>**gqlclient** | Low-level GraphQL HTTP client | JSON and multipart `Run`; `WithHTTPClient`, `UseMultipartForm` options | Stable | May 2026 |
| <a id="header"></a>**header** | HTTP header & MIME constants | ~18 string constants (content types, cache control, header keys) | Stable | Jun 2024 |
| <a id="instrument"></a>**instrument** | Prometheus metrics for HTTP, DB, scheduler | `MetricsHandler`, `HTTPRequestTimer`/`Counter`, `RegisterDBStats`, `DatabaseQueryTimer`, `SchedulerRunningTimer`/`Counter` | Stable | May 2026 |
| <a id="language"></a>**language** | Locale constants + HTTP status text | EN/ID/JA/DE constants; `HTTPStatusText(lang, code)` | Stable | May 2026 |
| <a id="localstorage"></a>**localstorage** | Bleve-backed full-text local index | `NewIndex`, `Index`, `Search`, `DeleteIndex` | Stable | May 2026 |
| <a id="logger"></a>**logger** | Structured logging on zerolog | Trace/Debug/Info/Warn/Error/Fatal/Panic, `Debugf`, context-field extraction | Stable | May 2026 |
| <a id="messaging"></a>**messaging** | Firebase Cloud Messaging | `SubscribeToTopic`, `UnsubscribeFromTopic`, `BroadcastToTopic`, `BatchSendDryRun` | Stable | May 2026 |
| <a id="nosql"></a>**nosql** | MongoDB wrapper | `Find`, `FindOne`, `InsertOne`, `UpdateOne`, `UpdateMany`, `Close` | Stable | May 2026 |
| <a id="null"></a>**null** | SQL-nullable JSON-friendly types | `Bool`, `Int64`, `Float64`, `String`, `Time` with `SqlNull` flag | Stable | Mar 2025 |
| <a id="num"></a>**num** | Numeric & matrix utilities | `SafeDivision`, `RandomString`, `RoundFloat`, `ExcelGenerateCoords`, `EmptyStringSlice` | Stable | Mar 2025 |
| <a id="operator"></a>**operator** | Bitwise & ternary helpers | `CheckBitOnPosition`, generic `Ternary[T comparable]` | Stable | May 2026 |
| <a id="parser"></a>**parser** | JSON + CSV parsing with schema validation | `JsonInterface` (5 marshal/unmarshal variants), `CsvInterface`, JSON-schema enforcement | Stable | Apr 2026 |
| <a id="pdf"></a>**pdf** | PDF manipulation | `Encrypt`, `RemovePassword`, `Merge`, `Split`, `AddTextWatermark`, `ExtractText`, `PageCount` | Stable | May 2026 |
| <a id="query"></a>**query** | SQL query/clause builder | Struct-tag-driven WHERE/ORDER builder, cursor pagination, typed converters | Stable | May 2026 |
| <a id="ratelimiter"></a>**ratelimiter** | Gin rate-limiting middleware | Per-path `ConfigPath`, `GinMiddleware`, ulule/limiter backend | Stable | Jun 2024 |
| <a id="redis"></a>**redis** | Redis client with distributed locks | `Get`, `SetEX`, `Lock`/`LockRelease` (redislock), `Del`, `Flush*`, `Ping`, `CRC16` | Stable | May 2026 |
| <a id="scheduler"></a>**scheduler** | gocron v2 wrapper | `Register` with duration/daily/weekly/monthly job types, `Start`/`Shutdown` | Stable | May 2026 |
| <a id="security"></a>**security** | Cryptographic primitives | AES-GCM encrypt/decrypt, PBKDF2, Scrypt password hashing, HMAC | Stable | May 2026 |
| <a id="slack"></a>**slack** | Slack message sender | `SendMessage` with attachments and attachment fields | Stable | Jun 2024 |
| <a id="sql"></a>**sql** | SQL DB abstraction with leader/follower | Multi-driver (MySQL/Postgres/SQLite), prepared statements, transactions, instrumentation | Stable | Apr 2026 |
| <a id="storage"></a>**storage** | AWS S3 wrapper | `Upload`, `Download`, `Delete`, `GetPresignedUrl[WithDuration]`, `CreateUrlByKey` | Stable | Jun 2024 |
| <a id="stringlib"></a>**stringlib** | Misc string utilities | `RandStringBytes` | Stable | Apr 2026 |
| <a id="tests"></a>**tests** | Shared gomock mocks for SDK packages | No top-level Go code; `tests/mock/<pkg>/` directories with generated mocks | Stable | May 2026 |
| <a id="tracker"></a>**tracker** | Prometheus push gateway + webhook | `Push`, `PushWebhook` with `Options`/`WebhookOptions` | Stable | Apr 2026 |
| <a id="translator"></a>**translator** | i18n via universal-translator | `Translate(ctx, key, params)`, EN/ID locale registration | Stable | Apr 2026 |

## Import Paths

```go
import (
    "github.com/downsized-devs/sdk-go/appcontext"
    "github.com/downsized-devs/sdk-go/audit"
    "github.com/downsized-devs/sdk-go/auth"
    "github.com/downsized-devs/sdk-go/character"
    "github.com/downsized-devs/sdk-go/checker"
    "github.com/downsized-devs/sdk-go/clock"
    "github.com/downsized-devs/sdk-go/codes"
    "github.com/downsized-devs/sdk-go/configbuilder"
    "github.com/downsized-devs/sdk-go/configreader"
    "github.com/downsized-devs/sdk-go/convert"
    "github.com/downsized-devs/sdk-go/dates"
    "github.com/downsized-devs/sdk-go/email"
    "github.com/downsized-devs/sdk-go/errors"
    "github.com/downsized-devs/sdk-go/featureflag"
    "github.com/downsized-devs/sdk-go/files"
    "github.com/downsized-devs/sdk-go/gqlclient"
    "github.com/downsized-devs/sdk-go/header"
    "github.com/downsized-devs/sdk-go/instrument"
    "github.com/downsized-devs/sdk-go/language"
    "github.com/downsized-devs/sdk-go/localstorage"
    "github.com/downsized-devs/sdk-go/logger"
    "github.com/downsized-devs/sdk-go/messaging"
    "github.com/downsized-devs/sdk-go/nosql"
    "github.com/downsized-devs/sdk-go/null"
    "github.com/downsized-devs/sdk-go/num"
    "github.com/downsized-devs/sdk-go/operator"
    "github.com/downsized-devs/sdk-go/parser"
    "github.com/downsized-devs/sdk-go/pdf"
    "github.com/downsized-devs/sdk-go/query"
    "github.com/downsized-devs/sdk-go/ratelimiter"
    "github.com/downsized-devs/sdk-go/redis"
    "github.com/downsized-devs/sdk-go/scheduler"
    "github.com/downsized-devs/sdk-go/security"
    "github.com/downsized-devs/sdk-go/slack"
    "github.com/downsized-devs/sdk-go/sql"
    "github.com/downsized-devs/sdk-go/storage"
    "github.com/downsized-devs/sdk-go/stringlib"
    "github.com/downsized-devs/sdk-go/tracker"
    "github.com/downsized-devs/sdk-go/translator"
)
```

## See also

- [DEPENDENCY_GRAPH.md](./DEPENDENCY_GRAPH.md) — internal/external deps and critical-path packages.
- [STABILITY.md](../STABILITY.md) — what each stability level promises.
- [PACKAGE_README_TEMPLATE.md](./PACKAGE_README_TEMPLATE.md) — template for per-package READMEs.
