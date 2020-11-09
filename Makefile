.PHONY: default fix lint test install build playground-log example-init example-simple

# --- dev --- #

default: fix lint test

fix:
	@ echo ">> fixing source code"
	@ gofmt -s -l -w internal services cmd
	@ go mod tidy
	@ echo ">> done"

lint:
	@ echo ">> running linter"
	@ golangci-lint run --skip-dirs internal
	@ echo ">> done"

test:
	@ echo ">> running tests"
	@ go test --count=1 `go list ./... | grep -v examples`
	@ echo ">> done"

test-github:
	@ echo ">> running tests (+github)"
	@ go test --count=1 -v --tags=github `go list ./... | grep -v examples`
	@ echo ">> done"

test-all:
	@ echo ">> running all tests"
	@ go test --count=1 -v --tags=github `go list ./... | grep -v examples`
	@ echo ">> done"

cover:
	@ echo ">> running tests and coverage"
	@ go test --count=1 --v --cover --coverprofile=cover.out `go list ./... | grep -v examples`
	@ go tool cover --html=cover.out
	@ echo ">> done"

version=snapshot-$(USER)-$(shell git rev-parse --short HEAD)

install:
	@ echo ">> installing cli (dev)"
	@ go install -ldflags "-X main.version=$(version)" ./cmd/copy-basta
	@ echo ">> done"

# --- release --- #

build:
	@ echo ">> building cli binaries"
	@ ./build.sh
	@ echo ">> done"

# --- playground --- #

playground-log: install
	@ echo ">> running logging demo"
	@ go run internal/playground/log/main.go
	@ echo ">> done"

# --- examples --- #

tmpl = ./.tmp/
gen = ./generated
exec = $(gen)/main.sh

example-init: install
	@ echo ">> running demo init command"
	@ rm -rf $(tmpl)
	@ copy-basta init --name=$(tmpl)
	@ rm -rf $(gen)
	@ copy-basta generate --src=$(tmpl) --dest=$(gen)
	@ sh $(exec)
	@ echo ">> done"

example-simple: install
	@ echo ">> running demo generate command"
	@ cd ./examples/simple; make
	@ echo ">> done"




