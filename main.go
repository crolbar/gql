package main

import (
	"database/sql"
	"fmt"
	"os"

	"gql/auth"
	"gql/panes/db_pane"
	"gql/panes/db_tables_pane"
	"gql/panes/main_pane"

	"gql/panes"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    keyMap KeyMap
    help   help.Model

    panes panes.Panes

    auth auth.Auth
    uri  string
    db   *sql.DB
}

type dbConnectMsg struct {db *sql.DB}

type KeyMap struct {
    Quit        key.Binding
    ChangeCreds key.Binding
}

func main() {
    m := model {
        keyMap: defaultKeyMap(),
        help:   help.New(),

        panes: panes.New(
            panes.WithDBPane(db_pane.New()),
            panes.WithDBTablesPane(db_tables_pane.New()),
            panes.WithMainPane(main_pane.New()),
        ),

        auth: auth.InitialAuth(),
        uri:  getDBUriFromCache(),
        db:   nil,
    }

    if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
        fmt.Println("Error running program:", err)
        os.Exit(1)
    }
}

func (m model) Init() tea.Cmd {
    return m.openMysql
}
