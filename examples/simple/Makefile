.PHONY: default clean generate build run

src = ./template
dest = ./generated
input = ./input.yaml
exec = simple

default: run

clean:
	@ echo ">> cleaning"
	@ rm -rf $(dest)
	@ echo ">> done"

simple-dynamic-generate: clean
	@ echo ">> dynamically generating new project from template"
	@ copy-basta generate --src=$(src) --dest=$(dest) --log-level=debug
	@ echo ">> done"

generate: clean
	@ echo ">> generating new project from template"
	@ copy-basta generate --src=$(src) --dest=$(dest) --input=$(input) --log-level=debug
	@ echo ">> done"

build: generate
	@ echo ">> building generated project"
	@ cd $(dest); go build -o $(exec) main.go
	@ echo ">> done"

run: build
	@ echo ">> running generated project"
	@ cd $(dest); ./$(exec)
	@ echo ">> done"
