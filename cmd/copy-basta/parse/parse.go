package parse

import (
	"fmt"
	"io/ioutil"
	"os"

	"copy-basta/cmd/copy-basta/load"
)

type ignorer interface {
	Ignore(string) bool
}

type passer interface {
	Pass(string) bool
}

type File struct {
	Path     string
	Mode     os.FileMode
	Template bool
	Content  []byte
}

func Parse(ignorer ignorer, passer passer, loadedFiles []load.File) ([]File, error) {
	err := validateFiles(loadedFiles)
	if err != nil {
		return nil, err
	}

	files, err := processFiles(ignorer, passer, loadedFiles)
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

func processFiles(ignorer ignorer, passer passer, loadedFiles []load.File) ([]File, error) {
	var files []File
	for _, loadedFile := range loadedFiles {
		if ignorer.Ignore(loadedFile.Path) {
			continue
		}
		content, err := ioutil.ReadAll(loadedFile.Reader)
		if err != nil {
			return nil, err
		}

		files = append(files, File{
			Path:     loadedFile.Path,
			Mode:     loadedFile.Mode,
			Template: !passer.Pass(loadedFile.Path),
			Content:  content,
		})
	}
	return files, nil
}
