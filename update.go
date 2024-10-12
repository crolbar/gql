package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) updateCurrDB() {
    m.currDB = m.DBTable.GetSelectedRow()[0]
    m.CreateDBTablesTable()
    m.updateMainTable()
}

func (m *model) updateMainTable() {
    m.currTable = m.DBTablesTable.GetSelectedRow()[0]
    m.CreateMainTable()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
    case dbConnectMsg:
        m.db = msg.db
        if (m.db != nil) {
            m.CreateDBTable()
            m.updateCurrDB()
            m.updateMainTable()
        }

    case tea.WindowSizeMsg:
        if (m.mainTable != nil) {
            //m.mainTable.Width = msg.Width
            m.mainTable.UpdateRenderedColums()
        }

        if (m.DBTablesTable != nil) {
            m.DBTablesTable.UpdateRenderedColums()
        }

        if (m.DBTable != nil) {
            m.DBTable.UpdateRenderedColums()
        }

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
        if (m.DBTable != nil) {
            *m.DBTable, cmd = m.DBTable.Update(msg)
        }

        if key, ok := msg.(tea.KeyMsg); ok {
            switch key.String() {
            case "enter", "j", "k":
                m.updateCurrDB()
            }
        }

    case DBTables:
        if (m.DBTablesTable != nil) {
            *m.DBTablesTable, cmd = m.DBTablesTable.Update(msg)
        }

        if key, ok := msg.(tea.KeyMsg); ok {
            switch key.String() {
            case "enter", "j", "k":
                m.updateMainTable()
            }
        }
    case Table:
        if (m.mainTable != nil) {
            *m.mainTable, cmd = m.mainTable.Update(msg)
        }
    }

	return m, cmd
}
