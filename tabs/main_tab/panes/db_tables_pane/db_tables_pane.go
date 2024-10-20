package db_tables_pane

import (
	"gql/table"
	"gql/tabs/main_tab/panes"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
            key.WithHelp(", esc", "back to db selection, "),
        ),
        SelectMainTable: key.NewBinding(
            key.WithKeys("enter"),
            key.WithHelp(", enter", "selected table, "),
        ),
        Update: key.NewBinding(
            key.WithKeys("j", "k"),
        ),
    }
}

func (km KeyMap) ShortHelp() []key.Binding {
    return []key.Binding{}
}

func (km KeyMap) FullHelp() [][]key.Binding {
    return [][]key.Binding{
        {km.SelectMainTable, km.SelectDBTable},
    }
}

func helpView(p panes.Panes) string {
    return lipgloss.JoinHorizontal(lipgloss.Right,
        p.DbTables.Table.HelpView(),
        p.DbTables.Help.View(p.DbTables.KeyMap),
    )
}

func update(p panes.Panes, msg tea.Msg) (panes.Panes, tea.Cmd) {
    var cmd tea.Cmd
    p.DbTables.Table, cmd = p.DbTables.Table.Update(msg)

    if cmd != nil {
        switch cmd().(type) {
        case table.Updated:
            return p, nil
        }
    }

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
            cmd = panes.RequireMainTableUpdate
        }
    }

    return p, cmd
}

func New() panes.Pane {
    return panes.NewPane(
        table.New(nil, nil, 32, 100),
        defaultKeyMap(),
        update,
        helpView,
    )
}
