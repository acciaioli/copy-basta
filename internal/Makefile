.PHONY: default clean generate build run

src = ./demo-template
dest = ./generated
bastaYAML = ./basta.yaml
demoBin = demo

default: clean

clean:
	@ echo ">> cleaning"
	@ rm -rf $(dest)
	@ echo ">> done"

generate: clean
	@ echo ">> generating new project from template"
	@ copy-basta generate --src=$(src) --dest=$(dest) --basta-yaml=$(bastaYAML)
	@ echo ">> done"

build: generate
	@ echo ">> building generated project"
	@ cd $(dest); go build -o $(demoBin) main.go
	@ echo ">> done"


run: build
	@ echo ">> running generated project"
	@ cd $(dest); ./$(demoBin)
	@ echo ">> done"
