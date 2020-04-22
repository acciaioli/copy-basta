package bootstrap

import (
	"fmt"
	"os"
	"path/filepath"

	"copy-basta/cmd/common/log"

	"copy-basta/cmd/common"
)

func Bootstrap(destDir string) error {
	err := bootstrap(destDir)
	if err != nil {
		cleanup(destDir)
	}
	return err
}

func bootstrap(destDir string) error {
	err := os.Mkdir(destDir, os.ModePerm)
	if err != nil {
		return err
	}

	_, err = bootstrapFile(destDir, common.IgnoreFile, ignoreText)
	if err != nil {
		return err
	}

	_, err = bootstrapFile(destDir, readmeFileName, readmeText)
	if err != nil {
		return err
	}

	_, err = bootstrapFile(destDir, common.SpecFile, specText)
	if err != nil {
		return err
	}

	scriptFile, err := bootstrapFile(destDir, fmt.Sprintf("%s%s", scriptFileName, common.TemplateExtension), scriptText)
	if err != nil {
		return err
	}
	err = scriptFile.Chmod(scriptFileChmodCode)
	if err != nil {
		return err
	}

	return nil
}

func bootstrapFile(destDir, filName, fileText string) (*os.File, error) {
	p := filepath.Join(destDir, filName)
	f, err := os.Create(p)
	if err != nil {
		return nil, err
	}
	_, err = f.WriteString(fileText)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func cleanup(destDir string) {
	if err := os.RemoveAll(destDir); err != nil {
		log.L.DebugWithData("external error", log.Data{"error": err.Error()})
	}
}

const (
	ignoreText = `
# ignored dirs
.git/

# ignored patterns
readme.md
.bastaignore
spec.yaml
`
	readmeFileName = "readme.md"
	readmeText     = `# template

This is the readme of the template. 

It will not be copied because it is featured in the .bastaignore file.

To generate a project from this template you should run:
*copy-basta generate --src=template-dir --src=new-project*

*--src* should be the directory containing this file

You should override this file with information that is relevant for your template!
`
	specText = `---
variables:
  - name: myName
    type: string
    description: your name so that you can be greeted
  - name: greet
    type: string
    description: your favorite greet expression
    default: hello
`
	scriptFileName      = "main.sh"
	scriptFileChmodCode = 0777
	scriptText          = `#!/bin/sh
# Your generated code bellow
echo {{.greet}} {{.myName}}!
`
)
