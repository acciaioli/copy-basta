package common_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"copy-basta/internal/common"
)

func Test_TrimRootDir(t *testing.T) {
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
			out := common.TrimRootDir(tt.in)
			require.Equal(t, out, tt.expected)
		})
	}
}
