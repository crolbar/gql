package main_pane

import (
	"gql/table"
	"gql/tabs"
	"gql/tabs/main_tab/panes"
	"gql/tabs/main_tab/panes/dialog_pane"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KeyMap struct {
    SelectDBTablesPane key.Binding
    DeleteSelectedRow  key.Binding
}

func defaultKeyMap() KeyMap {
    return KeyMap { 
        SelectDBTablesPane: key.NewBinding(
            key.WithKeys("esc"),
            key.WithHelp(", esc", "back to table selection, "),
        ),
        DeleteSelectedRow: key.NewBinding(
            key.WithKeys("d"),
        ),
    }
}

func (km KeyMap) ShortHelp() []key.Binding {
    return []key.Binding{}
}

func (km KeyMap) FullHelp() [][]key.Binding {
    return [][]key.Binding{
        {km.SelectDBTablesPane},
    }
}

func helpView(p panes.Panes) string {
    return lipgloss.JoinHorizontal(lipgloss.Right,
        p.Main.Table.HelpView(),
        p.Main.Help.View(p.Main.KeyMap),
    )
}

func update(p panes.Panes, msg tea.Msg) (panes.Panes, tea.Cmd) {
    var cmd tea.Cmd
    p.Main.Table, cmd = p.Main.Table.Update(msg)

    if cmd != nil {
        switch cmd().(type) {
        case table.UpdatedMsg:
            return p, nil
        }
    }

    keyMap := p.Main.KeyMap.(KeyMap)

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, keyMap.SelectDBTablesPane):
            p.SelectDBTables()
        case key.Matches(msg, keyMap.DeleteSelectedRow):
            cmd = dialog_pane.RequestConfirmation(tabs.DeleteSelectedRow)
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
