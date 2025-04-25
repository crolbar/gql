package auth

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
)

type Auth struct {
	textinput textinput.Model

	KeyMap KeyMap
	Help   help.Model
	err    error
}

func InitialAuth() Auth {
	ti := textinput.New()

	ti.Placeholder = "connection string / uri"
	ti.Width = 50
	ti.Focus()

	help := help.New()
	help.ShowAll = true

	return Auth{
		textinput: ti,

		KeyMap: DefaultKeyMap(),
		Help:   help,
		err:    nil,
	}
}

func (a *Auth) Reset(uri string) {
	a.textinput.Reset()
	a.textinput.SetValue(uri)
	a.err = nil
}
