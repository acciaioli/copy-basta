package ignore

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const testLines = `
# comment

myignoredfile.go

my-ignored-tree/

my-ignored-files/*

starts*
*ends
*mids*
`

func Test_NewIgnorer(t *testing.T) {
	expectedDirs := []string{
		"root/my-ignored-tree",
	}
	expectedPatterns := []string{
		"root/myignoredfile.go",
		"root/my-ignored-files/*",
		"root/starts*",
		"root/*ends",
		"root/*mids*",
	}
	ignorer, err := New("root", strings.NewReader(testLines))
	require.Nil(t, err)
	require.Equal(t, expectedDirs, ignorer.dirs)
	require.Equal(t, expectedPatterns, ignorer.patterns)
}

func Test_Ignorer_ignore(t *testing.T) {
	ignorer, err := New("root", strings.NewReader(testLines))
	require.Nil(t, err)

	tests := []struct {
		name     string
		filepath string
		matched  bool
	}{
		{
			name:     "LoadedFile - ignored",
			filepath: "root/myignoredfile.go",
			matched:  true,
		},
		{
			name:     "LoadedFile - not ignored",
			filepath: "root/myfile.go",
			matched:  false,
		},
		{
			name:     "tree - ignored",
			filepath: "root/my-ignored-tree/LoadedFile.go",
			matched:  true,
		},
		{
			name:     "tree nested - ignored",
			filepath: "root/my-ignored-tree/nested/LoadedFile.go",
			matched:  true,
		},
		{
			name:     "dir files - ignored",
			filepath: "root/my-ignored-files/LoadedFile.go",
			matched:  true,
		},
		{
			name:     "dir files nested - not ignored",
			filepath: "root/my-ignored-files/nested/LoadedFile.go",
			matched:  false,
		},
		{
			name:     "starts with - ignored",
			filepath: "root/starts-LoadedFile.go",
			matched:  true,
		},
		{
			name:     "starts with in dir - not ignored",
			filepath: "root/some-dir/starts-LoadedFile.go",
			matched:  false,
		},
		{
			name:     "ends with - ignored",
			filepath: "root/LoadedFile.go-ends",
			matched:  true,
		},
		{
			name:     "ends with in dir - not ignored",
			filepath: "root/some-dir/LoadedFile-ends.go",
			matched:  false,
		},
		{
			name:     "mids with - ignored",
			filepath: "root/LoadedFile-mids.go",
			matched:  true,
		},
		{
			name:     "mids with in dir - not ignored",
			filepath: "root/some-dir/LoadedFile-mids-any.go",
			matched:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched := ignorer.Ignore(tt.filepath)
			require.Equal(t, tt.matched, matched)
		})
	}
}
