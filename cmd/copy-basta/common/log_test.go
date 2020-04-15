package common

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewLogger(t *testing.T) {
	logger, err := NewLogger()
	require.Nil(t, err)
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
			expected: []string{msg, "[DEBUG]", string(ColorGray)},
		},
		{
			level: Info,
			logFunc: func(l *Logger) {
				l.InfoWithData(msg, LoggerData{"Value From": "LoggingData"})
			},
			expected: []string{msg, "[INFO]", string(ColorBlue), "Value From", "LoggingData"},
		},
		{
			level: Warn,
			logFunc: func(l *Logger) {
				l.Warn(msg)
			},
			expected: []string{msg, "[WARN]", "@", string(ColorOrange)},
		},
		{
			level: Error,
			logFunc: func(l *Logger) {
				l.ErrorWithData(msg, LoggerData{"Forty-Three": 43})
			},
			expected: []string{msg, "[ERROR]", "@", string(ColorRed), "Forty-Three", "43"},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s @ Debug", levelNames[tt.level]), func(t *testing.T) {
			b := &strings.Builder{}
			logger, err := NewLogger(WithLevel(Debug), WithWriter(b))
			require.Nil(t, err)
			tt.logFunc(logger)
			for _, s := range tt.expected {
				require.Contains(t, b.String(), s)
			}
		})
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s @ Fatal", levelNames[tt.level]), func(t *testing.T) {
			b := &strings.Builder{}
			logger, err := NewLogger(WithLevel(Fatal), WithWriter(b))
			require.Nil(t, err)
			tt.logFunc(logger)
			require.Equal(t, "", b.String())
		})
	}
}
