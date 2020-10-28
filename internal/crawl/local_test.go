package crawl_test

import (
	"copy-basta/internal/crawl"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_LocalCrawler_Load(t *testing.T) {
	root := "./test-files"
	defer func() { _ = os.RemoveAll(root) }()

	err := os.MkdirAll("./test-files/nested/empty", os.ModePerm)
	require.Nil(t, err)

	mdFile, err := os.Create("./test-files/nested/dummy.md")
	require.Nil(t, err)
	_, err = mdFile.Write([]byte("# Dummy\n\nThis file is useless.\n"))
	require.Nil(t, err)

	txtFile, err := os.Create("./test-files/example.txt")
	require.Nil(t, err)
	_, err = txtFile.Write([]byte("Hello {{.Name}}!\nThis is an example.\n"))
	require.Nil(t, err)

	/*
		./test-files
		├── example.txt
		└── nested
		    ├── empty
		    └── readme.md

		2 directories, 2 files
	*/

	expectedFiles := []crawl.File{
		{
			Path:   "example.txt",
			Mode:   0666 - 002, // default permission - umask
			Reader: strings.NewReader("Hello {{.Name}}!\nThis is an example.\n"),
		},
		{
			Path:   "nested/dummy.md",
			Mode:   0666 - 002, // default permission - umask
			Reader: strings.NewReader("# Dummy\n\nThis file is useless.\n"),
		},
	}

	crawler := crawl.NewLocalCrawler(root)
	files, err := crawler.Crawl()
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
