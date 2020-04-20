package write

import (
	"testing"

	"github.com/stretchr/testify/require"

	"copy-basta/cmd/common"
)

func Test_generateFromTemplate(t *testing.T) {
	input := common.InputVariables{
		"user":   "vasco",
		"number": 22,
	}

	tests := []struct {
		name               string
		rawPath            string
		rawContent         string
		expectedGenPath    string
		expectedGenContent string
	}{
		{
			name:               "template in path",
			rawPath:            "dir/{{.user}}.go",
			rawContent:         "package dir\n\nconst FavoriteNumber=19",
			expectedGenPath:    "dir/vasco.go",
			expectedGenContent: "package dir\n\nconst FavoriteNumber=19",
		},
		{
			name:               "template in content",
			rawPath:            "dir/maria.go",
			rawContent:         "package dir\n\nconst FavoriteNumber={{.number}}",
			expectedGenPath:    "dir/maria.go",
			expectedGenContent: "package dir\n\nconst FavoriteNumber=22",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			genPath, genContent, err := generateFromTemplate(tt.rawPath, tt.rawContent, input)
			require.Nil(t, err)
			require.Equal(t, tt.expectedGenPath, *genPath)
			require.Equal(t, tt.expectedGenContent, *genContent)
		})
	}
}

func Test_generateFromTemplate_error(t *testing.T) {
	input := common.InputVariables{
		"user":   "vasco",
		"number": 22,
	}

	tests := []struct {
		name       string
		rawPath    string
		rawContent string
	}{
		{
			name:       "template in path",
			rawPath:    "dir/{{.username}}.go",
			rawContent: "package dir\n\nconst FavoriteNumber=19",
		},
		{
			name:       "template in content",
			rawPath:    "dir/maria.go",
			rawContent: "package dir\n\nconst FavoriteNumber={{.favoritenumber}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := generateFromTemplate(tt.rawPath, tt.rawContent, input)
			require.NotNil(t, err)
		})
	}
}
