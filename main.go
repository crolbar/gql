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
    Main
)


type model struct {
    selectedPane  Pane

    currDB        string
    currTable     string

    db            *sql.DB
	DBTable       table.Table
	DBTablesTable table.Table

	mainTable     table.Table
}

type dbConnectMsg struct {db *sql.DB}
func (m model) Init() tea.Cmd { return OpenMysql }

func (m *model) selectedTable() table.Table {
    switch (m.selectedPane) {
    case DB:
        return m.DBTable
    case DBTables:
        return m.DBTablesTable
    case Main:
        return m.mainTable
    }

    panic("No pane for the table ?")
}

func (m *model) selectDBpane() {
    m.selectedPane = DB

    m.DBTable.Focus()
    m.DBTablesTable.DeFocus()
}
func (m *model) selectDBTablespane() {
    m.selectedPane = DBTables

    m.DBTable.DeFocus()
    m.mainTable.DeFocus()
    m.DBTablesTable.Focus()
}
func (m *model) selectMainpane() {
    m.selectedPane = Main

    m.DBTablesTable.DeFocus()
    m.mainTable.Focus()
}

func main() {
    m := model {
        selectedPane:  DB,
        DBTablesTable: table.New(nil, nil, 32, 100),
        DBTable:       table.New(nil, nil, 32, 100),
        mainTable:     table.New(nil, nil, 32, 100),
        db:            nil,
    }

    if _, err := tea.NewProgram(m).Run(); err != nil {
        fmt.Println("Error running program:", err)
        os.Exit(1)
    }
}
