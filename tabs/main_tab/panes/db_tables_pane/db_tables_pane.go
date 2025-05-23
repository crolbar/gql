package db_tables_pane

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
	SelectMainTable key.Binding
	Filter          key.Binding
	Delete          key.Binding
	Rename          key.Binding
	SendCustomQuery key.Binding
	RefreshDbTables key.Binding
}

func defaultKeyMap() KeyMap {
	return KeyMap{
		SelectDBTable: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp(", esc", "back to db selection, "),
		),
		SelectMainTable: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp(", enter", "selected table, "),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
		),
		Rename: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "rename table, "),
		),
		SendCustomQuery: key.NewBinding(
			key.WithKeys(":"),
		),
		RefreshDbTables: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh tables, "),
		),
	}
}

func (km KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (km KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{km.SelectMainTable, km.SelectDBTable},
		{km.Rename, km.RefreshDbTables},
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
		case table.UpdatedMsg:
			return p, nil
		case table.CursorMovedMsg:
			cmd = tabs.RequireMainTableUpdate
		}
	}

	keyMap := p.DbTables.KeyMap.(KeyMap)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyMap.RefreshDbTables):
			cmd = tabs.RequireDBTablesUpdate

		case key.Matches(msg, keyMap.SelectDBTable):
			p.SelectDB()

		case key.Matches(msg, keyMap.SelectMainTable):
			p.SelectMain()

		case key.Matches(msg, keyMap.Delete):
			cmd = dialog_pane.RequestConfirmation(tabs.DeleteSelectedDBTable)

		case key.Matches(msg, keyMap.Rename):
			cmd = dialog_pane.RequestValueUpdate(tabs.ChangeDbTableName)

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
