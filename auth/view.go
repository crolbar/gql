package auth

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
)


func (km KeyMap) ShortHelp() []key.Binding {
    return []key.Binding { km.Accept, km.Next, km.Prev, km.Quit }
}

func (km KeyMap) FullHelp() [][]key.Binding {
    return [][]key.Binding {
        {km.Accept, km.Cancel},
        {km.Next, km.Prev},
        {km.Quit},
    }
}

func (a Auth) View() string {
    err := ""

    if a.err != nil {
        err = "\nError: " + a.err.Error()
    }

    fields := fmt.Sprintf(
        "\n%s\n%s\n%s\n%s\n\n",
        a.username.View(),
        a.password.View(),
        a.host.View(),
        a.port.View(),
    )

    help := a.Help.View(a.KeyMap)

    return err + "\n" + fields + "\n" + help + "\n"
}
