package crawl

import (
	"os"
	"path/filepath"

	"copy-basta/internal/common"
)

type localCrawler struct {
	root string
}

func NewLocalCrawler(root string) Crawler {
	return &localCrawler{root: root}
}

func (c *localCrawler) Crawl() ([]File, error) {
	var files []File

	err := filepath.Walk(c.root, func(fPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		r, err := os.Open(fPath)
		if err != nil {
			return err
		}
		files = append(files, File{Path: common.TrimRootDir(fPath), Mode: info.Mode(), Reader: r})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
