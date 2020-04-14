package parse

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spin14/copy-basta/cmd/copy-basta/common"
)

func Parse(root string) ([]common.File, error) {
	var files []common.File

	ignorer, err := getIgnorer(root)
	if err != nil {
		return nil, err
	}

	if err := filepath.Walk(root, func(fPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if ignorer.ignore(fPath) {
			return nil
		}

		file, err := processFile(fPath, info)
		if err != nil {
			return err
		}
		if file != nil {
			files = append(files, *file)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	if err := validate(files); err != nil {
		return nil, err
	}

	return files, nil
}

func getIgnorer(root string) (*Ignorer, error) {
	ignoreFilePath := filepath.Join(root, common.IgnoreFile)

	if fInfo, err := os.Stat(ignoreFilePath); err == nil {
		if fInfo.IsDir() {
			return nil, fmt.Errorf("%s must no be a dir", ignoreFilePath)
		}
		file, err := os.Open(ignoreFilePath)
		if err != nil {
			return nil, err
		}
		return NewIgnorer(root, file)
	}

	return NewIgnorer(root, nil)
}

func processFile(filepath string, info os.FileInfo) (*common.File, error) {
	if info.IsDir() {
		return nil, nil
	}

	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	if path.Ext(filepath) == common.TemplateExtension {
		return &common.File{Path: trimRootDir(trimExtension(filepath)), Mode: info.Mode(), Template: true, Content: content}, nil
	}

	return &common.File{Path: trimRootDir(filepath), Mode: info.Mode(), Template: false, Content: content}, nil
}

func trimExtension(s string) string {
	return strings.TrimSuffix(s, common.TemplateExtension)
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
