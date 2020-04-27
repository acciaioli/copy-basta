package specificationdiskloader

import (
	"io"
	"os"
)

type Loader struct {
	specsYAML string
}

func New(specsYAML string) (*Loader, error) {
	return &Loader{specsYAML: specsYAML}, nil
}

func (l *Loader) LoadReader() (io.Reader, error) {
	return os.Open(l.specsYAML)
}
