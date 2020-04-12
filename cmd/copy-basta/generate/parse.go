package generate

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spin14/copy-basta/cmd/copy-basta/generate/common"
)

const (
	tmplExtension = ".basta"
)

func parse(root string) ([]common.File, error) {
	var files []common.File

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

func processFile(filepath string, info os.FileInfo) (*common.File, error) {
	if info.IsDir() {
		return nil, nil
	}

	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	if path.Ext(filepath) == tmplExtension {
		return &common.File{Path: trimRootDir(trimExtension(filepath)), Template: true, Content: content}, nil
	}

	return &common.File{Path: trimRootDir(filepath), Template: false, Content: content}, nil
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

func validate(files []common.File) error {
	paths := map[string]struct{}{}

	for _, file := range files {
		if _, found := paths[file.Path]; found {
			return fmt.Errorf("`%s` path found multiple times", file.Path)
		}
		paths[file.Path] = struct{}{}
	}

	return nil
}
