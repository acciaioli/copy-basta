package parse

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"copy-basta/cmd/copy-basta/common"
	"copy-basta/cmd/copy-basta/load"
)

type ignorer interface {
	Ignore(string) bool
}

type File struct {
	Path     string
	Mode     os.FileMode
	Template bool
	Content  []byte
}

func Parse(ignorer ignorer, loadedFiles []load.File) ([]File, error) {
	err := validateFiles(loadedFiles)
	if err != nil {
		return nil, err
	}

	files, err := processFiles(ignorer, loadedFiles)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func validateFiles(files []load.File) error {
	paths := map[string]struct{}{}

	for _, file := range files {
		if _, found := paths[file.Path]; found {
			return fmt.Errorf("`%s` path found multiple times", file.Path)
		}
		paths[file.Path] = struct{}{}
	}

	return nil
}

func processFiles(ignorer ignorer, loadedFiles []load.File) ([]File, error) {
	var files []File
	for _, loadedFile := range loadedFiles {
		if ignorer.Ignore(loadedFile.Path) {
			continue
		}
		file, err := processFile(loadedFile)
		if err != nil {
			return nil, err
		}
		files = append(files, *file)
	}
	return files, nil
}

func processFile(loadedFile load.File) (*File, error) {
	content, err := ioutil.ReadAll(loadedFile.Reader)
	if err != nil {
		return nil, err
	}

	if path.Ext(loadedFile.Path) == common.TemplateExtension {
		return &File{
			Path:     trimExtension(loadedFile.Path),
			Mode:     loadedFile.Mode,
			Template: true,
			Content:  content,
		}, nil
	}

	return &File{
		Path:     loadedFile.Path,
		Mode:     loadedFile.Mode,
		Template: false,
		Content:  content,
	}, nil
}

func trimExtension(s string) string {
	return strings.TrimSuffix(s, common.TemplateExtension)
}
