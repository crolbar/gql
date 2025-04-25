package main

import (
	"gql/auth"
	"gql/dbms"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = msg.Width
		m.height = msg.Height
	}

    if (m.requiresAuth()) {
        return m.authUpdate(msg)
    }

    return m.mainUpdate(msg)
}

func (m model) authUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg.(type) {
    case auth.CancelMsg:
        m.dbms = InitDBMS(getDBUriFromCache())
		if m.dbms == nil {
			return m, nil
		}
        return m, m.dbms.Open()
    case auth.Uri:
        m.dbms = InitDBMS(string(msg.(auth.Uri)))
        return m, m.dbms.Open()
    }

    m.auth, cmd = m.auth.Update(msg)

    return m, cmd
}

func defaultKeyMap() KeyMap {
    return KeyMap {
        Quit: key.NewBinding(
            key.WithKeys("q", "ctrl+c"),
            key.WithHelp("q/ctrl+c", "quit"),
        ),
        ChangeCreds: key.NewBinding(
            key.WithKeys("s"),
            key.WithHelp("s", "change uri"),
        ),
    }
}

func (m model) onWindowRisize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
    m.tabs.OnWindowResize(msg, m.dbms.HasDb())

    return m, nil
}

func (m model) mainUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
    case dbms.DbConnectMsg:
        m.onDBConnect(msg.Db)

    case tea.WindowSizeMsg:
        return m.onWindowRisize(msg)

    case tea.KeyMsg:
        switch {
        case key.Matches(msg, m.keyMap.Quit):
            if m.tabs.IsTyting() && msg.String() == "q" {
                break;
            }

            return m, tea.Quit
        case key.Matches(msg, m.keyMap.ChangeCreds):
            if m.tabs.IsTyting() {
                break;
            }

            m.changeCreds()
            return m, nil
        }
	}

    m.tabs, cmd = m.tabs.Update(m.dbms, msg)

	return m, cmd
}

func perc(per, num int) int {
    return int(float32(num) * (float32(per) / float32(100)))
}
