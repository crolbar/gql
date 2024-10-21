package db_pane

import (
	"gql/table"
	"gql/tabs/main_tab/panes"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KeyMap struct {
    SelectDBTable key.Binding
    Update        key.Binding
}

func defaultKeyMap() KeyMap {
    return KeyMap {
        SelectDBTable: key.NewBinding(
            key.WithKeys("enter"),
            key.WithHelp(", enter", "table selection, "),
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
        {km.SelectDBTable},
    }
}

func helpView(p panes.Panes) string {
    return lipgloss.JoinHorizontal(lipgloss.Right,
        p.Db.Table.HelpView(),
        p.Db.Help.View(p.Db.KeyMap),
    )
}

func update(p panes.Panes, msg tea.Msg) (panes.Panes, tea.Cmd) {
    var cmd tea.Cmd
    p.Db.Table, cmd = p.Db.Table.Update(msg)

    keyMap := p.Db.KeyMap.(KeyMap)

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, keyMap.SelectDBTable):
            p.SelectDBTables()
            fallthrough

        case key.Matches(msg, keyMap.Update):
            cmd = panes.RequireDBTablesUpdate
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
