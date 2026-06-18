package style

import (
	"fmt"
	"os"
	"strings"
)

var noColor = os.Getenv("NO_COLOR") != ""

func wrap(code, s string) string {
	if noColor {
		return s
	}
	return fmt.Sprintf("\033[%sm%s\033[0m", code, s)
}

func Bold(s string) string    { return wrap("1", s) }
func Dim(s string) string     { return wrap("2", s) }
func Cyan(s string) string    { return wrap("36", s) }
func Green(s string) string   { return wrap("32", s) }
func Red(s string) string     { return wrap("31", s) }
func Yellow(s string) string  { return wrap("33", s) }
func Header(s string) string  { return Bold(Cyan(s)) }
func Success(s string) string { return Green("✓ ") + s }
func Error(s string) string   { return Red("✗ ") + s }
func Warning(s string) string { return Yellow("! ") + s }

func Table(headers []string, rows [][]string) string {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}
	var b strings.Builder
	for i, h := range headers {
		fmt.Fprintf(&b, "%-*s  ", widths[i], Bold(h))
	}
	b.WriteString("\n")
	for i := range headers {
		b.WriteString(strings.Repeat("─", widths[i]))
		b.WriteString("  ")
	}
	b.WriteString("\n")
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) {
				fmt.Fprintf(&b, "%-*s  ", widths[i], cell)
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}
