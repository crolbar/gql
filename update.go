package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if (m.requiresAuth()) {
        return m.authUpdate(msg)
    }

    return m.mainUpdate(msg)
}

func (m model) authUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
    a, cmd, uri := m.auth.Update(msg)
    m.auth = a;

    if (uri != "") {
        m.uri = uri

        return m, m.openMysql
    }

    return m, cmd
}

func defaultKeyMap() KeyMap {
    return KeyMap {
        Quit: key.NewBinding(
            key.WithKeys("q", "ctrl+c"),
        ),
        ChangeCreds: key.NewBinding(
            key.WithKeys("s"),
        ),
    }
}

func (m model) mainUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
    case dbConnectMsg:
        m.onDBConnect(msg.db)

    case tea.WindowSizeMsg:
        //m.mainTable.UpdateRenderedColums()
        //m.DBTablesTable.UpdateRenderedColums()
        //m.DBTable.UpdateRenderedColums()

    case tea.KeyMsg:
        switch {
        case key.Matches(msg, m.keyMap.Quit):
            return m, tea.Quit
        case key.Matches(msg, m.keyMap.ChangeCreds):
            m.changeCreds()
            return m, nil
        }
	}

    m, cmd = m.getSelectedPane().Update(m, msg)

	return m, cmd
}
