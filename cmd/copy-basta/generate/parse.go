package generate

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
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

func processFile(filepath string, info os.FileInfo) (*file, error) {
	if info.IsDir() {
		return nil, nil
	}

	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	if path.Ext(filepath) == tmplExtension {
		return &file{path: trimRootDir(trimExtension(filepath)), template: true, content: content}, nil
	}

	return &file{path: trimRootDir(filepath), template: false, content: content}, nil
}

func trimExtension(s string) string {
	return strings.TrimSuffix(s, tmplExtension)
}

func trimRootDir(s string) string {
	ss := strings.Split(s, "/")
	if len(ss) == 1 {
		return ss[0]
	}
	return strings.Join(ss[1:], "/")
}

func validate(files []file) error {
	paths := map[string]struct{}{}

	for _, file := range files {
		if _, found := paths[file.path]; found {
			return fmt.Errorf("`%s` path found multiple times", file.path)
		}
		paths[file.path] = struct{}{}
	}

	return nil
}
