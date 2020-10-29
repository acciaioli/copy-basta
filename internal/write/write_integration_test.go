package write

import (
	"copy-basta/internal/load"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Integration_Write(t *testing.T) {
	root := "./test-generated"

	defer func() { _ = os.RemoveAll(root) }()

	files := []load.File{
		{
			Path:     "simple.md",
			Mode:     os.ModePerm,
			Template: false,
			Content:  []byte("# useless readme\n"),
		},
		{
			Path:     "nested/file.txt",
			Mode:     os.ModePerm,
			Template: false,
			Content:  []byte("this file is nested\n"),
		},
		{
			Path:     "template.go",
			Mode:     os.ModePerm,
			Template: true,
			Content:  []byte("package generated\n\nconst Version = \"{{ .version }}\"\n"),
		},
		{
			Path:     "string.txt",
			Mode:     os.ModePerm,
			Template: true,
			Content: []byte(`upper: {{ .myString | stringsToUpper }}
lower: {{ .myString | stringsToLower }}
title: {{ .myString | stringsTitle }}`),
		},
	}

	tVars := map[string]interface{}{"version": "v0.1.4", "myString": "my string"}

	err := Write(root, files, tVars)
	require.Nil(t, err)

	simpleMD, err := ioutil.ReadFile(filepath.Join(root, files[0].Path))
	require.Nil(t, err)
	require.Equal(t, simpleMD, files[0].Content)

	nested, err := ioutil.ReadFile(filepath.Join(root, files[1].Path))
	require.Nil(t, err)
	require.Equal(t, nested, files[1].Content)

	templateGO, err := ioutil.ReadFile(filepath.Join(root, files[2].Path))
	require.Nil(t, err)
	require.Equal(t, templateGO, []byte("package generated\n\nconst Version = \"v0.1.4\"\n"))

	templateTXT, err := ioutil.ReadFile(filepath.Join(root, files[3].Path))
	require.Nil(t, err)
	require.Equal(t, templateTXT, []byte("upper: MY STRING\nlower: my string\ntitle: My String"))
}
