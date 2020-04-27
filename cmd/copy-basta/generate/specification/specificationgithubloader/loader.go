package specificationgithubloader

import (
	"bytes"
	"io"

	"copy-basta/cmd/copy-basta/clients/github"
)

type Loader struct {
	specsYAML string
	ghc       *github.Client
}

func New(specsYAML string, ghc *github.Client) (*Loader, error) {
	return &Loader{specsYAML: specsYAML, ghc: ghc}, nil
}

func (l *Loader) LoadReader() (io.Reader, error) {
	b, err := l.ghc.GetContentsFileData(l.specsYAML)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}
