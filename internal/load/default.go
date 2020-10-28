package load

import (
	"fmt"
	"io/ioutil"
	"os"

	"copy-basta/internal/crawl"
)

// A Loader takes a list of crawled Files and loads the
// necessary ones into memory
type Loader interface {
	Load([]crawl.File) ([]File, error)
}

type File struct {
	Path     string
	Mode     os.FileMode
	Template bool
	Content  []byte
}

type ignorer interface {
	Ignore(string) bool
}

type passer interface {
	Pass(string) bool
}

type loader struct {
	ignorer ignorer
	passer  passer
}

func New(ignorer ignorer, passer passer) (Loader, error) {
	if ignorer == nil {
		return nil, fmt.Errorf("ignorer can't be nil")
	}
	if passer == nil {
		return nil, fmt.Errorf("passer can't be nil")
	}

	return &loader{
		ignorer: ignorer,
		passer:  passer,
	}, nil
}

func (l *loader) Load(crawledFiles []crawl.File) ([]File, error) {
	err := validateFiles(crawledFiles)
	if err != nil {
		return nil, err
	}

	files, err := processFiles(l.ignorer, l.passer, crawledFiles)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func validateFiles(files []crawl.File) error {
	paths := map[string]struct{}{}

	for _, file := range files {
		if _, found := paths[file.Path]; found {
			return fmt.Errorf("`%s` path found multiple times", file.Path)
		}
		paths[file.Path] = struct{}{}
	}

	return nil
}

func processFiles(ignorer ignorer, passer passer, crawledFiles []crawl.File) ([]File, error) {
	var files []File
	for _, crawledFile := range crawledFiles {
		if ignorer.Ignore(crawledFile.Path) {
			continue
		}
		content, err := ioutil.ReadAll(crawledFile.Reader)
		if err != nil {
			return nil, err
		}

		files = append(files, File{
			Path:     crawledFile.Path,
			Mode:     crawledFile.Mode,
			Template: !passer.Pass(crawledFile.Path),
			Content:  content,
		})
	}
	return files, nil
}
