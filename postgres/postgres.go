package postgres

import (
	"database/sql"
	"fmt"
	"gql/table"
	"log"
	"strings"
	"time"

	"gql/dbms"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/lib/pq"
)

type Model struct { Db *sql.DB }

func (m Model) HasDb() bool {
    return m.Db != nil
}

func (m *Model) SetDb(db *sql.DB) {
    m.Db = db
}

func (m Model) Open(uri string) tea.Cmd {
    db, err := sql.Open("postgres", uri)

    err = db.Ping()
	if err != nil {
        fmt.Println(err)
		return func() tea.Msg {
            return dbms.DbConnectMsg{}
        }
	}
    return func() tea.Msg {
        return dbms.DbConnectMsg{Db: db}
    }
}

func (m Model) GetDatabases(
    whereClause string,
) ([]table.Column, []table.Row, error) {
    query := "select datname from pg_database"

    if whereClause != "" {
        query = fmt.Sprintf("%s where datname = '%s'", query, whereClause)
    }

    res, err := m.Db.Query(query)
    if err != nil {
        return nil, nil, err
    }

    var rows []table.Row

    width := 0

    for res.Next() {
        var dbName string

        if err := res.Scan(&dbName); err != nil {
            return nil, nil, err
        }

        rows = append(rows, []string{dbName})

        l := len(dbName)
        if (l > width) {
            width = l
        }
    }

    title := fmt.Sprintf("Databases")

    cols := []table.Column {
        {
            Title: title,
            Width: max(min(width, 20), 10),
        }, 
    }

    return cols, rows, nil
}


func (m Model) GetDBTables(
    dbName, // useless
    whereClause string,
) ([]table.Column, []table.Row, error) {
    query := fmt.Sprintf("select tablename from pg_tables where schemaname = 'public'")

    if whereClause != "" {
        query = fmt.Sprintf("%s and tablename = '%s'", query, whereClause)
    }

    rowsRes, err := m.Db.Query(query)
	if err != nil {
        return nil, nil, err
	}

    var rows []table.Row
	for rowsRes.Next() {
		var tableName string

		if err := rowsRes.Scan(&tableName); err != nil {
            return nil, nil, err
		}

		rows = append(rows, []string{tableName})
	}

    title := fmt.Sprintf("tables in %s", dbName)

    cols := []table.Column {
        {
            Title: title,
            Width: max(len(title), 20),
        }, 
    }

    return cols, rows, nil
}


func (m Model) GetTable(
    currDB, // useless
    selTable,
    whereClause string,
) ([]table.Column, []table.Row, error) {
    query := fmt.Sprintf("select * from public.%s", selTable)

    if whereClause != "" {
        query = fmt.Sprintf("%s where %s", query, whereClause)
    }

    rowsRes, err := m.Db.Query(query)
	if err != nil {
        return nil, nil, err
	}

	columnsRes, err := rowsRes.Columns()
	if err != nil {
        return nil, nil, err
	}

	values        := make([]interface{}, len(columnsRes))
	valuePointers := make([]interface{}, len(columnsRes))

    var rows []table.Row

    for rowsRes.Next() {
        for i := range columnsRes {
			valuePointers[i] = &values[i]
		}

        var currRow table.Row

        if err := rowsRes.Scan(valuePointers...); err != nil {
            return nil, nil, err
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
                text = strings.ReplaceAll(text, "\n", "\\n") // replace new lines with "\n"
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
        columns = append(columns, table.Column {
            Title: col,
            Width: 10,
        })
    }


    return columns, rows, nil
}

func (m Model) GetDescribe(
    currDB,
    selTable string,
) ([]table.Column, []table.Row, error) {
    rowsRes, err := m.Db.Query(fmt.Sprintf(
        `SELECT column_name,
            data_type,
            character_maximum_length,
            column_default,
            is_nullable
        FROM information_schema.columns 
        WHERE table_name = '%s'`,
    selTable))

    if err != nil {
        return nil, nil, err
	}

	columnsRes, err := rowsRes.Columns()
	if err != nil {
        return nil, nil, err
	}

	values        := make([]interface{}, len(columnsRes))
	valuePointers := make([]interface{}, len(columnsRes))

    var rows []table.Row

    maxColWidth := make([]int, len(columnsRes))

    for rowsRes.Next() {
        for i := range columnsRes {
			valuePointers[i] = &values[i]
		}

        var currRow table.Row

        if err := rowsRes.Scan(valuePointers...); err != nil {
            return nil, nil, err
		}

        var text string

        for i := range columnsRes {
			switch val := values[i].(type) {
			case nil:
                text = "NULL"
            case string:
                text = strings.ReplaceAll(val, "\n", "\\n")
			case []byte:
                text = string(val)
                text = strings.ReplaceAll(text, "\n", "\\n")
            case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
                text = fmt.Sprintf("%d", val)
            case float32, float64:
                text = fmt.Sprintf("%f", val)
			default:
                log.Fatalf("Found type that's not supported, val: %v, type: %T", val, val)
			}

            maxColWidth[i] = max(maxColWidth[i], len(text))

            currRow = append(currRow, text)
        }

        rows = append(rows, currRow)
    }

    columns := make([]table.Column, 0, len(columnsRes))
    for i, col := range columnsRes {
        columns = append(columns, table.Column {
            Title: col,
            Width: max(maxColWidth[i], len(col)),
        })
    }

    return columns, rows, nil
}


func (m Model) GetUser() (string, error) {
    res, err := m.Db.Query(fmt.Sprintf("select user"))
	if err != nil {
        return "", err
	}

    var user string
	for res.Next() {
		if err := res.Scan(&user); err != nil {
            return "", err
		}
	}

    return user, nil
}

func (m Model) DeleteDB(dbName string) error {
    _, err := m.Db.Query(
        fmt.Sprintf(
            "drop database %s",
            dbName,
        ),
    )
    return err
}

func (m Model) DeleteDBTable(dbName, selTable string) error {
    _, err := m.Db.Query(
        fmt.Sprintf(
            "drop table %s",
            selTable,
        ),
    )
    return err
}

func (m Model) DeleteRow(
    dbName,
    tableName string,
    row table.Row,
    cols []table.Column,
) error {
    _, err := m.Db.Query(
        fmt.Sprintf(
            "delete from %s where %s",
            tableName,
            buildWhereClause(row, cols),
        ),
    )

    return err
}

func (m Model) UpdateCell(
    dbName,
    tableName string,
    row table.Row,
    cols []table.Column,
    selectedCol int,
    value string,
) error {
    _, err := m.Db.Query(
        fmt.Sprintf(
            "update %s set %s = '%s' where %s",
            tableName,
            cols[selectedCol].Title,
            value,
            buildWhereClause(row, cols),
        ),
    )

    return err
}

func (m Model) ChangeDbTableName(
    dbName,
    tableName string,
    value string,
) error {
    _, err := m.Db.Query(
        fmt.Sprintf(
            "alter table %s rename to %s",
            tableName,
            value,
        ),
    )

    return err
}

func (m Model) SendQuery(query string) error {
    _, err := m.Db.Query(query)
    return err
}

func buildWhereClause(
    row table.Row,
    cols []table.Column,
) string {
    var sb strings.Builder

    for i := 0; i < len(cols); i++ {
        if row[i] == "NULL" {
            sb.WriteString(fmt.Sprintf("%s IS NULL", cols[i].Title))
        } else {
            col := strings.ReplaceAll(row[i], "'", "\\'") // replace "'" with "\'"

            sb.WriteString(fmt.Sprintf("%s = '%s'", cols[i].Title, col))
        }

        if i != len(cols) - 1 {
            sb.WriteString(" and ")
        }
    }

    return sb.String()
}
