package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"copy-basta/internal/common/errors"
)

func TestErrorBuilder(t *testing.T) {
	errBuilder := errors.NewErrorBuilder("test")

	tests := []struct {
		name                  string
		format                string
		a                     []interface{}
		expectedUserError     string
		expectedInternalError string
	}{
		{
			name:                  "empty",
			format:                "",
			a:                     nil,
			expectedUserError:     "user error\n[test] ",
			expectedInternalError: "internal error\n[test] ",
		},
		{
			name:                  "without format",
			format:                "you did bad",
			a:                     nil,
			expectedUserError:     "user error\n[test] you did bad",
			expectedInternalError: "internal error\n[test] you did bad",
		},
		{
			name:                  "with format",
			format:                "you did bad %d times",
			a:                     []interface{}{2},
			expectedUserError:     "user error\n[test] you did bad 2 times",
			expectedInternalError: "internal error\n[test] you did bad 2 times",
		},
		{
			name:                  "with format err",
			format:                "%w",
			a:                     []interface{}{fmt.Errorf("you bad boy")},
			expectedUserError:     "user error\n[test] you bad boy",
			expectedInternalError: "internal error\n[test] you bad boy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userErr := errBuilder.NewUserErrorf(tt.format, tt.a...)
			require.NotNil(t, userErr)
			require.Equal(t, tt.expectedUserError, userErr.Error())

			internalErr := errBuilder.NewInternalErrorf(tt.format, tt.a...)
			require.NotNil(t, internalErr)
			require.Equal(t, tt.expectedInternalError, internalErr.Error())
		})
	}
}
