package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if (m.requiresAuth()) {
        return m.authUpdate(msg)
    }

    return m.mainUpdate(msg)
}

func (m model) authUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
    a, cmd, uri := m.auth.Update(msg)
    m.auth = a;

    if (uri != "") {
        m.uri = uri

        return m, m.OpenMysql
    }

    return m, cmd
}

func (m model) mainUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
    case dbConnectMsg:
        m.onDBConnect(msg.db)

    case tea.WindowSizeMsg:
        //m.mainTable.UpdateRenderedColums()
        //m.DBTablesTable.UpdateRenderedColums()
        //m.DBTable.UpdateRenderedColums()

    case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c":
			return m, tea.Quit

        case "s":
            m.changeCreds()
            return m, nil
        }
	}

    switch m.selectedPane {
    case DB:
        m.DBTable, cmd = m.DBTable.Update(msg)

        if key, ok := msg.(tea.KeyMsg); ok {
            switch key.String() {
            case "enter":
                m.selectDBTablespane()
                fallthrough
            case "j", "k":
                m.updateCurrDB()
            }
        }

    case DBTables:
        m.DBTablesTable, cmd = m.DBTablesTable.Update(msg)

        if key, ok := msg.(tea.KeyMsg); ok {
            switch key.String() {
            case "esc":
                m.selectDBpane()
            case "enter":
                m.selectMainpane()
                fallthrough
            case "j", "k":
                m.updateMainTable()
            }
        }
    case Main:
        m.mainTable, cmd = m.mainTable.Update(msg)

        if key, ok := msg.(tea.KeyMsg); ok {
            switch key.String() {
            case "esc":
                m.selectDBTablespane()
            }
        }
    }

	return m, cmd
}
