package main

import (
	"fmt"
	"os"

    "gql/table"
    "gql/data"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle()

type model struct {
	table table.Table
}


func (m model) Init() tea.Cmd { return nil }


func main() {
    t := table.New(data.ColumnsBig, data.RowsBig, 32, 148)
    m := model{ t, }

    if _, err := tea.NewProgram(m).Run(); err != nil {
        fmt.Println("Error running program:", err)
        os.Exit(1)
    }
}
