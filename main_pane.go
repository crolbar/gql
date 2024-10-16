package main

import (
	"gql/table"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type mainTablePaneKeyMap struct {
    SelectDBTablesPane key.Binding
}

func mainPaneNew() pane {
    keyMap := mainTablePaneKeyMap {
        SelectDBTablesPane: key.NewBinding(
            key.WithKeys("esc"),
        ),
    }

    update := func(m model, msg tea.Msg) (model, tea.Cmd) {
        var cmd tea.Cmd
        m.mainPane.table, cmd = m.mainPane.table.Update(msg)

        keyMap := m.mainPane.KeyMap.(mainTablePaneKeyMap)

        switch msg := msg.(type) {
        case tea.KeyMsg:
            switch {
            case key.Matches(msg, keyMap.SelectDBTablesPane):
                m.selectDBTablesPane()
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
