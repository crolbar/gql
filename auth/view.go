package auth

import (
	"fmt"
)

func (a Auth) View() string {
    err := ""

    if a.err != nil {
        err = "\nError: " + a.err.Error()
    }

    fields := fmt.Sprintf(
        "\n%s\n%s\n%s\n%s\n\n%s",
        a.username.View(),
        a.password.View(),
        a.host.View(),
        a.port.View(),
        "([shift]tab for prev/next | esc to quit)",
    )

    return err + "\n" + fields + "\n"
}
