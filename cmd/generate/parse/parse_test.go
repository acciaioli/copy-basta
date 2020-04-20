package parse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_validate_ok(t *testing.T) {
	files := []file{
		{path: "a.go"},
		{path: "a.md"},
		{path: "b.txt"},
	}
	err := validateFiles(files)
	require.Nil(t, err)
}

func Test_validate_err(t *testing.T) {
	files := []file{
		{path: "a.go"},
		{path: "a.md"},
		{path: "a.go"},
	}
	err := validateFiles(files)
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
