package crawl

import (
	"io"
	"os"
)

// A Crawler crawls the source project and returns all its files
type Crawler interface {
	Crawl() ([]File, error)
}

type File struct {
	Path   string
	Mode   os.FileMode
	Reader io.Reader
}
