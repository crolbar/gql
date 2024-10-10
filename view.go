package main

import "fmt"

func (m model) View() string {
    log := fmt.Sprintf(
        "Height: %d, yOff: %d, xOff: %d, cursor: %d, dbg: %s",
        m.table.Height,
        m.table.YOffset,
        m.table.XOffset,
        m.table.Cursor,
        m.table.Dbg,
    )

    return log + "\n" + baseStyle.Render(m.table.View()) + "\n"// + m.table.HelpView() + "\n" 
}

