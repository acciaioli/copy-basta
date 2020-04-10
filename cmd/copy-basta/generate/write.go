package generate

import (
	"os"
	"path"
	"text/template"
)

func write(root string, files []file, templateVars map[string]interface{}) error {
	for _, file := range files {
		fp, err := createFile(path.Join(root, file.path))
		if err != nil {
			return err
		}

		if file.template {
			t, err := newTemplate(file.path, string(file.content))
			if err != nil {
				return err
			}

			err = t.Execute(fp, templateVars)
			if err != nil {
				return err
			}

		} else {
			_, err := fp.Write(file.content)
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
