package disk

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"copy-basta/cmd/copy-basta/common"
	"copy-basta/cmd/copy-basta/generate/parse"
)

func Test_Integration_Parse_Local(t *testing.T) {
	root := "./test-files"
	defer func() { _ = os.RemoveAll(root) }()

	err := os.MkdirAll("./test-files/nested/empty", os.ModePerm)
	require.Nil(t, err)

	dummyMD, err := os.Create("./test-files/nested/dummy.md")
	require.Nil(t, err)
	_, err = dummyMD.Write([]byte("# dummy\n\nThis LoadedFile is useless.\n"))
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

	expectedFiles := []common.File{
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
			Content:  []byte("# dummy\n\nThis LoadedFile is useless.\n"),
		},
	}

	loader, err := NewLoader(root)
	require.Nil(t, err)
	files, err := parse.Parse(loader)
	require.Nil(t, err)

	require.Equal(t, len(expectedFiles), len(files))
	require.Equal(t, expectedFiles[0], files[0])
	require.Equal(t, expectedFiles[1], files[1])
}
