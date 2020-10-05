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

type testIgnorer struct{}

func (i *testIgnorer) Ignore(s string) bool {
	return strings.Contains(s, "ingore")
}

type testPasser struct{}

func (i *testPasser) Pass(s string) bool {
	return strings.Contains(s, "pass")
}

func Test_processFiles(t *testing.T) {
	loadedFiles := []load.File{
		{
			Path:   "ignore.go",
			Mode:   0123,
			Reader: strings.NewReader("ignore.go"),
		},
		{
			Path:   "pass.go",
			Mode:   0123,
			Reader: strings.NewReader("pass.go"),
		},
		{
			Path:   "template.cpp",
			Mode:   0123,
			Reader: strings.NewReader("template.cpp"),
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

	files, err := processFiles(&testIgnorer{}, &testPasser{}, loadedFiles)
	require.Nil(t, err)
	require.Equal(t, expectedFiles, files)
}