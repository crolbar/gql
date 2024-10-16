package auth

import (
	"gql/util"

	tea "github.com/charmbracelet/bubbletea"
)

func (a Auth) Update(msg tea.Msg) (Auth, tea.Cmd, string) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return a, tea.Quit, ""

        case tea.KeyEnter:
            uri := a.createUri()
            err := util.CheckMysql(uri);

            if err == nil {
                util.WriteToCacheFile(uri)
                return a, nil, uri
            }

            a.err = err

        case tea.KeyTab:
            if (a.port.Focused()) {
                a.focusUsername()
            } else if (a.username.Focused()) {
                a.focusPassword()
            } else if (a.password.Focused()) {
                a.focusHost()
            } else if (a.host.Focused()) {
                a.focusPort()
            }
        case tea.KeyShiftTab:
            if (a.port.Focused()) {
                a.focusHost()
            } else if (a.username.Focused()) {
                a.focusPort()
            } else if (a.password.Focused()) {
                a.focusUsername()
            } else if (a.host.Focused()) {
                a.focusPassword()
            }
		}
    }

	var cmds []tea.Cmd
	var cmd tea.Cmd

	a.username, cmd = a.username.Update(msg)
    cmds = append(cmds, cmd)

	a.password, cmd = a.password.Update(msg)
    cmds = append(cmds, cmd)

	a.host, cmd = a.host.Update(msg)
    cmds = append(cmds, cmd)

	a.port, cmd = a.port.Update(msg)
    cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...), ""
}
