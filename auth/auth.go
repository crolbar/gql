package auth

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func Init() string {
    m := initialModel()
	p := tea.NewProgram(m)

    fm, err := p.Run(); 
	if err != nil {
		log.Fatal(err)
	}

    if fm.(model).accept {
        return fm.(model).createUri()
    }

    return ""
}

type (
	errMsg error
)

type model struct {
	username  textinput.Model
	password  textinput.Model
	host      textinput.Model
	port      textinput.Model
    accept    bool
	err       error
}

func initTextInput(placeholder string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 156
	ti.Width = 50

    return ti
}

func initialModel() model {
    unameInput := initTextInput("username")
    unameInput.Focus()

	return model{
		username:  unameInput,
		password:  initTextInput("password"),
		host:      initTextInput("host"),
		port:      initTextInput("port"),
        accept:    false,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) createUri() string {
    username := m.username.Value()
    password := m.password.Value()
    host := m.host.Value()
    port := m.port.Value()

    return fmt.Sprintf("%s:%s@(%s:%s)/", username, password, host, port)
}

func(m *model) focusUsername() {
    m.username.Focus()
    m.password.Blur()
    m.host.Blur()
    m.port.Blur()
}

func(m *model) focusPassword() {
    m.username.Blur()
    m.password.Focus()
    m.host.Blur()
    m.port.Blur()
}

func(m *model) focusHost() {
    m.username.Blur()
    m.password.Blur()
    m.host.Focus()
    m.port.Blur()
}

func(m *model) focusPort() {
    m.username.Blur()
    m.password.Blur()
    m.host.Blur()
    m.port.Focus()
}
