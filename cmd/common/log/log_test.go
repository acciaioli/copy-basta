package log

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"copy-basta/cmd/common"

	"github.com/stretchr/testify/require"
)

func Test_NewLogger(t *testing.T) {
	logger := NewLogger()
	require.Equal(t, logger.level, Warn)
	require.Equal(t, logger.writer, os.Stdout)
}

func Test_Logger(t *testing.T) {
	msg := "Hello!"

	tests := []struct {
		level    Level
		logFunc  func(Logger)
		expected []string
	}{
		{
			level: Debug,
			logFunc: func(l Logger) {
				l.Debug(msg)
			},
			expected: []string{msg, "[DEBUG]", string(common.ColorGray)},
		},
		{
			level: Info,
			logFunc: func(l Logger) {
				l.InfoWithData(msg, Data{"Value From": "LoggingData"})
			},
			expected: []string{msg, "[INFO]", string(common.ColorBlue), "Value From", "LoggingData"},
		},
		{
			level: Warn,
			logFunc: func(l Logger) {
				l.Warn(msg)
			},
			expected: []string{msg, "[WARN]", string(common.ColorOrange)},
		},
		{
			level: Error,
			logFunc: func(l Logger) {
				l.ErrorWithData(msg, Data{"Forty-Three": 43})
			},
			expected: []string{msg, "[ERROR]", string(common.ColorRed), "Forty-Three", "43"},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s @ Debug", tt.level.String()), func(t *testing.T) {
			w := &strings.Builder{}
			logger := NewLogger()
			logger.SetLevel(Debug)
			logger.SetWriter(w)
			tt.logFunc(logger)
			for _, s := range tt.expected {
				require.Contains(t, w.String(), s)
			}
		})
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s @ Fatal", tt.level.String()), func(t *testing.T) {
			w := &strings.Builder{}
			logger := NewLogger()
			logger.SetLevel(Fatal)
			logger.SetWriter(w)
			tt.logFunc(logger)
			require.Equal(t, "", w.String())
		})
	}
}
