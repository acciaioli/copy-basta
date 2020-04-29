package parse

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"copy-basta/cmd/copy-basta/load"
)

func Test_validate_ok(t *testing.T) {
	files := []load.File{
		{Path: "a.go"},
		{Path: "a.md"},
		{Path: "b.txt"},
	}
	err := validateFiles(files)
	require.Nil(t, err)
}

func Test_validate_err(t *testing.T) {
	files := []load.File{
		{Path: "a.go"},
		{Path: "a.md"},
		{Path: "a.go"},
	}
	err := validateFiles(files)
	require.NotNil(t, err)
}

type testIgnorer struct {
	paths []string
}

func (i *testIgnorer) Ignore(s string) bool {
	for _, p := range i.paths {
		if s == p {
			return true
		}
	}
	return false
}

func Test_processFiles(t *testing.T) {
	i := testIgnorer{paths: []string{"a.go", "d.python"}}
	loadedFiles := []load.File{
		{
			Path:   "a.go",
			Mode:   0123,
			Reader: strings.NewReader("a.go"),
		},
		{
			Path:   "b.go",
			Mode:   0123,
			Reader: strings.NewReader("b.go"),
		},
		{
			Path:   "c.go.basta",
			Mode:   0123,
			Reader: strings.NewReader("c.go.basta"),
		},
	}

	expectedFiles := []File{
		{
			Path:     "b.go",
			Mode:     0123,
			Content:  []byte("b.go"),
			Template: false,
		},
		{
			Path:     "c.go",
			Mode:     0123,
			Content:  []byte("c.go.basta"),
			Template: true,
		},
	}

	files, err := processFiles(&i, loadedFiles)
	require.Nil(t, err)
	require.Equal(t, expectedFiles, files)
}

func Test_trimExtension(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		expected string
	}{
		{name: "simple", in: "example.txt", expected: "example.txt"},
		{name: ".txt", in: "example.txt.basta", expected: "example.txt"},
		{name: ".go", in: "example.go.basta", expected: "example.go"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := trimExtension(tt.in)
			require.Equal(t, out, tt.expected)
		})
	}
}
