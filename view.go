package main

import "fmt"

//"strings"
//
//"github.com/charmbracelet/lipgloss"

func (m model) View() string {
    //logg("\n=========START================\n" + m.table.View() + "\n=========================END============================\n\n\n")

    //log := fmt.Sprintf("y: %d\n", m.table.Cursor())

    log := fmt.Sprintf(
        "vpHeight: %d, y: %d, start: %d, end: %d, yOffset: %d\n",
        m.table.Viewport.Height / 2,
        m.table.Cursor(),
        m.table.Start,
        m.table.End,
        m.table.Viewport.YOffset,
    )

    return log + "\n" + baseStyle.Render(m.table.View()) + "\n"// + m.table.HelpView() + "\n" 
}

