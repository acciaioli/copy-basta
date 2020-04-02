.PHONY: default fix lint test

default: fix lint test

fix:
	@ echo ">> fixing"
	@ gofmt -l -w .
	@ go mod tidy
	@ echo ">> done"

lint:
	@ echo ">> running linter"
	@ golangci-lint run
	@ echo ">> done"


test:
	@ echo ">> running tests"
	@ go test `go list ./... | grep -v internal`
	@ echo ">> done"

cover:
	@ echo ">> running tests and coverage"
	@ go test --count=1 --v --cover --coverprofile=cover.out `go list ./... | grep -v internal`
	@ go tool cover --html=cover.out
	@ echo ">> done"


# CLI
cli-install:
	@ echo ">> installing cli"
	@ go install ./cmd/copy-basta-cli
	@ echo ">> done"
