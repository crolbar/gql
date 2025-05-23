package main

import (
	"fmt"
	"os"

	"gql/auth"
	"gql/dbms"

	"gql/tabs"

	"gql/tabs/main_tab/panes/db_pane"
	"gql/tabs/main_tab/panes/db_tables_pane"
	"gql/tabs/main_tab/panes/main_pane"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	keyMap KeyMap
	help   help.Model

	tabs tabs.Tabs

	auth auth.Auth
	user string
	dbms dbms.DBMS

	width  int
	height int
}

type KeyMap struct {
	Quit        key.Binding
	ChangeCreds key.Binding
}

func main() {
	help := help.New()

	help.ShowAll = true
	help.Styles.FullKey = lipgloss.NewStyle().Bold(true)
	help.Styles.FullDesc = lipgloss.NewStyle().Italic(true)
	help.Styles.FullSeparator = lipgloss.NewStyle()
	help.FullSeparator = ""

	uri := getDBUriFromCache()
	m := model{
		keyMap: defaultKeyMap(),
		help:   help,

		tabs: tabs.New(
			db_pane.New(),
			db_tables_pane.New(),
			main_pane.New(),
		),

		auth: auth.InitialAuth(),
		dbms: InitDBMS(uri),
	}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	if m.requiresAuth() {
		return nil
	}

	return m.dbms.Open()
}
