package specification

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_newFromReader(t *testing.T) {
	yml := `---
ignore:
  - myfileA.py
  - myDirB/
pass-through:
  - myFileC.cpp
  - myExpressionD*
on-overwrite:
  exclude:
   - manuallyUpdated.txt
variables:
  - name: stringHello
    type: string
    description: used to greet
    default: hello
  - name: int75
    type: number
    default: 75
    description: an integer
`

	tests := []struct {
		name      string
		overwrite bool
	}{
		{
			name:      "no overwrite",
			overwrite: false,
		},
		{
			name:      "yes overwrite",
			overwrite: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(yml)
			s, err := newFromReader(r, tt.overwrite)
			require.Nil(t, err)
			require.NotNil(t, s.Ignorer)
			require.NotNil(t, s.Passer)
			require.NotNil(t, s.Variables)

			require.True(t, s.Ignorer.Ignore("myfileA.py"))
			require.True(t, s.Ignorer.Ignore("myDirB/file"))
			require.False(t, s.Ignorer.Ignore("somethingElse.go"))
			if tt.overwrite {
				require.True(t, s.Ignorer.Ignore("manuallyUpdated.txt"))
			} else {
				require.False(t, s.Ignorer.Ignore("manuallyUpdated.txt"))
			}

			require.True(t, s.Passer.Pass("myFileC.cpp"))
			require.True(t, s.Passer.Pass("myExpressionD.h"))
			require.False(t, s.Passer.Pass("somethingElse.go"))

			require.Equal(t, len(s.Variables), 2)
		})
	}
}
