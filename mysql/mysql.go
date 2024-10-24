package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"gql/table"

	_ "github.com/go-sql-driver/mysql"
)

func GetDatabases(
    db *sql.DB,
    whereClause string,
) ([]table.Column, []table.Row, error) {
    query := "show databases"

    if whereClause != "" {
        query = fmt.Sprintf("%s where Database = '%s'", query, whereClause)
    }

    rowsRes, err := db.Query(query)
	if err != nil {
        return nil, nil, err
	}

    var rows []table.Row

    width := 0

	for rowsRes.Next() {
		var dbName string

		if err := rowsRes.Scan(&dbName); err != nil {
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

func GetTables(
    db *sql.DB,
    dbName,
    whereClause string,
) ([]table.Column, []table.Row, error) {
    query := fmt.Sprintf("show tables from %s", dbName)

    if whereClause != "" {
        query = fmt.Sprintf("%s where Tables_in_%s = '%s'", query, dbName, whereClause)
    }

    rowsRes, err := db.Query(query)
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

func GetTable(
    db *sql.DB,
    currDB,
    selTable,
    whereClause string,
) ([]table.Column, []table.Row, error) {
    query := fmt.Sprintf("select * from %s.%s", currDB, selTable)

    if whereClause != "" {
        query = fmt.Sprintf("%s where %s", query, whereClause)
    }

    rowsRes, err := db.Query(query)
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

func GetDescribe(
    db *sql.DB,
    currDB,
    selTable string,
) ([]table.Column, []table.Row, error) {
    rowsRes, err := db.Query(fmt.Sprintf("describe %s.%s", currDB, selTable))
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

func GetUser(db *sql.DB) (string, error) {
    res, err := db.Query(fmt.Sprintf("select user()"))
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

func DeleteDB(db *sql.DB, dbName string) error {
    _, err := db.Query(
        fmt.Sprintf(
            "drop database %s",
            dbName,
        ),
    )
    return err
}

func DeleteDBTable(db *sql.DB, dbName, selTable string) error {
    _, err := db.Query(
        fmt.Sprintf(
            "drop table %s.%s",
            dbName,
            selTable,
        ),
    )
    return err
}

func DeleteRow(
    db *sql.DB,
    dbName,
    tableName string,
    row table.Row,
    cols []table.Column,
) error {
    _, err := db.Query(
        fmt.Sprintf(
            "delete from %s.%s where %s",
            dbName,
            tableName,
            buildWhereClause(row, cols),
        ),
    )

    return err
}

func UpdateCell(
    db *sql.DB,
    dbName,
    tableName string,
    row table.Row,
    cols []table.Column,
    selectedCol int,
    value string,
) error {
    _, err := db.Query(
        fmt.Sprintf(
            "update %s.%s set %s = '%s' where %s",
            dbName,
            tableName,
            cols[selectedCol].Title,
            value,
            buildWhereClause(row, cols),
        ),
    )

    return err
}

func ChangeDbTableName(
    db *sql.DB,
    dbName,
    tableName string,
    value string,
) error {
    _, err := db.Query(
        fmt.Sprintf(
            "alter table %s.%s rename to %s.%s",
            dbName,
            tableName,
            dbName,
            value,
        ),
    )

    return err
}

func buildWhereClause(
    row table.Row,
    cols []table.Column,
) string {
    var sb strings.Builder

    for i := 0; i < len(cols); i++ {
        if row[i] == "NULL" {
            sb.WriteString(fmt.Sprintf("`%s` IS NULL", cols[i].Title))
        } else {
            col := strings.ReplaceAll(row[i], "'", "\\'") // replace "'" with "\'"

            sb.WriteString(fmt.Sprintf("`%s` = '%s'", cols[i].Title, col))
        }

        if i != len(cols) - 1 {
            sb.WriteString(" and ")
        }
    }

    return sb.String()
}

func SendQuery(db *sql.DB, query string) error {
    _, err := db.Query(query)
    return err
}


func getTableFromQueryRes(res *sql.Rows) ([]table.Column, []table.Row, error) {
	columnsRes, err := res.Columns()
	if err != nil {
        return nil, nil, err
	}

	values        := make([]interface{}, len(columnsRes))
	valuePointers := make([]interface{}, len(columnsRes))

    var rows []table.Row
    for res.Next() {
        for i := range columnsRes {
			valuePointers[i] = &values[i]
		}

        var currRow table.Row

        if err := res.Scan(valuePointers...); err != nil {
            return nil, nil, err
		}

        for i := range columnsRes {
			switch val := values[i].(type) {
			case nil:
                currRow = append(currRow, "NULL")
			case []byte:
                text := string(val)
                text = strings.ReplaceAll(text, "\n", "\\n")
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
