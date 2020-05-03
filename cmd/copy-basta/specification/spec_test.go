package specification

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
  - name: stringHi
    type: string
  - name: float21_6
    type: number
`),
			expectedSpec: Spec{
				Variables: []Variable{
					{
						Name: "stringHi",
						Type: func() *openAPIType { v := openAPIType("string"); return &v }(),
					},
					{
						Name: "float21_6",
						Type: func() *openAPIType { v := openAPIType("number"); return &v }(),
					},
				},
			},
		},
		{
			name: "complete",
			r: strings.NewReader(`---
variables:
  - name: stringHello
    type: string
    description: used to greet
    default: hello
  - name: int75
    type: number
    default: 75
    description: an integer`),
			expectedSpec: Spec{
				Variables: []Variable{
					{
						Name:        "stringHello",
						Type:        func() *openAPIType { v := openAPIType("string"); return &v }(),
						Default:     "hello",
						Description: func() *string { s := "used to greet"; return &s }(),
					},
					{
						Name:        "int75",
						Type:        func() *openAPIType { v := openAPIType("number"); return &v }(),
						Default:     75,
						Description: func() *string { s := "an integer"; return &s }(),
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
