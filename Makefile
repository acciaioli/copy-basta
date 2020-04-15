.PHONY: default fix lint test install demo

default: fix lint test

fix:
	@ echo ">> fixing"
	@ gofmt -l -w .
	@ go mod tidy
	@ echo ">> done"

lint:
	@ echo ">> running linter"
	@ golangci-lint run --skip-dirs internal
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

install:
	@ echo ">> installing cli"
	@ go install ./cmd/copy-basta
	@ echo ">> done"

demo-generate: install
	@ echo ">> running demo generate command"
	@ cd ./internal; make run
	@ echo ">> done"


tmpl = ./.tmp/
gen = ./generated
exec = $(gen)/main.sh

demo-init: install
	@ echo ">> running demo init command"
	@ rm -rf $(tmpl)
	@ copy-basta init --name=$(tmpl)
	@ rm -rf $(gen)
	@ copy-basta generate --src=$(tmpl) --dest=$(gen)
	@ sh $(exec)
	@ echo ">> done"

demo-logs: install
	@ echo ">> running logging demo"
	@ cd internal/demo-log; go run main.go
	@ echo ">> done"
