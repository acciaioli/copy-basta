package parse

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"copy-basta/cmd/copy-basta/common"
)

type LoadedFile struct {
	Path   string
	Mode   os.FileMode
	Reader io.Reader
}

type Loader interface {
	LoadFiles() ([]LoadedFile, error)
}

func Parse(loader Loader) ([]common.File, error) {
	loadedFiles, err := loader.LoadFiles()
	if err != nil {
		return nil, err
	}

	err = validateFiles(loadedFiles)
	if err != nil {
		return nil, err
	}

	files, err := processFiles(loadedFiles)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func validateFiles(files []LoadedFile) error {
	paths := map[string]struct{}{}

	for _, file := range files {
		if _, found := paths[file.Path]; found {
			return fmt.Errorf("`%s` path found multiple times", file.Path)
		}
		paths[file.Path] = struct{}{}
	}

	return nil
}

func processFiles(files []LoadedFile) ([]common.File, error) {
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

func processFile(file LoadedFile) (*common.File, error) {
	content, err := ioutil.ReadAll(file.Reader)
	if err != nil {
		return nil, err
	}

	if path.Ext(file.Path) == common.TemplateExtension {
		return &common.File{
			Path:     trimRootDir(trimExtension(file.Path)),
			Mode:     file.Mode,
			Template: true,
			Content:  content,
		}, nil
	}

	return &common.File{
		Path:     trimRootDir(file.Path),
		Mode:     file.Mode,
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
