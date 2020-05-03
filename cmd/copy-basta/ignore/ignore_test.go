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

func Test_newFromReader(t *testing.T) {
	expectedDirs := []string{
		"my-ignored-tree",
	}
	expectedPatterns := []string{
		"myignoredfile.go",
		"my-ignored-files/*",
		"starts*",
		"*ends",
		"*mids*",
	}
	ignorer, err := newFromReader(strings.NewReader(testLines))
	require.Nil(t, err)
	require.Equal(t, expectedDirs, ignorer.dirs)
	require.Equal(t, expectedPatterns, ignorer.patterns)
}

func Test_Ignorer_ignore(t *testing.T) {
	ignorer, err := newFromReader(strings.NewReader(testLines))
	require.Nil(t, err)

	tests := []struct {
		name     string
		filepath string
		matched  bool
	}{
		{
			name:     "LoadedFile - ignored",
			filepath: "myignoredfile.go",
			matched:  true,
		},
		{
			name:     "LoadedFile - not ignored",
			filepath: "myfile.go",
			matched:  false,
		},
		{
			name:     "tree - ignored",
			filepath: "my-ignored-tree/LoadedFile.go",
			matched:  true,
		},
		{
			name:     "tree nested - ignored",
			filepath: "my-ignored-tree/nested/LoadedFile.go",
			matched:  true,
		},
		{
			name:     "dir files - ignored",
			filepath: "my-ignored-files/LoadedFile.go",
			matched:  true,
		},
		{
			name:     "dir files nested - not ignored",
			filepath: "my-ignored-files/nested/LoadedFile.go",
			matched:  false,
		},
		{
			name:     "starts with - ignored",
			filepath: "starts-LoadedFile.go",
			matched:  true,
		},
		{
			name:     "starts with in dir - not ignored",
			filepath: "some-dir/starts-LoadedFile.go",
			matched:  false,
		},
		{
			name:     "ends with - ignored",
			filepath: "LoadedFile.go-ends",
			matched:  true,
		},
		{
			name:     "ends with in dir - not ignored",
			filepath: "some-dir/LoadedFile-ends.go",
			matched:  false,
		},
		{
			name:     "mids with - ignored",
			filepath: "LoadedFile-mids.go",
			matched:  true,
		},
		{
			name:     "mids with in dir - not ignored",
			filepath: "some-dir/LoadedFile-mids-any.go",
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
