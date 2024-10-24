package filter_pane

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

type Filter struct {
    keyMap KeyMap
    textinput textinput.Model
}

func Init() Filter {
    ti := textinput.New()
    ti.Prompt = "> where "
    return Filter {
        keyMap:    defaultKeyMap(),
        textinput: ti,
    }
}


type CancelMsg struct{}
func Cancel() tea.Msg { return CancelMsg{} }

type AcceptMsg struct {
    Txt string
}
func Accept(txt string) tea.Cmd { 
    return func() tea.Msg {
        return AcceptMsg{
            Txt: txt,
        } 
    }
}


func (f Filter) Update(msg tea.Msg) (Filter, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, f.keyMap.Accept):
            return f, Accept(f.textinput.Value())
        case key.Matches(msg, f.keyMap.Cancel):
            f.reset()
            return f, Cancel
        }
    }

    f.textinput, cmd = f.textinput.Update(msg)

    return f, cmd
}

func (f Filter) View() string {
    return f.textinput.View()
}

func (f *Filter) SetWidth(width int) {
    f.textinput.Width = (width / 2) - len(f.textinput.Prompt)
}

func (f *Filter) Focus() {
    f.textinput.Focus()
}

func (f *Filter) DeFocus() {
    f.textinput.Blur()
}

func (f *Filter) UpdateValue(value string) {
    f.textinput.SetValue(value)
}

func (f *Filter) UpdatePrefix(prefix string) {
    f.textinput.Prompt = "> where " + prefix
}

func (f *Filter) reset() {
    f.textinput.Reset()
}
