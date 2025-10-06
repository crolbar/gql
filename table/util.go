package table

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func isBetween(v, i, j int) bool {
	min := min(i, j)
	max := max(i, j)
	return v >= min && v <= max
}

func iff(cond bool, t, f string) string {
	if cond {
		return t
	}
	return f
}

func clamp(v, low, high int) int {
	return min(max(v, low), high)
}

func truncate(s string, width int) string {
	if lipgloss.Width(s) <= width {
		return s
	}

	var ret strings.Builder
	for i, r := range []rune(s) {
		if i >= width-1 {
			break
		}
		ret.WriteRune(r)
	}

	ret.WriteRune('…')
	return ret.String()
}
