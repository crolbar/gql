package main

import (
	"database/sql"
	"fmt"
	"os"

	"gql/auth"
	"gql/table"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    selectedPane  Pane

    uri           string
    currDB        string
    currTable     string

    db            *sql.DB
	DBTable       table.Table
	DBTablesTable table.Table
	mainTable     table.Table

    auth          auth.Auth
}

type dbConnectMsg struct {db *sql.DB}

func main() {
    m := model {
        selectedPane:  DB,
        DBTablesTable: table.New(nil, nil, 32, 100),
        DBTable:       table.New(nil, nil, 32, 100),
        mainTable:     table.New(nil, nil, 32, 100),
        db:            nil,
        uri:           getDBUriFromCache(),
        auth:          auth.InitialAuth(),
    }

    if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
        fmt.Println("Error running program:", err)
        os.Exit(1)
    }
}

func (m model) Init() tea.Cmd {
    return m.OpenMysql
}
