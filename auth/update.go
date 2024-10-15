package auth

import (
	"gql/util"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
        case tea.KeyEnter:
            uri := m.createUri()
            err := util.CheckMysql(uri);

            if err == nil {
                util.WriteToCacheFile(uri)
                m.accept = true
                return m, tea.Quit
            }

            m.err = err
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

        case tea.KeyTab:
            if (m.port.Focused()) {
                m.focusUsername()
            } else if (m.username.Focused()) {
                m.focusPassword()
            } else if (m.password.Focused()) {
                m.focusHost()
            } else if (m.host.Focused()) {
                m.focusPort()
            }
        case tea.KeyShiftTab:
            if (m.port.Focused()) {
                m.focusHost()
            } else if (m.username.Focused()) {
                m.focusPort()
            } else if (m.password.Focused()) {
                m.focusUsername()
            } else if (m.host.Focused()) {
                m.focusPassword()
            }
		}
    }

	m.username, cmd = m.username.Update(msg)
	m.password, cmd = m.password.Update(msg)
	m.host, cmd = m.host.Update(msg)
	m.port, cmd = m.port.Update(msg)
	return m, cmd
}
