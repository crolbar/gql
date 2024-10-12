package main

import (
	"database/sql"
	"fmt"
	"os"

	"gql/table"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle()

type Pane int
const (
    DB Pane = iota
    DBTables
    Table
)


type model struct {
    selectedPane  Pane

    currDB        string
    currTable     string

    db            *sql.DB
	DBTable       *table.Table
	DBTablesTable *table.Table

	mainTable     *table.Table
}

type dbConnectMsg struct {db *sql.DB}
func (m model) Init() tea.Cmd { return OpenMysql }

func main() {
    m := model {
        selectedPane:  DB,
        DBTablesTable: nil,
        DBTable:       nil,
        mainTable:     nil,
        db:            nil,
    }

    if _, err := tea.NewProgram(m).Run(); err != nil {
        fmt.Println("Error running program:", err)
        os.Exit(1)
    }
}
