package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"copy-basta/cmd/common"
)

var L Logger

func init() {
	L = NewLogger()
}

const (
	Debug Level = iota
	Info
	Warn
	Error
	Fatal

	debugS = "DEBUG"
	infoS  = "INFO"
	warnS  = "WARN"
	errorS = "ERROR"
	fatalS = "FATAL"
)

func StringToLevel(s string) (Level, error) {
	switch strings.ToUpper(s) {
	case debugS:
		return Debug, nil
	case infoS:
		return Info, nil
	case warnS:
		return Warn, nil
	case errorS:
		return Error, nil
	case fatalS:
		return Fatal, nil
	default:
		return Fatal, fmt.Errorf("%s is not a known level", s)
	}
}

type Level int

func (lvl Level) String() string {
	return fmt.Sprintf("[%s]", []string{
		debugS,
		infoS,
		warnS,
		errorS,
		fatalS,
	}[lvl])
}

type Data map[string]interface{}

type Logger struct {
	level         Level
	writer        io.Writer
	trace         bool
	levelColors   map[Level]common.Color
	levelBGColors map[Level]common.BGColor
}

func NewLogger() Logger {
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
		levelBGColors: map[Level]common.BGColor{},
	}

	l.DebugWithData("new logger created", Data{"level": l.level, "writer is stdout": l.writer == os.Stdout})
	return l
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) SetWriter(writer io.Writer) {
	l.writer = writer
}

func (l *Logger) EnableTrace() {
	l.trace = true
}

func (l *Logger) DisableTrace() {
	l.trace = true
}

func (l *Logger) SetColor(level Level, color common.Color) {
	l.levelColors[level] = color
}

func (l *Logger) SetBGColor(level Level, bgColor common.BGColor) {
	l.levelBGColors[level] = bgColor
}

func (l *Logger) Debug(msg string) {
	l.log(Debug, nil, msg)
}

func (l *Logger) DebugWithData(msg string, data Data) {
	l.log(Debug, data, msg)
}

func (l *Logger) Info(msg string) {
	l.log(Info, nil, msg)
}

func (l *Logger) InfoWithData(msg string, data Data) {
	l.log(Info, data, msg)
}

func (l *Logger) Warn(msg string) {
	l.log(Warn, nil, msg)
}

func (l *Logger) WarnWithData(msg string, data Data) {
	l.log(Warn, data, msg)
}

func (l *Logger) Error(msg string) {
	l.log(Error, nil, msg)
}

func (l *Logger) ErrorWithData(msg string, data Data) {
	l.log(Error, data, msg)
}

func (l *Logger) Fatal(msg string) {
	l.log(Fatal, nil, msg)
	os.Exit(1)
}

func (l *Logger) FatalWithData(msg string, data Data) {
	l.log(Fatal, data, msg)
}

func (l *Logger) log(level Level, data Data, userMsg string) {
	if l.level > level {
		return
	}

	color := l.color(level)
	bgColor := l.colorBG(level)
	levelMsg := common.ColoredFormat(color, common.TextFormatBold, bgColor, level.String())

	lineBuilder := strings.Builder{}
	lineBuilder.WriteString(fmt.Sprintf("%s	%s", levelMsg, userMsg))
	if l.trace || level == Debug {
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
