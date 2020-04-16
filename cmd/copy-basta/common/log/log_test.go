package log

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/spin14/copy-basta/cmd/copy-basta/common"

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
		logFunc  func(*Logger)
		expected []string
	}{
		{
			level: Debug,
			logFunc: func(l *Logger) {
				l.Debug(msg)
			},
			expected: []string{msg, "[DEBUG]", string(common.ColorGray)},
		},
		{
			level: Info,
			logFunc: func(l *Logger) {
				l.InfoWithData(msg, LoggerData{"Value From": "LoggingData"})
			},
			expected: []string{msg, "[INFO]", string(common.ColorBlue), "Value From", "LoggingData"},
		},
		{
			level: Warn,
			logFunc: func(l *Logger) {
				l.Warn(msg)
			},
			expected: []string{msg, "[WARN]", string(common.ColorOrange)},
		},
		{
			level: Error,
			logFunc: func(l *Logger) {
				l.ErrorWithData(msg, LoggerData{"Forty-Three": 43})
			},
			expected: []string{msg, "[ERROR]", string(common.ColorRed), "Forty-Three", "43"},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s @ Debug", levelNames[tt.level]), func(t *testing.T) {
			b := &strings.Builder{}
			logger := NewLogger(WithLevel(Debug), WithWriter(b))
			tt.logFunc(logger)
			for _, s := range tt.expected {
				require.Contains(t, b.String(), s)
			}
		})
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s @ Fatal", levelNames[tt.level]), func(t *testing.T) {
			b := &strings.Builder{}
			logger := NewLogger(WithLevel(Fatal), WithWriter(b))
			tt.logFunc(logger)
			require.Equal(t, "", b.String())
		})
	}
}
