package auth

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
)

func (km KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{km.Accept, km.Quit}
}

func (km KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{km.Accept, km.Cancel},
		{km.Quit},
	}
}

func (a Auth) View() string {
	err := ""

	if a.err != nil {
		err = "\nError: " + a.err.Error()
	}

	fields := fmt.Sprintf(
		"\n%s\n\n",
		a.textinput.View(),
	)

	help := a.Help.View(a.KeyMap)

	return err + "\n" + fields + "\n" + help + "\n"
}
