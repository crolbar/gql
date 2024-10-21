package tabs

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

func (t Tabs) View() string {
    tabs          := make([]string, Describe)
    separator     := " | "
    inactiveColor := lipgloss.Color("240")

    for i := Main; i <= Describe; i++ {
        s := lipgloss.NewStyle()

        if t.selected != i {
            s = s.Foreground(inactiveColor)
        }

        tab := fmt.Sprintf("[%d]", i+1)
        tabs = append(tabs, s.Render(tab))

        if i != Describe {
            tabs = append(tabs, separator)
        }
    }

    join := lipgloss.JoinHorizontal(lipgloss.Left, tabs...)

    return lipgloss.NewStyle().
        Border(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color("240")).
        Render(join)
}
