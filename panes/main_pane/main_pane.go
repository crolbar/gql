package main_pane

import (
	"database/sql"
	"gql/panes"
	"gql/table"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
    SelectDBTablesPane key.Binding
}

func defaultKeyMap() KeyMap {
    return KeyMap { 
        SelectDBTablesPane: key.NewBinding(
            key.WithKeys("esc"),
        ),
    }
}

func update(p panes.Panes, db *sql.DB, msg tea.Msg) (panes.Panes, tea.Cmd) {
    var cmd tea.Cmd
    p.Main.Table, cmd = p.Main.Table.Update(msg)

    keyMap := p.Main.KeyMap.(KeyMap)

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, keyMap.SelectDBTablesPane):
            p.SelectDBTables()
        }
    }

    return p, cmd
}

func New() panes.Pane {
    return panes.NewPane(
        table.New(nil, nil, 32, 100),
        defaultKeyMap(),
        help.New(),
        update,
    )
}
