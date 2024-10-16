package main

import (
	"database/sql"
	"fmt"
	"os"

	"gql/auth"
	"gql/table"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type pane struct {
    table  table.Table
    KeyMap interface{}
    Help   help.Model
    Update func(m model, msg tea.Msg) (model, tea.Cmd)
}

type model struct {
    keyMap        KeyMap


    selectedPane  Pane
    dbPane        pane
    dbTablesPane  pane
    mainPane      pane

    db            *sql.DB
    uri           string
    auth          auth.Auth

    currDB        string
    currDBTable     string
}

type dbConnectMsg struct {db *sql.DB}

type KeyMap struct {
    Quit        key.Binding
    ChangeCreds key.Binding
}

func main() {
    m := model {
        keyMap:        defaultKeyMap(),

        selectedPane:  DB,
        dbPane:        dbPaneNew(),
        dbTablesPane:  dbTablesPaneNew(),
        mainPane:      mainPaneNew(),

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
    return m.openMysql
}
