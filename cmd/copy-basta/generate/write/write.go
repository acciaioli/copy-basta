package write

import (
	"log"
	"os"
	"path"
	"text/template"

	"github.com/spin14/copy-basta/cmd/copy-basta/generate/common"
)

func Write(destDir string, files []common.File, input common.InputVariables) error {
	err := write(destDir, files, input)
	if err != nil {
		cleanup(destDir)
	}
	return err
}

func write(destDir string, files []common.File, input common.InputVariables) error {
	for _, file := range files {
		fp, err := createFile(path.Join(destDir, file.Path))
		if err != nil {
			return err
		}

		if file.Template {
			t, err := newTemplate(file.Path, string(file.Content))
			if err != nil {
				return err
			}

			err = t.Execute(fp, input)
			if err != nil {
				return err
			}

		} else {
			_, err := fp.Write(file.Content)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func cleanup(destDir string) {
	if err := os.RemoveAll(destDir); err != nil {
		log.Print("[ERROR] cleanup fail")
	}
}

func newTemplate(name string, t string) (*template.Template, error) {
	return template.New(name).Option("missingkey=error").Parse(t)
}

func createFile(filepath string) (*os.File, error) {
	err := os.MkdirAll(path.Dir(filepath), os.ModePerm)
	if err != nil {
		return nil, err
	}
	return os.Create(filepath)
}
