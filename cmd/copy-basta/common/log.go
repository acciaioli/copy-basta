package common

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
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
		return Fatal, fmt.Errorf("invalid log level representation `%s`", s)
	}
}

type Logger struct {
	level         Level
	writer        io.Writer
	levelColors   map[Level]Color
	levelBGColors map[Level]BGColor
}

type LoggerData map[string]interface{}

type LoggerOpt func(*Logger) error

func WithLevel(level Level) LoggerOpt {
	return func(l *Logger) error {
		l.level = level
		return nil
	}
}

func WithLevelS(level string) LoggerOpt {
	return func(l *Logger) error {
		level, err := sToLevel(level)
		if err != nil {
			return err
		}
		l.level = level
		return nil
	}
}

func WithWriter(writer io.Writer) LoggerOpt {
	return func(l *Logger) error {
		l.writer = writer
		return nil
	}
}

func NewLogger(opts ...LoggerOpt) (*Logger, error) {
	l := Logger{
		level:  Warn,
		writer: os.Stdout,
		levelColors: map[Level]Color{
			Debug: ColorGray,
			Info:  ColorBlue,
			Warn:  ColorOrange,
			Error: ColorRed,
			Fatal: ColorRed,
		},
	}
	for _, o := range opts {
		if err := o(&l); err != nil {
			return nil, err
		}
	}

	l.DebugWithData("new logger created", LoggerData{"level": l.level, "writer is stdout": l.writer == os.Stdout})
	return &l, nil
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
	levelMsg := ColoredFormat(color, TextFormatBold, bgColor, levelNames[level])
	runtimeMsg := ""
	if level > Info {
		if _, fn, fl, ok := runtime.Caller(2); ok {
			runtimeMsg = fmt.Sprintf("@ %s:%d", fn, fl)
		}
	}

	if _, err := fmt.Fprintf(l.writer, "%s	%s	%s\n", levelMsg, userMsg, runtimeMsg); err != nil {
		panic(err)
	}
	for k, v := range data {
		fmtK := ColoredFormat(color, TextFormatNormal, bgColor, k)
		if _, err := fmt.Fprintf(l.writer, "        %s: %v\n", fmtK, v); err != nil {
			panic(err)
		}
	}
}

func (l *Logger) color(level Level) Color {
	if color, ok := l.levelColors[level]; ok {
		return color
	}
	return ColorNone
}

func (l *Logger) colorBG(level Level) BGColor {
	if color, ok := l.levelBGColors[level]; ok {
		return color
	}
	return BGColorNone
}
