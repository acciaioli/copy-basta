package generate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_processFile(t *testing.T) {
	root := "./test-files/"
	/*
		$ tree ./test-files/
		test-files
		├── dummy.md
		└── example.go.basta

		0 directories, 2 files
	*/
	expectedFile := []file{
		{
			path:     "test-files/dummy.md",
			template: false,
			content:  []byte("# dummy\n\nThis file is useless.\n"),
		},
		{
			path:     "test-files/example.go",
			template: true,
			content:  []byte("Hello {{ .Name }}!\nThis is an example.\n"),
		},
	}

	files, err := parse(root)
	require.Nil(t, err)

	require.Equal(t, len(expectedFile), len(files))
	require.Equal(t, expectedFile[0], files[0])
	require.Equal(t, expectedFile[1], files[1])
}

func Test_trimExtension(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		expected string
	}{
		{name: "simple", in: "example.basta", expected: "example"},
		{name: ".go", in: "example.go.basta", expected: "example.go"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := trimExtension(tt.in)
			require.Equal(t, out, tt.expected)
		})
	}
}
