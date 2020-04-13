package parse

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const testLines = `
# comment

myignoredfile.go

my-ignored-dir/*

starts*
*ends
*mids*
`

func Test_NewIgnorer(t *testing.T) {
	expectedPatterns := []string{
		"root/myignoredfile.go",
		"root/my-ignored-dir/*",
		"root/starts*",
		"root/*ends",
		"root/*mids*",
	}
	ignorer, err := NewIgnorer("root", strings.NewReader(testLines))
	require.Nil(t, err)
	require.Equal(t, expectedPatterns, ignorer.patterns)
}

func Test_Ignorer_ignore(t *testing.T) {
	ignorer, err := NewIgnorer("root", strings.NewReader(testLines))
	require.Nil(t, err)

	tests := []struct {
		name     string
		filepath string
		matched  bool
	}{
		{
			name:     "file matched",
			filepath: "root/myignoredfile.go",
			matched:  true,
		},
		{
			name:     "file not matched",
			filepath: "root/myfile.go",
			matched:  false,
		},
		{
			name:     "file matched in dir",
			filepath: "root/my-ignored-dir/file.go",
			matched:  true,
		},
		{
			name:     "file not matched in dir nested dir",
			filepath: "root/my-ignored-dir/nested/file.go",
			matched:  false,
		},
		{
			name:     "file matched starts with",
			filepath: "root/starts-file.go",
			matched:  true,
		},
		{
			name:     "file not matched starts with in dir",
			filepath: "root/some-dir/starts-file.go",
			matched:  false,
		},
		{
			name:     "file matched ends with",
			filepath: "root/file.go-ends",
			matched:  true,
		},
		{
			name:     "file not matched ends with in dir",
			filepath: "root/some-dir/file-ends.go",
			matched:  false,
		},
		{
			name:     "file matched mids with",
			filepath: "root/file-mids.go",
			matched:  true,
		},
		{
			name:     "file not matched mids with in dir",
			filepath: "root/some-dir/file-mids-any.go",
			matched:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched := ignorer.ignore(tt.filepath)
			require.Equal(t, tt.matched, matched)
		})
	}
}
