package db_pane

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
	SelectDBTable   key.Binding
	Delete          key.Binding
	Filter          key.Binding
	SendCustomQuery key.Binding
	RefreshDBs      key.Binding
}

func defaultKeyMap() KeyMap {
	return KeyMap{
		SelectDBTable: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp(", enter", "table selection, "),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
		),
		SendCustomQuery: key.NewBinding(
			key.WithKeys(":"),
		),
		RefreshDBs: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp(", r", "refresh databases, "),
		),
	}
}

func (km KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (km KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{km.SelectDBTable, km.RefreshDBs},
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

	if cmd != nil {
		switch cmd().(type) {
		case table.CursorMovedMsg:
			cmd = tabs.RequireDBTablesUpdate
		}
	}

	keyMap := p.Db.KeyMap.(KeyMap)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyMap.RefreshDBs):
			cmd = tabs.RequireDBTableUpdate

		case key.Matches(msg, keyMap.SelectDBTable):
			p.SelectDBTables()

		case key.Matches(msg, keyMap.Delete):
			cmd = dialog_pane.RequestConfirmation(tabs.DeleteSelectedDB) // we have to pass cmd not the msg !!

		case key.Matches(msg, keyMap.Filter):
			cmd = tabs.FocusFilter

		case key.Matches(msg, keyMap.SendCustomQuery):
			cmd = tabs.FocusCmd
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
