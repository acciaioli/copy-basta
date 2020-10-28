package specification

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_newFromReader(t *testing.T) {
	tests := []struct {
		name string
		r    *strings.Reader
	}{
		{
			name: "variables",
			r: strings.NewReader(`---
variables:
  - name: stringHi
    dtype: string
  - name: float21_6
    dtype: number
`),
		},
		{
			name: "complete",
			r: strings.NewReader(`---
ignore:
  - myfileA
  - myDirB/
pass:
  - myFileC
  - myExpressionD*
variables:
  - name: stringHello
    dtype: string
    description: used to greet
    default: hello
  - name: int75
    dtype: number
    default: 75
    description: an integer
`),
		},
		{
			name: "no variables",
			r: strings.NewReader(`---
variables:
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := newFromReader(tt.r)
			require.Nil(t, err)
			require.NotNil(t, spec.Ignorer)
			require.NotNil(t, spec.Passer)
			require.NotNil(t, spec.Variables)
		})
	}
}
