package write

import (
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/spin14/copy-basta/cmd/copy-basta/common/log"

	"github.com/spin14/copy-basta/cmd/copy-basta/common"
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

		fp, err := createFile(filepath.Join(destDir, file.Path), file.Mode)
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
		log.TheLogger.DebugWithData("external error", log.LoggerData{"error": err.Error()})
	}
}

func newTemplate(name string, t string) (*template.Template, error) {
	return template.New(name).Option("missingkey=error").Parse(t)
}

func createFile(filepath string, mode os.FileMode) (*os.File, error) {
	err := os.MkdirAll(path.Dir(filepath), os.ModePerm)
	if err != nil {
		return nil, err
	}
	f, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	err = f.Chmod(mode)
	if err != nil {
		return nil, err
	}
	return f, nil
}
