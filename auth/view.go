package auth

import "fmt"

func (m model) View() string {
    err := ""

    if m.err != nil {
        err = "\nError: " + m.err.Error()
    }

    fields := fmt.Sprintf(
		"\n%s\n%s\n%s\n%s\n\n%s",
		m.username.View(),
		m.password.View(),
		m.host.View(),
		m.port.View(),
		"([shift]tab for prev/next | esc to quit)",
	)

    return err + "\n" + fields + "\n"
}
