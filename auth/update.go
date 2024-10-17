package auth

import (
	"gql/util"

    "github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)


type KeyMap struct {
    Quit    key.Binding
    Cancel  key.Binding
    Accept  key.Binding
    Next    key.Binding
    Prev    key.Binding
}

func DefaultKeyMap() KeyMap {
    return KeyMap {
        Quit: key.NewBinding(
            key.WithKeys("ctrl+c"),
            key.WithHelp("ctrl+c", "quit"),
        ),
        Cancel: key.NewBinding(
            key.WithKeys("esc"),
            key.WithHelp("esc", "cancel"),
        ),
        Accept: key.NewBinding(
            key.WithKeys("enter"),
            key.WithHelp("enter", "accept"),
        ),
        Next: key.NewBinding(
            key.WithKeys("tab", "down"),
            key.WithHelp("tab/down", "next field"),
        ),
        Prev: key.NewBinding(
            key.WithKeys("shift+tab", "up"),
            key.WithHelp("shift+tab/up", "prev field"),
        ),
    }
}

type Uri string

func (a Auth) Accept() (Auth, tea.Cmd) {
    uri := a.createUri()
    err := util.CheckMysql(uri);

    if err != nil {
        a.err = err
        return a, nil
    }

    util.WriteToCacheFile(uri)
    return a,
    func() tea.Msg {
        return Uri(uri)
    }
}

func (a *Auth) NextField() {
    if (a.port.Focused()) {
        a.focusUsername()
    } else if (a.username.Focused()) {
        a.focusPassword()
    } else if (a.password.Focused()) {
        a.focusHost()
    } else if (a.host.Focused()) {
        a.focusPort()
    }
}

func (a *Auth) PrevField() {
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

type CancelMsg struct{}

func (a Auth) Update(msg tea.Msg) (Auth, tea.Cmd) {
	switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, a.KeyMap.Quit):
            return a, tea.Quit

        case key.Matches(msg, a.KeyMap.Cancel):
            return a,
            func() tea.Msg {
                return CancelMsg{}
            }

        case key.Matches(msg, a.KeyMap.Accept):
            return a.Accept()

        case key.Matches(msg, a.KeyMap.Next):
            a.NextField()

        case key.Matches(msg, a.KeyMap.Prev):
            a.PrevField()
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

	return a, tea.Batch(cmds...)
}
