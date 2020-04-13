package write

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/spin14/copy-basta/cmd/copy-basta/common"
	"github.com/stretchr/testify/require"
)

func Test_write(t *testing.T) {
	root := "./test-generated"

	defer func() { _ = os.RemoveAll(root) }()

	files := []common.File{
		{
			Path:     "simple.md",
			Template: false,
			Content:  []byte("# useless readme\n"),
		},
		{
			Path:     "nested/file.txt",
			Template: false,
			Content:  []byte("this file is nested\n"),
		},
		{
			Path:     "template.go",
			Template: true,
			Content:  []byte("package generated\n\nconst Version = \"{{ .Version }}\"\n"),
		},
	}

	tVars := map[string]interface{}{"Version": "v0.1.4"}

	err := write(root, files, tVars)
	require.Nil(t, err)

	simpleMD, err := ioutil.ReadFile(path.Join(root, files[0].Path))
	require.Nil(t, err)
	require.Equal(t, simpleMD, files[0].Content)

	nested, err := ioutil.ReadFile(path.Join(root, files[1].Path))
	require.Nil(t, err)
	require.Equal(t, nested, files[1].Content)

	templateGO, err := ioutil.ReadFile(path.Join(root, files[2].Path))
	require.Nil(t, err)
	require.Equal(t, templateGO, []byte("package generated\n\nconst Version = \"v0.1.4\"\n"))
}
