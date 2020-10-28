// +build github

package crawl

import (
	"testing"

	"github.com/stretchr/testify/require"

	"copy-basta/internal/clients/github"
	"copy-basta/internal/common/log"
)

func Test_GithubLoader_Load(t *testing.T) {
	lvl, err := log.ToLevel(log.Debug)
	require.Nil(t, err)
	log.L.SetLevel(lvl)
	repo := "acciaioli/gorilla-mux-hello-world-basta-template"
	ghc, err := github.NewClient(repo)
	require.Nil(t, err)
	crawler := NewGithubCrawler(ghc)

	files, err := crawler.Crawl()
	require.Nil(t, err)

	expected := map[string]bool{
		".gitignore": false,
		"README.md":  false,
		"basta.yaml": false,
		"go.mod":     false,
		"main.go":    false,
	}

	for _, file := range files {
		if _, ok := expected[file.Path]; ok {
			expected[file.Path] = true
		}
	}

	for _, fileFound := range expected {
		require.True(t, fileFound)
	}
}
