// +build github

package load

import (
	"testing"

	"github.com/stretchr/testify/require"

	"copy-basta/cmd/copy-basta/clients/github"
	"copy-basta/cmd/copy-basta/common/log"
)

func Test_GithubLoader_Load(t *testing.T) {
	log.L.SetLevel(log.Debug)
	repo := "acciaioli/server-basta-template"
	ghc, err := github.NewClient(repo)
	require.Nil(t, err)
	loader, err := NewGithubLoader(ghc)
	require.Nil(t, err)

	files, err := loader.Load()
	require.Nil(t, err)

	gitIgnoreFound := false
	makefileFound := false
	mainFound := false

	for _, file := range files {
		switch file.Path {
		case "server-basta-template-master/.gitignore":
			gitIgnoreFound = true
		case "server-basta-template-master/Makefile.basta":
			makefileFound = true
		case "server-basta-template-master/main.go.basta":
			mainFound = true
		default:
			require.NotEqual(t, "", file.Path)
		}
	}

	require.True(t, gitIgnoreFound)
	require.True(t, makefileFound)
	require.True(t, mainFound)
}
