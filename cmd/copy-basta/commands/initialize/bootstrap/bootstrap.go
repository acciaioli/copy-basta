package bootstrap

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spin14/copy-basta/cmd/copy-basta/common"
)

func Bootstrap(destDir string) error {
	err := bootstrap(destDir)
	if err != nil {
		cleanup(destDir)
	}
	return err
}

func bootstrap(destDir string) error {
	if err := os.Mkdir(destDir, os.ModePerm); err != nil {
		return err
	}

	ignorePath := filepath.Join(destDir, common.IgnoreFile)
	if f, err := os.Create(ignorePath); err != nil {
		return err
	} else {
		if _, err := f.WriteString(ignoreText); err != nil {
			return err
		}
	}

	specPath := filepath.Join(destDir, common.SpecFile)
	if f, err := os.Create(specPath); err != nil {
		return err
	} else {
		if _, err := f.WriteString(specText); err != nil {
			return err
		}
	}

	scriptPath := filepath.Join(destDir, fmt.Sprintf("%s%s", scriptFileName, common.TemplateExtension))
	if f, err := os.Create(scriptPath); err != nil {
		return err
	} else {
		if _, err := f.WriteString(scriptText); err != nil {
			return err
		}
		if err := f.Chmod(scriptFileChmodCode); err != nil {
			return err
		}
	}

	return nil
}

func cleanup(destDir string) {
	if err := os.RemoveAll(destDir); err != nil {
	}
}

const (
	ignoreText = `
# ignored dirs
.git/

# ignored patterns
ignore-me.md
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
