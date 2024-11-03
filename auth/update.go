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
}

func DefaultKeyMap() KeyMap {
    return KeyMap {
        Quit: key.NewBinding(
            key.WithKeys("ctrl+c"),
            key.WithHelp("    ctrl+c", "quit"),
        ),
        Cancel: key.NewBinding(
            key.WithKeys("esc"),
            key.WithHelp("esc", "cancel"),
        ),
        Accept: key.NewBinding(
            key.WithKeys("enter"),
            key.WithHelp("enter", "accept"),
        ),
    }
}

type Uri string

func (a Auth) Accept() (Auth, tea.Cmd) {
    uri := a.textinput.Value()
    err := util.CheckDBMS(uri);

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
        }
    }

	var cmd tea.Cmd
	a.textinput, cmd = a.textinput.Update(msg)
	return a, cmd
}
