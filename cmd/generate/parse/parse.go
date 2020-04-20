package parse

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"copy-basta/cmd/common"
)

type file struct {
	path string
	mode os.FileMode
	r    io.Reader
}

func Parse(root string) ([]common.File, error) {
	ignorer, err := getIgnorer(root)
	if err != nil {
		return nil, err
	}

	files, err := getFiles(root, ignorer)
	if err != nil {
		return nil, err
	}

	err = validateFiles(files)
	if err != nil {
		return nil, err
	}

	cFiles, err := processFiles(files)
	if err != nil {
		return nil, err
	}

	return cFiles, nil
}

func getIgnorer(root string) (*ignorer, error) {
	ignoreFilePath := filepath.Join(root, common.IgnoreFile)

	if fInfo, err := os.Stat(ignoreFilePath); err == nil {
		if fInfo.IsDir() {
			return nil, fmt.Errorf("%s must no be a dir", ignoreFilePath)
		}
		file, err := os.Open(ignoreFilePath)
		if err != nil {
			return nil, err
		}
		return newIgnorer(root, file)
	}

	return newIgnorer(root, nil)
}

func getFiles(root string, ignorer *ignorer) ([]file, error) {
	var files []file

	err := filepath.Walk(root, func(fPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if ignorer.ignore(fPath) {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		r, err := os.Open(fPath)
		if err != nil {
			return err
		}
		files = append(files, file{path: fPath, mode: info.Mode(), r: r})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func validateFiles(files []file) error {
	paths := map[string]struct{}{}

	for _, file := range files {
		if _, found := paths[file.path]; found {
			return fmt.Errorf("`%s` path found multiple times", file.path)
		}
		paths[file.path] = struct{}{}
	}

	return nil
}

func processFiles(files []file) ([]common.File, error) {
	var cFiles []common.File
	for _, file := range files {
		cFile, err := processFile(file)
		if err != nil {
			return nil, err
		}
		cFiles = append(cFiles, *cFile)
	}
	return cFiles, nil
}

func processFile(file file) (*common.File, error) {
	content, err := ioutil.ReadAll(file.r)
	if err != nil {
		return nil, err
	}

	if path.Ext(file.path) == common.TemplateExtension {
		return &common.File{
			Path:     trimRootDir(trimExtension(file.path)),
			Mode:     file.mode,
			Template: true,
			Content:  content,
		}, nil
	}

	return &common.File{
		Path:     trimRootDir(file.path),
		Mode:     file.mode,
		Template: false,
		Content:  content,
	}, nil
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
