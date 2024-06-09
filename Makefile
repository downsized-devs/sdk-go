.PHONY: run-tests
run-tests:
	@go test -v -failfast `go list ./...` -cover

.PHONY: run-tests-report
run-tests-report:
	@go test -v -failfast `go list ./...` -cover -coverprofile=coverage.out -json > test-report.out

.PHONY: run-integ-tests
run-integ-tests:
	@go test -v -failfast `go list ./...` -cover -tags=integration

.PHONY: run-integ-tests-report
run-integ-tests-report:
	@go test -v -failfast `go list ./...` -cover  -tags=integration -coverprofile=coverage.out -json > test-report.out

.PHONY:
mock-install:
	@go install go.uber.org/mock/mockgen@v0.4.0

.PHONY: mock
mock:
	@`go env GOPATH`/bin/mockgen -source ./$(util)/$(subutil).go -destination ./tests/mock/$(util)/$(subutil).go

.PHONY: mock-all
mock-all:
	@make mock util=auth subutil=auth
	@make mock util=configbuilder subutil=configbuilder
	@make mock util=configreader subutil=configreader
	@make mock util=instrument subutil=instrument
	@make mock util=logger subutil=logger
	@make mock util=parser subutil=parser
	@make mock util=parser subutil=csv
	@make mock util=parser subutil=excel
	@make mock util=parser subutil=json
	@make mock util=query subutil=sql_builder
	@make mock util=sql subutil=sql
	@make mock util=sql subutil=sql_tx
	@make mock util=sql subutil=sql_stmt
	@make mock util=sql subutil=sql_cmd
	@make mock util=storage subutil=storage
	@make mock util=translator subutil=translator
	@make mock util=email subutil=email
	@make mock util=email subutil=email_template
	@make mock util=redis subutil=redis
	@make mock util=slack subutil=slack
	@make mock util=featureflag subutil=feature_flag
	@make mock util=ratelimiter subutil=rate_limiter
	@make mock util=timelib subutil=timelib
	@make mock util=security subutil=security
	@make mock util=tracker subutil=tracker
	@make mock util=messaging subutil=messaging
	@make mock util=pdf subutil=pdf
	@make mock util=local_storage subutil=local_storage
