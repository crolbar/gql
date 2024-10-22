package dialog_pane

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)


type KeyMap struct {
    Cancel key.Binding
    Accept key.Binding
}

func defaultKeyMap() KeyMap {
    return KeyMap { 
        Cancel: key.NewBinding(
            key.WithKeys("esc"),
        ),
        Accept: key.NewBinding(
            key.WithKeys("enter"),
        ),
    }
}

type Dialog struct {
    confirmation bool
	textinput    textinput.Model
	err          string
    returnCmd    tea.Cmd
    keyMap       KeyMap
}

func InitDialog() Dialog {
	ti := textinput.New()

	ti.CharLimit   = 156
	ti.Width       = 25
    ti.Focus()

    return Dialog {
        textinput: ti,
        keyMap:    defaultKeyMap(),
        err:       "",
    }
}

type CancelMsg struct{}
func Cancel() tea.Msg { return CancelMsg{} }

type RequestConfirmationMsg struct {
    Cmd tea.Cmd
}
func RequestConfirmation(cmd tea.Cmd) tea.Cmd {
    return func() tea.Msg {
        return RequestConfirmationMsg {
            Cmd: cmd,
        }
    }
}

func (d Dialog) handleConfirmationAccept() (Dialog, tea.Cmd) {
    if d.textinput.Value() == "yes" {
        d.reset()
        return d, d.returnCmd
    } else if d.textinput.Value() == "no" {
        d.reset()
        return d, Cancel
    } else {
        d.err = "yes or no"
    }
    return d, nil
}

func (d Dialog) Update(msg tea.Msg) (Dialog, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, d.keyMap.Accept):
            if d.confirmation {
                return d.handleConfirmationAccept()
            }

        case key.Matches(msg, d.keyMap.Cancel):
            d.reset()
            return d, Cancel
        }
    }

    d.textinput, cmd = d.textinput.Update(msg)

    return d, cmd
}

func (d *Dialog) SetupConfirmation(cmd tea.Cmd) {
    d.confirmation          = true
    d.returnCmd             = cmd
    d.textinput.Placeholder = "yes/no"
}

func (d *Dialog) reset() {
    d.confirmation          = false
    d.textinput.Placeholder = ""
    d.textinput.Reset()
}

func (d Dialog) TextInputView() string {
    return d.textinput.View() + "\n" + d.err
}
