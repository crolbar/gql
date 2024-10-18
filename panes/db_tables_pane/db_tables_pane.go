package db_tables_pane

import (
	"database/sql"
	"gql/panes"
	"gql/table"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
    SelectDBTable   key.Binding
    SelectMainTable key.Binding
    Update          key.Binding
}

func defaultKeyMap() KeyMap {
    return KeyMap {
        SelectDBTable: key.NewBinding(
            key.WithKeys("esc"),
        ),
        SelectMainTable: key.NewBinding(
            key.WithKeys("enter"),
        ),
        Update: key.NewBinding(
            key.WithKeys("j", "k"),
        ),
    }
}

func update(p panes.Panes, db *sql.DB, msg tea.Msg) (panes.Panes, tea.Cmd) {
    var cmd tea.Cmd
    p.DbTables.Table, cmd = p.DbTables.Table.Update(msg)

    keyMap := p.DbTables.KeyMap.(KeyMap)

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, keyMap.SelectDBTable):
            p.SelectDB()

        case key.Matches(msg, keyMap.SelectMainTable):
            p.SelectMain()
            fallthrough

        case key.Matches(msg, keyMap.Update):
            p.UpdateMainTable(db)
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