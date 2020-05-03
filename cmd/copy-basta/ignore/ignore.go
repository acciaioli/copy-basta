package ignore

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"copy-basta/cmd/copy-basta/load"
)

type Ignorer struct {
	dirs     []string
	patterns []string
}

func New(ignoreFileName string, files []load.File) (*Ignorer, error) {
	var ignoreFile *load.File
	for _, f := range files {
		if f.Path == ignoreFileName {
			ignoreFile = &f
			break
		}
	}
	if ignoreFile == nil {
		return nil, fmt.Errorf("specification: failed to find spec file (%s)", ignoreFileName)
	}

	return newFromReader(ignoreFile.Reader)
}

func newFromReader(r io.Reader) (*Ignorer, error) {
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

		if strings.HasSuffix(line, "/") {
			// completely excluded dir
			i.dirs = append(i.dirs, strings.TrimSuffix(line, "/"))
		} else {
			// patterns
			if _, err := filepath.Match(line, ""); err != nil {
				return nil, err
			}

			i.patterns = append(i.patterns, line)
		}
	}
	return &i, nil
}

func (i *Ignorer) Ignore(s string) bool {
	for _, dir := range i.dirs {
		target := s
		for {
			target = filepath.Dir(target)
			if target == "." {
				break
			}
			if target == dir {
				return true
			}
		}
	}

	for _, pattern := range i.patterns {
		matched, err := filepath.Match(pattern, s)
		if err != nil {
			return false
		}
		if matched {
			return matched
		}
	}

	return false
}
