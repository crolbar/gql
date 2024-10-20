package util

import (
	"strings"
	"github.com/charmbracelet/lipgloss"
)

func MaxLine(s string) int {
    split := strings.Split(s, "\n")

    m := 0
    for i := 0; i < len(split); i++ {
        currSize := lipgloss.Width(split[i])
        if currSize > m {
            m = currSize
        }
    }

    return m
}

