package main

import (
	"gql/auth"

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
    var cmd tea.Cmd

    switch msg.(type) {
    case auth.CancelMsg:
        m.uri = getDBUriFromCache()
        return m, m.openMysql
    case auth.Uri:
        m.uri = string(msg.(auth.Uri))
        return m, m.openMysql
    }

    m.auth, cmd = m.auth.Update(msg)

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
        m.width = msg.Width
        m.height = msg.Height

        width  := msg.Width
        height := msg.Height

        height = perc(80, height)

        m.panes.Main.Table.SetMaxSize(perc(60, width), height)

        // only height because we are using only one column
        // and are setting the max width in from UpdateDBTablesTable(): db_util
        m.panes.Db.Table.SetMaxHeight(height)
        m.panes.DbTables.Table.SetMaxHeight(height)

        if (m.db != nil) {
            m.panes.Db.Table.UpdateOffset()
            m.panes.DbTables.Table.UpdateOffset()
            m.panes.Main.Table.UpdateOffset()
        }

        return m, nil

    case tea.KeyMsg:
        switch {
        case key.Matches(msg, m.keyMap.Quit):
            return m, tea.Quit
        case key.Matches(msg, m.keyMap.ChangeCreds):
            m.changeCreds()
            return m, nil
        }
	}

    m.panes, cmd = m.panes.Update(m.db, msg)

	return m, cmd
}

func perc(per, num int) int {
    return int(float32(num) * (float32(per) / float32(100)))
}
