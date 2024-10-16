package auth

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
)

type Auth struct {
	username  textinput.Model
	password  textinput.Model
	host      textinput.Model
	port      textinput.Model

	err       error
}

func initTextInput(placeholder string) textinput.Model {
	ti := textinput.New()

	ti.Placeholder = placeholder
	ti.CharLimit   = 156
	ti.Width       = 50

    return ti
}

func InitialAuth() Auth {
    unameInput := initTextInput("username")
    unameInput.Focus()

	return Auth {
		username:  unameInput,
		password:  initTextInput("password"),
		host:      initTextInput("host"),
		port:      initTextInput("port"),
		err:       nil,
	}
}

func (a *Auth) ResetAll() {
    a.username.Reset()
    a.password.Reset()
    a.host.Reset()
    a.port.Reset()
    a.focusUsername()
    a.err = nil
}

func (a Auth) createUri() string {
    username := a.username.Value()
    password := a.password.Value()
    host     := a.host.Value()
    port     := a.port.Value()

    return fmt.Sprintf("%s:%s@(%s:%s)/", username, password, host, port)
}

func(a *Auth) focusUsername() {
    a.username.Focus()
    a.password.Blur()
    a.host.Blur()
    a.port.Blur()
}

func(a *Auth) focusPassword() {
    a.username.Blur()
    a.password.Focus()
    a.host.Blur()
    a.port.Blur()
}

func(a *Auth) focusHost() {
    a.username.Blur()
    a.password.Blur()
    a.host.Focus()
    a.port.Blur()
}

func(a *Auth) focusPort() {
    a.username.Blur()
    a.password.Blur()
    a.host.Blur()
    a.port.Focus()
}
