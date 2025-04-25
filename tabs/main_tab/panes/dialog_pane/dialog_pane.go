package dialog_pane

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KeyMap struct {
	Cancel key.Binding
	Accept key.Binding
}

func defaultKeyMap() KeyMap {
	return KeyMap{
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
	help         string
	err          string
	returnCmd    tea.Cmd
	keyMap       KeyMap
}

func Init() Dialog {
	ti := textinput.New()

	ti.CharLimit = 156
	ti.Width = 80
	ti.Focus()

	return Dialog{
		textinput: ti,
		keyMap:    defaultKeyMap(),
		err:       "",
	}
}

type CancelMsg struct{}

func Cancel() tea.Msg { return CancelMsg{} }

type AcceptValueUpdateMsg struct {
	Cmd   tea.Cmd
	Value string
}

func AcceptValueUpdate(cmd tea.Cmd, value string) tea.Cmd {
	return func() tea.Msg {
		return AcceptValueUpdateMsg{
			Cmd:   cmd,
			Value: value,
		}
	}
}

type RequestValueUpdateMsg struct {
	Cmd tea.Cmd
}

func RequestValueUpdate(cmd tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		return RequestValueUpdateMsg{
			Cmd: cmd,
		}
	}
}

type RequestConfirmationMsg struct {
	Cmd tea.Cmd
}

func RequestConfirmation(cmd tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		return RequestConfirmationMsg{
			Cmd: cmd,
		}
	}
}

func (d Dialog) handleConfirmationAccept() (Dialog, tea.Cmd) {
	if d.textinput.Value() == "yes" {
		return d, d.returnCmd
	} else if d.textinput.Value() == "no" {
		d.Reset()
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

			return d, AcceptValueUpdate(d.returnCmd, d.textinput.Value())
		case key.Matches(msg, d.keyMap.Cancel):
			d.Reset()
			return d, Cancel
		}
	}

	d.textinput, cmd = d.textinput.Update(msg)

	return d, cmd
}

func (d *Dialog) SetupConfirmation(cmd tea.Cmd, help string) {
	d.confirmation = true
	d.returnCmd = cmd
	d.textinput.Placeholder = "yes/no"
	d.help = help
}

func (d *Dialog) SetupValueUpdate(cmd tea.Cmd, help, currVal string) {
	d.confirmation = false
	d.returnCmd = cmd
	d.textinput.SetValue(currVal)
	d.help = help
}

func (d *Dialog) GetHelpMsg() string {
	return d.help
}

func (d *Dialog) SetWidth(
	width,
	dbPaneWidth,
	dbTablesPaneWidth,
	mainPaneWidth int,
) {
	tablesWidth := dbPaneWidth + dbTablesPaneWidth + mainPaneWidth
	dialogWidth := (width - tablesWidth) - ((1 + 1) + 2 + 1)

	d.textinput.Width = dialogWidth
}

func (d *Dialog) Reset() {
	d.confirmation = false
	d.textinput.Placeholder = ""
	d.err = ""
	d.textinput.Reset()
}

func (d Dialog) TextInputView(style lipgloss.Style) string {
	err := style.
		Foreground(lipgloss.Color("1")).
		Render(d.err)

	return "\n" + d.textinput.View() + "\n\n" + err
}

func (d *Dialog) SetError(err string) {
	d.err = err
}
