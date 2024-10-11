package main

import (
	"fmt"
	"gql/util"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.table.Width = msg.Width
        m.table.UpdateRenderedColums()

    case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c":
			return m, tea.Quit

        case " ":
            util.Logg(fmt.Sprintf("%v", m.table.GetSelectedRows()))
        }

	}

    m.table.Dbg = ""
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}
