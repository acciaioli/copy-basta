package parse

import (
	"io"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"
)

type Ignorer struct {
	patterns []string
}

func NewIgnorer(root string, r io.Reader) (*Ignorer, error) {
	i := Ignorer{}
	if r == nil {
		return &i, nil
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(b), "\n")

	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		pattern := path.Join(root, line)
		if _, err := filepath.Match(pattern, ""); err != nil {
			return nil, err
		}

		i.patterns = append(i.patterns, pattern)
	}
	return &i, nil
}

func (i *Ignorer) ignore(s string) bool {
	for _, pattern := range i.patterns {
		matched, err := filepath.Match(pattern, s)
		if err != nil {
			log.Println("danger danger....")
			return false
		}
		if matched {
			return matched
		}
	}
	return false
}
