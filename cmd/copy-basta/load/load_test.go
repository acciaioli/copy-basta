package load

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DiskLoader_Load(t *testing.T) {
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

	expectedFiles := []File{
		{
			Path:   "example.txt.basta",
			Mode:   0666 - 002, // default permission - umask
			Reader: strings.NewReader("Hello {{.Name}}!\nThis is an example.\n"),
		},
		{
			Path:   "nested/dummy.md",
			Mode:   0666 - 002, // default permission - umask
			Reader: strings.NewReader("# dummy\n\nThis LoadedFile is useless.\n"),
		},
	}

	loader, err := NewDiskLoader(root)
	require.Nil(t, err)
	files, err := loader.Load()
	require.Nil(t, err)

	require.Equal(t, len(expectedFiles), len(files))

	for i := range files {
		require.Equal(t, expectedFiles[i].Path, files[i].Path)
		require.Equal(t, expectedFiles[i].Mode, files[i].Mode)
		expectedR, err := ioutil.ReadAll(expectedFiles[i].Reader)
		require.Nil(t, err)
		actualR, err := ioutil.ReadAll(files[i].Reader)
		require.Nil(t, err)
		require.Equal(t, expectedR, actualR)
	}
}
