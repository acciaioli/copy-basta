package disk

import (
	"copy-basta/cmd/common"
	"copy-basta/cmd/common/log"
	"copy-basta/cmd/generate/parse"
	"copy-basta/cmd/generate/parse/ignore"
	"fmt"
	"os"
	"path/filepath"
)

type Loader struct {
	root string
}

func NewLoader(root string) (*Loader, error) {
	return &Loader{root: root}, nil
}

func (l *Loader) LoadFiles() ([]parse.LoadedFile, error) {
	ignorer, err := l.getIgnorer()
	if err != nil {
		return nil, err
	}

	files, err := l.walkFiles(ignorer)
	if err != nil {
		return nil, err
	}

	return files, err
}

func (l *Loader) getIgnorer() (*ignore.Ignorer, error) {
	ignoreFilePath := filepath.Join(l.root, common.IgnoreFile)

	if fInfo, err := os.Stat(ignoreFilePath); err == nil {
		if fInfo.IsDir() {
			return nil, fmt.Errorf("%s must no be a dir", ignoreFilePath)
		}
		file, err := os.Open(ignoreFilePath)
		if err != nil {
			return nil, err
		}
		return ignore.New(l.root, file)
	}

	return ignore.New(l.root, nil)
}

func (l *Loader) walkFiles(ignorer *ignore.Ignorer) ([]parse.LoadedFile, error) {
	var files []parse.LoadedFile

	err := filepath.Walk(l.root, func(fPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if ignorer.Ignore(fPath) {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		r, err := os.Open(fPath)
		if err != nil {
			log.L.DebugWithData("external error", log.Data{"error": err.Error()})
			return err
		}
		files = append(files, parse.LoadedFile{Path: fPath, Mode: info.Mode(), Reader: r})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
