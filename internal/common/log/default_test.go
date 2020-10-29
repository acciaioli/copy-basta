package log

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"copy-basta/internal/common"

	"github.com/stretchr/testify/require"
)

func Test_NewLogger(t *testing.T) {
	logger := NewLogger()
	require.Equal(t, logger.level, warnLevel)
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
			level: debugLevel,
			logFunc: func(l Logger) {
				l.Debug(msg)
			},
			expected: []string{msg, "[debug]", string(common.ColorGray)},
		},
		{
			level: infoLevel,
			logFunc: func(l Logger) {
				l.InfoWithData(msg, Data{"Value From": "LoggingData"})
			},
			expected: []string{msg, "[info]", string(common.ColorBlue), "Value From", "LoggingData"},
		},
		{
			level: warnLevel,
			logFunc: func(l Logger) {
				l.Warn(msg)
			},
			expected: []string{msg, "[warn]", string(common.ColorOrange)},
		},
		{
			level: errorLevel,
			logFunc: func(l Logger) {
				l.ErrorWithData(msg, Data{"Forty-Three": 43})
			},
			expected: []string{msg, "[error]", string(common.ColorRed), "Forty-Three", "43"},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s @ debug", tt.level.String()), func(t *testing.T) {
			w := &strings.Builder{}
			logger := NewLogger()
			logger.SetLevel(debugLevel)
			logger.SetWriter(w)
			tt.logFunc(logger)
			for _, s := range tt.expected {
				require.Contains(t, w.String(), s)
			}
		})
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s @ fatal", tt.level.String()), func(t *testing.T) {
			w := &strings.Builder{}
			logger := NewLogger()
			logger.SetLevel(fatalLevel)
			logger.SetWriter(w)
			tt.logFunc(logger)
			require.Equal(t, "", w.String())
		})
	}
}
