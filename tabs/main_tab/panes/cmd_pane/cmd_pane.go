package cmd_pane

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	Cancel        key.Binding
	Accept        key.Binding
	PrevInHistory key.Binding
	NextInHistory key.Binding
}

func defaultKeyMap() KeyMap {
	return KeyMap{
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
		),
		Accept: key.NewBinding(
			key.WithKeys("enter"),
		),
		PrevInHistory: key.NewBinding(
			key.WithKeys("up", "ctrl+p"),
		),
		NextInHistory: key.NewBinding(
			key.WithKeys("down", "ctrl+n"),
		),
	}
}

type Cmd struct {
	keyMap      KeyMap
	textinput   textinput.Model
	history     []string
	currItemIdx int
}

func Init() Cmd {
	ti := textinput.New()
	ti.Prompt = ": "
	return Cmd{
		keyMap:      defaultKeyMap(),
		textinput:   ti,
		history:     make([]string, 0),
		currItemIdx: 0,
	}
}

type CancelMsg struct{}

func Cancel() tea.Msg { return CancelMsg{} }

type AcceptMsg struct {
	Query string
}

func Accept(query string) tea.Cmd {
	return func() tea.Msg {
		return AcceptMsg{
			Query: query,
		}
	}
}

func (c Cmd) Update(msg tea.Msg) (Cmd, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, c.keyMap.Accept):
			value := c.textinput.Value()
			c.history = append(c.history, value)
			c.reset()
			return c, Accept(value)
		case key.Matches(msg, c.keyMap.Cancel):
			c.reset()
			return c, Cancel

		case key.Matches(msg, c.keyMap.PrevInHistory):
			c.SelectPrevInHistory()
		case key.Matches(msg, c.keyMap.NextInHistory):
			c.SelectNextInHistory()
		}
	}

	c.textinput, cmd = c.textinput.Update(msg)

	return c, cmd
}

func (c *Cmd) SelectPrevInHistory() {
	if len(c.history) > 0 {
		if c.currItemIdx == -1 {
			c.currItemIdx = max(len(c.history)-1, 0)
		} else {
			c.currItemIdx = max(c.currItemIdx-1, 0)
		}

		value := c.history[c.currItemIdx]
		c.textinput.SetValue(value)
		c.textinput.SetCursor(len(value))
	}
}

func (c *Cmd) SelectNextInHistory() {
	if len(c.history) > 0 {
		if c.currItemIdx != -1 {
			c.currItemIdx = min(c.currItemIdx+1, len(c.history)-1)
		}

		value := c.history[c.currItemIdx]
		c.textinput.SetValue(value)
		c.textinput.SetCursor(len(value))
	}
}

func (c Cmd) View() string {
	return c.textinput.View()
}

func (c *Cmd) Focus() {
	c.textinput.Focus()
	c.currItemIdx = -1
}

func (c *Cmd) DeFocus() {
	c.textinput.Blur()
}

func (c *Cmd) reset() {
	c.textinput.Reset()
}
