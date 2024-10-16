package main

import (
	"gql/table"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type dbTablesPaneKeyMap struct {
    SelectDBTable   key.Binding
    SelectMainTable key.Binding
    Update          key.Binding
}


func dbTablesPaneNew() pane {
    keyMap := dbTablesPaneKeyMap {
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

    update := func(m model, msg tea.Msg) (model, tea.Cmd) {
        var cmd tea.Cmd
        m.dbTablesPane.table, cmd = m.dbTablesPane.table.Update(msg)

        keyMap := m.dbTablesPane.KeyMap.(dbTablesPaneKeyMap)

        switch msg := msg.(type) {
        case tea.KeyMsg:
            switch {
            case key.Matches(msg, keyMap.SelectDBTable):
                m.selectDBPane()

            case key.Matches(msg, keyMap.SelectMainTable):
                m.selectMainPane()
                fallthrough

            case key.Matches(msg, keyMap.Update):
                m.updateMainTable()
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
