# `tests` — shared test mocks

`import "github.com/downsized-devs/sdk-go/tests/mock/<pkg>"`

**Stability:** Stable — see [STABILITY.md](../STABILITY.md)

This directory does not contain a Go package of its own. It hosts `mock/<pkg>/` subdirectories with `gomock`-generated mocks for the public `Interface` of every other package in the monorepo.

## Layout

```
tests/
    mock/
        audit/
        auth/
        configbuilder/
        configreader/
        email/
        featureflag/
        gqlclient/
        instrument/
        local_storage/
        localstorage/
        logger/
        messaging/
        parser/
        pdf/
        query/
        ratelimiter/
        redis/
        security/
        slack/
        sql/
        … (one directory per mockable package)
```

## Using a mock in your tests

```go
import (
    "go.uber.org/mock/gomock"
    mock_logger "github.com/downsized-devs/sdk-go/tests/mock/logger"
)

ctrl := gomock.NewController(t)
defer ctrl.Finish()
log := mock_logger.NewMockInterface(ctrl)
log.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
```

## Regenerating mocks

From the repo root:

```bash
make mock-all
```

This re-runs `mockgen` for every `Interface` in the SDK. Commit the regenerated files alongside any `Interface` changes.

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md). If you add a new package with an `Interface`, add a `mockgen` invocation to the Makefile and regenerate.

## Related Packages

Every package with an `Interface` has a corresponding mock directory here.
