package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) updateCurrDB() {
    m.currDB = m.DBTable.GetSelectedRow()[0]
    m.UpdateDBTablesTable()
    m.updateMainTable()
}

func (m *model) updateMainTable() {
    m.currTable = m.DBTablesTable.GetSelectedRow()[0]
    m.UpdateMainTable()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
    case dbConnectMsg:
        m.db = msg.db
        if (m.db != nil) {
            m.UpdateDBTable()
            m.updateCurrDB()
        }

    case tea.WindowSizeMsg:
        //m.mainTable.UpdateRenderedColums()
        //m.DBTablesTable.UpdateRenderedColums()
        //m.DBTable.UpdateRenderedColums()

    case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c":
			return m, tea.Quit

        case "1":
            m.selectedPane = DB
        case "2":
            m.selectedPane = DBTables
        case "3":
            m.selectedPane = Table
        }
	}

    switch m.selectedPane {
    case DB:
        m.DBTable, cmd = m.DBTable.Update(msg)

        if key, ok := msg.(tea.KeyMsg); ok {
            switch key.String() {
            case "enter", "j", "k":
                m.updateCurrDB()
            }
        }

    case DBTables:
        m.DBTablesTable, cmd = m.DBTablesTable.Update(msg)

        if key, ok := msg.(tea.KeyMsg); ok {
            switch key.String() {
            case "enter", "j", "k":
                m.updateMainTable()
            }
        }
    case Table:
        m.mainTable, cmd = m.mainTable.Update(msg)
    }

	return m, cmd
}
