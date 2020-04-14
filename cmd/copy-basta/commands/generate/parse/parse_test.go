package parse

import (
	"os"
	"testing"

	"github.com/spin14/copy-basta/cmd/copy-basta/common"

	"github.com/stretchr/testify/require"
)

func Test_processFile(t *testing.T) {
	root := "./test-files"
	defer func() { _ = os.RemoveAll(root) }()

	err := os.MkdirAll("./test-files/nested/empty", os.ModePerm)
	require.Nil(t, err)

	dummyMD, err := os.Create("./test-files/nested/dummy.md")
	require.Nil(t, err)
	_, err = dummyMD.Write([]byte("# dummy\n\nThis file is useless.\n"))
	require.Nil(t, err)

	exampleBasta, err := os.Create("./test-files/example.txt.basta")
	require.Nil(t, err)
	_, err = exampleBasta.Write([]byte("Hello {{.Name}}!\nThis is an example.\n"))
	require.Nil(t, err)

	/*
		./test-files
		├── example.go
		└── nested
		    ├── empty
		    └── readme.md

		2 directories, 2 files
	*/

	expectedFile := []common.File{
		{
			Path:     "example.txt",
			Mode:     0666 - 002, // default permission - umask
			Template: true,
			Content:  []byte("Hello {{.Name}}!\nThis is an example.\n"),
		},
		{
			Path:     "nested/dummy.md",
			Mode:     0666 - 002, // default permission - umask
			Template: false,
			Content:  []byte("# dummy\n\nThis file is useless.\n"),
		},
	}

	files, err := Parse(root)
	require.Nil(t, err)

	require.Equal(t, len(expectedFile), len(files))
	require.Equal(t, expectedFile[0], files[0])
	require.Equal(t, expectedFile[1], files[1])
}

func Test_validate_ok(t *testing.T) {
	files := []common.File{
		{Path: "a.go"},
		{Path: "a.md"},
		{Path: "b.txt"},
	}
	err := validate(files)
	require.Nil(t, err)
}

func Test_validate_err(t *testing.T) {
	files := []common.File{
		{Path: "a.go"},
		{Path: "a.md"},
		{Path: "a.go"},
	}
	err := validate(files)
	require.NotNil(t, err)
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

func Test_trimRootDir(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		expected string
	}{
		{name: "simple", in: "example.txt", expected: "example.txt"},
		{name: "dir", in: "dir/example.go", expected: "example.go"},
		{name: "nested", in: "nested/x/dir/example.json", expected: "x/dir/example.json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := trimRootDir(tt.in)
			require.Equal(t, out, tt.expected)
		})
	}
}
