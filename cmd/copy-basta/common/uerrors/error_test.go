package uerrors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Error_Error(t *testing.T) {
	tests := []struct {
		name          string
		constructor   func() error
		expectedError string
	}{
		{
			name: "Internal",
			constructor: func() error {
				return NewInternalError(FromErr(errors.New("my error")))
			},
			expectedError: "Internal Error",
		},
		{
			name: "Input",
			constructor: func() error {
				return NewInputError("My Error Message", FromString("my error str"))
			},
			expectedError: "User Input Error: My Error Message",
		},
		{
			name: "Template",
			constructor: func() error {
				return NewTemplateError("My Error Message", FromString("my error str"))
			},
			expectedError: "Template Error: My Error Message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.constructor()
			require.NotNil(t, err)
			require.Equal(t, tt.expectedError, err.Error())
		})
	}
}
