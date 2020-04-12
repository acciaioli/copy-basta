package generate

import (
	"os"
	"path"
	"text/template"

	"github.com/spin14/copy-basta/cmd/copy-basta/generate/common"
)

func write(root string, files []common.File, templateVars map[string]interface{}) error {
	for _, file := range files {
		fp, err := createFile(path.Join(root, file.Path))
		if err != nil {
			return err
		}

		if file.Template {
			t, err := newTemplate(file.Path, string(file.Content))
			if err != nil {
				return err
			}

			err = t.Execute(fp, templateVars)
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
