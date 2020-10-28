package write

import (
	"copy-basta/internal/common/errors"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"copy-basta/internal/common"
	"copy-basta/internal/common/log"
	"copy-basta/internal/load"
)

// A Writer is responsible for executing the templates
// and writing the the loaded files to their destination
type Writer interface {
	Write(files []load.File, input common.InputVariables) error
}

type diskWriter struct {
	errorBuilder errors.ErrorBuilder
	destDir      string
}

func NewDiskWriter(destDir string) Writer {
	return &diskWriter{
		destDir: destDir,
	}
}

func (d diskWriter) Write(files []load.File, input common.InputVariables) error {
	err := write(d.destDir, files, input)
	if err != nil {
		cleanup(d.destDir)
		return d.errorBuilder.NewInternalError(err)
	}
	return nil
}

func Write(destDir string, files []load.File, input common.InputVariables) error {
	err := write(destDir, files, input)
	if err != nil {
		cleanup(destDir)
	}
	return err
}

func write(destDir string, files []load.File, input common.InputVariables) error {
	for _, file := range files {
		fpath := filepath.Join(destDir, file.Path)

		if !file.Template {
			err := writeFile(fpath, file.Mode, file.Content)
			if err != nil {
				return err
			}
			continue
		}

		genPath, genContent, err := generateFromTemplate(fpath, string(file.Content), input)
		if err != nil {
			return err
		}
		err = writeFile(*genPath, file.Mode, []byte(*genContent))
		if err != nil {
			return err
		}
	}
	return nil
}

func cleanup(destDir string) {
	if err := os.RemoveAll(destDir); err != nil {
		log.L.DebugWithData("external error", log.Data{"error": err.Error()})
	}
}

func generateFromTemplate(rawPath string, rawContent string, input common.InputVariables) (*string, *string, error) {
	w := strings.Builder{}
	pathT, err := newTemplate("pathTemplate").Parse(rawPath)
	if err != nil {
		return nil, nil, err
	}
	err = pathT.Execute(&w, input)
	if err != nil {
		return nil, nil, err
	}
	generatedPath := w.String()

	w.Reset()
	contentT, err := newTemplate("contentTemplate").Parse(rawContent)
	if err != nil {
		return nil, nil, err
	}
	err = contentT.Execute(&w, input)
	if err != nil {
		return nil, nil, err
	}
	generatedContent := w.String()

	return &generatedPath, &generatedContent, nil
}

func writeFile(fpath string, mode os.FileMode, content []byte) error {
	err := os.MkdirAll(path.Dir(fpath), os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(fpath)
	if err != nil {
		return err
	}
	err = f.Chmod(mode)
	if err != nil {
		return err
	}
	_, err = f.Write(content)
	if err != nil {
		return err
	}
	return nil
}

func newTemplate(name string) *template.Template {
	return template.New("t").
		Option("missingkey=error").
		Funcs(template.FuncMap{
			"stringsToUpper": strings.ToUpper,
			"stringsToLower": strings.ToLower,
			"stringsTitle":   strings.Title,
		})
}
