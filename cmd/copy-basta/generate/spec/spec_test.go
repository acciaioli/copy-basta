package spec

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_newFromReader(t *testing.T) {
	tests := []struct {
		name         string
		r            *strings.Reader
		expectedSpec Spec
	}{
		{
			name: "simple",
			r: strings.NewReader(`---
variables:
  stringHi:
    type: string
  float21_6:
    type: number
`),
			expectedSpec: Spec{
				Variables: map[string]Variable{
					"stringHi": {
						Type: "string",
					},
					"float21_6": {
						Type: "number",
					},
				},
			},
		},
		{
			name: "complete",
			r: strings.NewReader(`---
variables:
  stringHello:
    type: string
    description: used to greet
    default: hello
  int75:
    type: number
    default: 75
    description: an integer`),
			expectedSpec: Spec{
				Variables: map[string]Variable{
					"stringHello": {
						Type:        "string",
						Default:     "hello",
						Description: "used to greet",
					},
					"int75": {
						Type:        "number",
						Default:     75,
						Description: "an integer",
					},
				},
			},
		},
		{
			name: "no variables",
			r: strings.NewReader(`---
variables:
`),
			expectedSpec: Spec{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := newFromReader(tt.r)
			require.Nil(t, err)
			require.Equal(t, tt.expectedSpec, *spec)
		})
	}
}

func Test_newFromReader_error(t *testing.T) {
	tests := []struct {
		name string
		r    *strings.Reader
		e    error
	}{
		{
			name: "missing type",
			r: strings.NewReader(`---
variables:
  missingType:
    description: string
`),
			e: errors.New("variable validate error: type is required"),
		},
		{
			name: "invalid type",
			r: strings.NewReader(`---
variables:
  missingType:
    type: not-a-string
`),
			e: errors.New("variable validate error: not-a-string is not a valid type"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newFromReader(tt.r)
			require.Error(t, err)
			t.Log(err.Error())
		})
	}
}