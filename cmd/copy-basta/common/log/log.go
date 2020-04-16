package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/spin14/copy-basta/cmd/copy-basta/common"
)

var Log Logger

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
	Fatal
)

var levelNames = map[Level]string{
	Debug: "[DEBUG]",
	Info:  "[INFO]",
	Warn:  "[WARN]",
	Error: "[ERROR]",
	Fatal: "[FATAL]",
}

func sToLevel(s string) (Level, error) {
	switch strings.ToLower(s) {
	case "debug":
		return Debug, nil
	case "info":
		return Info, nil
	case "warn":
		return Warn, nil
	case "error":
		return Error, nil
	case "fatal":
		return Fatal, nil
	default:
		return Fatal, fmt.Errorf("invalid log-level string representation `%s`", s)
	}
}

type Logger struct {
	level         Level
	writer        io.Writer
	trace         bool
	levelColors   map[Level]common.Color
	levelBGColors map[Level]common.BGColor
}

type LoggerData map[string]interface{}

type LoggerOpt func(*Logger)

func WithLevel(level Level) LoggerOpt {
	return func(l *Logger) {
		l.level = level
	}
}

func WithWriter(writer io.Writer) LoggerOpt {
	return func(l *Logger) {
		l.writer = writer
	}
}

func WithTraceData() LoggerOpt {
	return func(l *Logger) {
		l.trace = true
	}
}

func NewLogger(opts ...LoggerOpt) *Logger {
	l := Logger{
		level:  Warn,
		writer: os.Stdout,
		trace:  false,
		levelColors: map[Level]common.Color{
			Debug: common.ColorGray,
			Info:  common.ColorBlue,
			Warn:  common.ColorOrange,
			Error: common.ColorRed,
			Fatal: common.ColorRed,
		},
	}
	for _, o := range opts {
		o(&l)
	}

	l.DebugWithData("new logger created", LoggerData{"level": l.level, "writer is stdout": l.writer == os.Stdout})
	return &l
}

func (l *Logger) Debug(msg string) {
	l.log(Debug, nil, msg)
}

func (l *Logger) DebugWithData(msg string, data LoggerData) {
	l.log(Debug, data, msg)
}

func (l *Logger) Info(msg string) {
	l.log(Info, nil, msg)
}

func (l *Logger) InfoWithData(msg string, data LoggerData) {
	l.log(Info, data, msg)
}

func (l *Logger) Warn(msg string) {
	l.log(Warn, nil, msg)
}

func (l *Logger) WarnWithData(msg string, data LoggerData) {
	l.log(Warn, data, msg)
}

func (l *Logger) Error(msg string) {
	l.log(Error, nil, msg)
}

func (l *Logger) ErrorWithData(msg string, data LoggerData) {
	l.log(Error, data, msg)
}

func (l *Logger) Fatal(msg string) {
	l.log(Fatal, nil, msg)
	os.Exit(1)
}

func (l *Logger) FatalWithData(msg string, data LoggerData) {
	l.log(Fatal, data, msg)
}

func (l *Logger) log(level Level, data LoggerData, userMsg string) {
	if l.level > level {
		return
	}

	color := l.color(level)
	bgColor := l.colorBG(level)
	levelMsg := common.ColoredFormat(color, common.TextFormatBold, bgColor, levelNames[level])

	lineBuilder := strings.Builder{}
	lineBuilder.WriteString(fmt.Sprintf("%s	%s", levelMsg, userMsg))
	if l.trace {
		if _, fn, fl, ok := runtime.Caller(2); ok {
			lineBuilder.WriteString(fmt.Sprintf("	@ %s:%d", fn, fl))
		}
	}
	lineBuilder.WriteString("\n")

	for k, v := range data {
		fmtK := common.ColoredFormat(color, common.TextFormatNormal, bgColor, k)
		lineBuilder.WriteString(fmt.Sprintf("        %s: %v\n", fmtK, v))
	}

	if _, err := fmt.Fprint(l.writer, lineBuilder.String()); err != nil {
		panic(err)
	}
}

func (l *Logger) color(level Level) common.Color {
	if color, ok := l.levelColors[level]; ok {
		return color
	}
	return common.ColorNone
}

func (l *Logger) colorBG(level Level) common.BGColor {
	if color, ok := l.levelBGColors[level]; ok {
		return color
	}
	return common.BGColorNone
}
