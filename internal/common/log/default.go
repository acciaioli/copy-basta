package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"copy-basta/internal/common"
)

type Level uint

func (l Level) String() string {
	switch l {
	case debugLevel:
		return Debug
	case infoLevel:
		return Info
	case warnLevel:
		return Warn
	case errorLevel:
		return Error
	case fatalLevel:
		return Fatal
	default:
		return ""
	}
}

const (
	noLevel Level = iota
	debugLevel
	infoLevel
	warnLevel
	errorLevel
	fatalLevel
)

const (
	Debug = "debug"
	Info  = "info"
	Warn  = "warn"
	Error = "error"
	Fatal = "fatal"
)

func ToLevel(lvl string) (Level, error) {
	switch lvl {
	case Debug:
		return debugLevel, nil
	case Info:
		return infoLevel, nil
	case Warn:
		return warnLevel, nil
	case Error:
		return errorLevel, nil
	case Fatal:
		return fatalLevel, nil
	default:
		return noLevel, fmt.Errorf("unknown level '%s'", lvl)
	}
}

var L Logger

func init() {
	L = NewLogger()
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
		level:  warnLevel,
		writer: os.Stdout,
		trace:  false,
		levelColors: map[Level]common.Color{
			debugLevel: common.ColorGray,
			infoLevel:  common.ColorBlue,
			warnLevel:  common.ColorOrange,
			errorLevel: common.ColorRed,
			fatalLevel: common.ColorRed,
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
	l.log(debugLevel, nil, msg)
}

func (l *Logger) DebugWithData(msg string, data Data) {
	l.log(debugLevel, data, msg)
}

func (l *Logger) Info(msg string) {
	l.log(infoLevel, nil, msg)
}

func (l *Logger) InfoWithData(msg string, data Data) {
	l.log(infoLevel, data, msg)
}

func (l *Logger) Warn(msg string) {
	l.log(warnLevel, nil, msg)
}

func (l *Logger) WarnWithData(msg string, data Data) {
	l.log(warnLevel, data, msg)
}

func (l *Logger) Error(msg string) {
	l.log(errorLevel, nil, msg)
}

func (l *Logger) ErrorWithData(msg string, data Data) {
	l.log(errorLevel, data, msg)
}

func (l *Logger) Fatal(msg string) {
	l.log(fatalLevel, nil, msg)
	os.Exit(1)
}

func (l *Logger) FatalWithData(msg string, data Data) {
	l.log(fatalLevel, data, msg)
}

func (l *Logger) log(level Level, data Data, userMsg string) {
	if l.level > level {
		return
	}

	color := l.color(level)
	bgColor := l.colorBG(level)
	levelMsg := common.ColoredFormat(color, common.TextFormatBold, bgColor, fmt.Sprintf("[%s]", level.String()))

	lineBuilder := strings.Builder{}
	lineBuilder.WriteString(fmt.Sprintf("%s	%s", levelMsg, userMsg))
	if l.trace || level == debugLevel {
		if _, fn, fl, ok := runtime.Caller(2); ok {
			lineBuilder.WriteString(fmt.Sprintf("\n        @ %s:%d", fn, fl))
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
