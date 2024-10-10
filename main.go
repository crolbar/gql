package main

import (
	"fmt"
	"os"

    "gql/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle()

type model struct {
	table table.Table
}


func (m model) Init() tea.Cmd { return nil }


func main() {
    t := table.New(columnsBig, rowsBig, 32, 200)
    m := model{ t, }

    if _, err := tea.NewProgram(m).Run(); err != nil {
        fmt.Println("Error running program:", err)
        os.Exit(1)
    }
}
