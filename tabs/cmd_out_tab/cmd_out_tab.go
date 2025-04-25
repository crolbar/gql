package cmd_out_tab

import (
	"database/sql"
	"fmt"
	"gql/table"
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type CmdOutTab struct {
	Table table.Table
}

func New() CmdOutTab {
	t := table.New(nil, nil, 100, 100)
	t.Focus()
	return CmdOutTab{
		Table: t,
	}
}

func (t CmdOutTab) Update(msg tea.Msg) (CmdOutTab, tea.Cmd) {
	var cmd tea.Cmd
	t.Table, cmd = t.Table.Update(msg)
	return t, cmd
}

func (t *CmdOutTab) UpdateTable(rowsRes *sql.Rows) {
	columnsRes, err := rowsRes.Columns()
	if err != nil {
		return
	}

	values := make([]interface{}, len(columnsRes))
	valuePointers := make([]interface{}, len(columnsRes))

	var rows []table.Row

	for rowsRes.Next() {
		for i := range columnsRes {
			valuePointers[i] = &values[i]
		}

		var currRow table.Row

		if err := rowsRes.Scan(valuePointers...); err != nil {
			return
		}

		for i := range columnsRes {
			switch val := values[i].(type) {
			case nil:
				currRow = append(currRow, "NULL")
			case string:
				currRow = append(currRow, strings.ReplaceAll(val, "\n", "\\n"))
			case time.Time:
				currRow = append(currRow, val.Format("2006-01-02 15:04:05.999999-07"))
			case bool:
				if val {
					currRow = append(currRow, "true")
				} else {
					currRow = append(currRow, "false")
				}
			case []byte:
				text := string(val)
				text = strings.ReplaceAll(text, "\\", "\\\\") // replace "\" with "\\"
				text = strings.ReplaceAll(text, "\n", "\\n")  // replace new lines with "\n"
				currRow = append(currRow, text)
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				currRow = append(currRow, fmt.Sprintf("%d", val))
			case float32, float64:
				currRow = append(currRow, fmt.Sprintf("%f", val))
			default:
				log.Fatalf("Found type that's not supported, val: %v, type: %T", val, val)
			}
		}

		rows = append(rows, currRow)
	}

	columns := make([]table.Column, 0, len(columnsRes))
	for _, col := range columnsRes {
		columns = append(columns, table.Column{
			Title: col,
			Width: 10,
		})
	}

	t.Table.SetRows(rows)
	t.Table.SetColumns(columns)
}

func (t CmdOutTab) View() string {
	return t.Table.View()
}

func (t *CmdOutTab) OnWindowResize(msg tea.WindowSizeMsg, isConnected bool) {
	width := msg.Width
	height := int(0.8 * float64(msg.Height))

	t.Table.SetMaxWidth(width)
	t.Table.SetMaxHeight(height)

	t.Table.UpdateOffset()
}
