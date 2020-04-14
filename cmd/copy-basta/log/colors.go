package log

import "fmt"

type Color string
type BGColor string
type TextFormat string

const (
	ColorNone    Color = ""
	ColorBlack   Color = "30"
	ColorRed     Color = "31"
	ColorGreen   Color = "32"
	ColorOrange  Color = "33"
	ColorBlue    Color = "34"
	ColorMagenta Color = "35"
	ColorCyan    Color = "36"
	ColorGray    Color = "37"
)

const (
	BGColorNone    BGColor = ""
	BGColorBlack   BGColor = "40"
	BGColorRed     BGColor = "41"
	BGColorGreen   BGColor = "42"
	BGColorOrange  BGColor = "43"
	BGColorBlue    BGColor = "44"
	BGColorMagenta BGColor = "45"
	BGColorCyan    BGColor = "46"
	BGColorGray    BGColor = "47"
)

const (
	TextFormatNormal     TextFormat = "0"
	TextFormatBold       TextFormat = "1"
	TextFormatUnderlined TextFormat = "4"
)

// BGColor;TextFormat;Color
const formatExpression = "\033[%s;%s;%sm%s\033[0m"

func ColoredFormat(color Color, tFMT TextFormat, bgColor BGColor, msg string) string {
	return fmt.Sprintf(formatExpression, bgColor, tFMT, color, msg)
}
