.PHONY: default clean generate build run

src = ./demo-template
dest = ./generated
input = ./basta.yaml
demoBin = demo

default: clean

clean:
	@ echo ">> cleaning"
	@ rm -rf $(dest)
	@ echo ">> done"

dynamic-generate: clean
	@ echo ">> dynamically generating new project from template"
	@ copy-basta generate --src=$(src) --dest=$(dest)
	@ echo ">> done"

generate: clean
	@ echo ">> generating new project from template"
	@ copy-basta generate --src=$(src) --dest=$(dest) --input=$(input)
	@ echo ">> done"

build: generate
	@ echo ">> building generated project"
	@ cd $(dest); go build -o $(demoBin) main.go
	@ echo ">> done"


run: build
	@ echo ">> running generated project"
	@ cd $(dest); ./$(demoBin)
	@ echo ">> done"
