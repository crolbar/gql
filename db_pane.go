package main

import (
	"gql/table"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type dbPaneKeyMap struct {
    SelectDBTable key.Binding
    Update        key.Binding
}


func dbPaneNew() pane {
    keyMap := dbPaneKeyMap {
        SelectDBTable: key.NewBinding(
            key.WithKeys("enter"),
        ),
        Update: key.NewBinding(
            key.WithKeys("j", "k"),
        ),
    }

    update := func(m model, msg tea.Msg) (model, tea.Cmd) {
        var cmd tea.Cmd
        m.dbPane.table, cmd = m.dbPane.table.Update(msg)

        keyMap := m.dbPane.KeyMap.(dbPaneKeyMap)

        switch msg := msg.(type) {
        case tea.KeyMsg:
            switch {
            case key.Matches(msg, keyMap.SelectDBTable):
                m.selectDBTablesPane()
                fallthrough

            case key.Matches(msg, keyMap.Update):
                m.updateDBTablesTable()
            }
        }

        return m, cmd
    }

    return pane {
        KeyMap: keyMap,
        table: table.New(nil, nil, 32, 100),
        //Help: nil,
        Update: update,
    }
}
