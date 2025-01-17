package logger

import "fmt"

// Foreground colors.
const (
	Black color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Background colors.
// color represents a text color.
type color uint8

// Add adds the coloring to the given string.
func (c color) Add(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}
