// +build integration

package parse

import (
	"os"
	"testing"

	"copy-basta/cmd/common"

	"github.com/stretchr/testify/require"
)

func Test_Integration_Parse(t *testing.T) {
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
