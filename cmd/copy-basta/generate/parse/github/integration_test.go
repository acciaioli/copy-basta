// +build github

package github

import (
	"copy-basta/cmd/copy-basta/common"
	"copy-basta/cmd/copy-basta/generate/parse"
	"testing"

	"github.com/stretchr/testify/require"

	"copy-basta/cmd/copy-basta/common/log"
)

func Test_Integration_Parse_Github(t *testing.T) {
	log.L.SetLevel(log.Debug)
	repo := "acciaioli/server-basta-template"
	loader, err := NewLoader(repo)
	require.Nil(t, err)

	files, err := parse.Parse(loader)
	require.Nil(t, err)

	var gitIgnoreFile common.File
	var makefileFile common.File
	var mainFile common.File

	for _, file := range files {
		switch file.Path {
		case ".gitignore":
			gitIgnoreFile = file
		case "Makefile":
			makefileFile = file
		case "main.go":
			mainFile = file
		default:
			require.NotEqual(t, "", file.Path)
		}
	}
	require.NotNil(t, gitIgnoreFile)
	require.False(t, gitIgnoreFile.Template)
	require.NotNil(t, makefileFile)
	require.True(t, makefileFile.Template)
	require.NotNil(t, mainFile)
	require.True(t, mainFile.Template)
}
