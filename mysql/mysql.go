package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"gql/table"

	_ "github.com/go-sql-driver/mysql"
)

func GetDatabases(db *sql.DB) ([]table.Column, []table.Row) {
    rowsRes, err := db.Query("show databases;")
	if err != nil {
		log.Fatal(err)
	}

    var rows []table.Row

    width := 0

	for rowsRes.Next() {
		var dbName string

		if err := rowsRes.Scan(&dbName); err != nil {
			log.Fatal(err)
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
            Width: min(width, 20),
        }, 
    }

    return cols, rows
}

func GetTables(db *sql.DB, dbName string) ([]table.Column, []table.Row) {
    rowsRes, err := db.Query(fmt.Sprintf("show tables from %s;", dbName))
	if err != nil {
		log.Fatal(err)
	}

    var rows []table.Row
	for rowsRes.Next() {
		var tableName string

		if err := rowsRes.Scan(&tableName); err != nil {
			log.Fatal(err)
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

    return cols, rows
}

func GetTable(db *sql.DB, currDB, selTable string) ([]table.Column, []table.Row) {
    rowsRes, err := db.Query(fmt.Sprintf("select * from %s.%s", currDB, selTable))
	if err != nil {
		log.Fatal(err)
	}

	columnsRes, err := rowsRes.Columns()
	if err != nil {
		log.Fatal(err)
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
			log.Fatal(err)
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


    return columns, rows
}

func GetDescribe(db *sql.DB, currDB, selTable string) ([]table.Column, []table.Row) {
    rowsRes, err := db.Query(fmt.Sprintf("describe %s.%s", currDB, selTable))
	if err != nil {
		log.Fatal(err)
	}

	columnsRes, err := rowsRes.Columns()
	if err != nil {
		log.Fatal(err)
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
			log.Fatal(err)
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

    return columns, rows
}

func GetUser(db *sql.DB) string {
    res, err := db.Query(fmt.Sprintf("select user()"))
	if err != nil {
		log.Fatal(err)
	}

    var user string
	for res.Next() {
		if err := res.Scan(&user); err != nil {
			log.Fatal(err)
		}
	}

    return user
}

func DeleteDB(db *sql.DB, dbName string) error {
    _, err := db.Query(fmt.Sprintf("drop database %s", dbName))
    return err
}

func DeleteRow(
    db *sql.DB,
    dbName,
    tableName string,
    row table.Row,
    cols []table.Column,
) error {
    var sb strings.Builder

    for i := 0; i < len(cols); i++ {
        if row[i] == "NULL" {
            sb.WriteString(fmt.Sprintf("`%s` IS NULL", cols[i].Title))
        } else {
            col := strings.ReplaceAll(row[i], "\\", "\\\\") // replace "\" with "\\"
            col = strings.ReplaceAll(col, "\\\\n", "\\n")   // replace "\\n" with "\n"
            col = strings.ReplaceAll(col, "'", "\\'")       // replace "'" with "\'"

            sb.WriteString(fmt.Sprintf("`%s` = '%s' COLLATE utf8mb4_general_ci", cols[i].Title, col))
        }

        if i != len(cols) - 1 {
            sb.WriteString(" and ")
        }
    }

    _, err := db.Query(
        fmt.Sprintf(
            "delete from %s.%s where %s",
            dbName,
            tableName,
            sb.String(),
        ),
    )
    if err != nil {
        log.Fatal(err)
    }

    return err
}

func getTableFromQueryRes(res *sql.Rows) ([]table.Column, []table.Row) {
	columnsRes, err := res.Columns()
	if err != nil {
		log.Fatal(err)
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
			log.Fatal(err)
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

    return columns, rows
}
