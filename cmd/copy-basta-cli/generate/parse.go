package generate

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	tmplExtension = ".basta"
)

func parse(root string) ([]file, error) {
	var files []file

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		file, err := processFile(path, info)
		if err != nil {
			return err
		}
		if file != nil {
			files = append(files, *file)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	err = validate(files)

	return files, nil
}

func processFile(path string, info os.FileInfo) (*file, error) {
	if info.IsDir() {
		return nil, nil
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(path, tmplExtension) {
		return &file{path: trimExtension(path), template: true, content: content}, nil
	}

	return &file{path: path, template: false, content: content}, nil
}

func trimExtension(s string) string {
	return strings.TrimRight(s, tmplExtension)
}

func validate(files []file) error {
	panic("check all files are unique")
}
