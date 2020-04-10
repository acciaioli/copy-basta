package generate

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_write(t *testing.T) {
	root := "./test-files/generated"

	defer func() { _ = os.RemoveAll(root) }()

	files := []file{
		{
			path:     "simple.md",
			template: false,
			content:  []byte("# useless readme\n"),
		},
		{
			path:     "nested/file.txt",
			template: false,
			content:  []byte("this file is nested\n"),
		},
		{
			path:     "template.go",
			template: true,
			content:  []byte("package generated\n\nconst Version = \"{{ .Version }}\"\n"),
		},
	}

	tVars := map[string]interface{}{"Version": "v0.1.4"}

	err := write(root, files, tVars)
	require.Nil(t, err)

	simpleMD, err := ioutil.ReadFile(path.Join(root, files[0].path))
	require.Nil(t, err)
	require.Equal(t, simpleMD, files[0].content)

	nested, err := ioutil.ReadFile(path.Join(root, files[1].path))
	require.Nil(t, err)
	require.Equal(t, nested, files[1].content)

	templateGO, err := ioutil.ReadFile(path.Join(root, files[2].path))
	require.Nil(t, err)
	require.Equal(t, templateGO, []byte("package generated\n\nconst Version = \"v0.1.4\"\n"))
}
