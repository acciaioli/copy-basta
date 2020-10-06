package specification

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var testPatterns = []string{
	"my-file.go",
	"my-tree/",
	"my-files/*",
	"starts*",
	"*ends",
	"*contains*",
}

func Test_Ignorer_ignore(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		matched  bool
	}{
		{
			name:     "LoadedFile - matched",
			filepath: "my-file.go",
			matched:  true,
		},
		{
			name:     "LoadedFile - not matched",
			filepath: "myfile.go",
			matched:  false,
		},
		{
			name:     "tree - matched",
			filepath: "my-tree/LoadedFile.go",
			matched:  true,
		},
		{
			name:     "tree nested - matched",
			filepath: "my-tree/nested/LoadedFile.go",
			matched:  true,
		},
		{
			name:     "dir files - matched",
			filepath: "my-files/LoadedFile.go",
			matched:  true,
		},
		{
			name:     "dir files nested - not matched",
			filepath: "my-files/nested/LoadedFile.go",
			matched:  false,
		},
		{
			name:     "starts with - matched",
			filepath: "starts-LoadedFile.go",
			matched:  true,
		},
		{
			name:     "starts with in dir - not matched",
			filepath: "some-dir/starts-LoadedFile.go",
			matched:  false,
		},
		{
			name:     "ends with - matched",
			filepath: "LoadedFile.go-ends",
			matched:  true,
		},
		{
			name:     "ends with in dir - not matched",
			filepath: "some-dir/LoadedFile-ends.go",
			matched:  false,
		},
		{
			name:     "contains with - matched",
			filepath: "LoadedFile-contains.go",
			matched:  true,
		},
		{
			name:     "contains with in dir - not matched",
			filepath: "some-dir/LoadedFile-contains-any.go",
			matched:  false,
		},
	}

	ignorer, err := NewPatternMatcher(testPatterns)
	require.Nil(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched := ignorer.Match(tt.filepath)
			require.Equal(t, tt.matched, matched)
		})
	}
}
