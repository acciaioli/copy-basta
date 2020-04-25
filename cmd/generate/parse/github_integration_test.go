// +build integration
// +build github

package parse

import (
	"copy-basta/cmd/common/log"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GitHubParser_Parse(t *testing.T) {
	log.L.SetLevel(log.Debug)

	repoRef := "github.com/acciaioli/copy-basta"
	p, err := NewGitHubParser(repoRef)
	require.Nil(t, err)
	files, err := p.Parse()
	require.Nil(t, err)
	for i, f := range files {
		require.NotEqual(t, "", f.path)
		require.NotEqual(t, uint32(0), f.mode)
		require.NotNil(t, f.r)
		log.L.DebugWithData(fmt.Sprintf("parsed file %d", i), log.Data{"path": f.path})
	}
}
