.PHONY: default clean generate build run

templateDir = ./demo-template
outputDir = ./generated
bastaYAML = ./basta.yaml
demoBin = demo

default: clean

clean:
	@ echo ">> cleaning"
	@ rm -rf $(outputDir)
	@ echo ">> done"

generate: clean
	@ echo ">> generating new project from template"
	@ copy-basta generate --template-root=$(templateDir) --project-name=$(outputDir) --basta-yaml=$(bastaYAML)
	@ echo ">> done"

build: generate
	@ echo ">> building generated project"
	@ cd $(outputDir); go build -o $(demoBin) main.go
	@ echo ">> done"


run: build
	@ echo ">> running generated project"
	@ cd $(outputDir); ./$(demoBin)
	@ echo ">> done"
