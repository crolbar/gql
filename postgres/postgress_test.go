package postgres_test

import (
	"fmt"
	"gql/dbms"
	"gql/postgres"
	"gql/table"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
    // postgresql
    // postgres
    uri := "postgresql://crolbar:aoeu@localhost:5432/t?sslmode=disable"
    m := postgres.Model{Uri: uri}
    cmd := m.Open()
    m.SetDb(cmd().(dbms.DbConnectMsg).Db)

    testDatabases(m)
    testDBTables(m)
    testTable(m)
    testDescribe(m)
    testUser(t, m)
    testDelDB(t, m)
    testDelDBTable(t, m)
    testDelRow(t, m)
    testUpdateCell(t, m)
    testChangeDbTableName(t, m)
    testSendQuery(t, m)
}

func testDatabases(m postgres.Model) {
    cols, rows, err := m.GetDatabases("")
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(cols)
    fmt.Println(rows)
}

func testDBTables(m postgres.Model) {
    cols, rows, err := m.GetDBTables("t", "")
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(cols)
    fmt.Println(rows)
}

func testTable(m postgres.Model) {
    cols, rows, err := m.GetTable("", "t1", "")
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(cols)
    fmt.Println(rows)
}

func testDescribe(m postgres.Model) {
    cols, rows, err := m.GetDescribe("", "t2")
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(cols)
    fmt.Println(rows)
}

func testUser(t *testing.T, m postgres.Model) {
    user, err := m.GetUser()
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(user) 
    assert.Equal(t, "crolbar", user)
}

func testDelDB(t *testing.T, m postgres.Model) {
    _, err := m.Db.Exec("create database tmp")
    assert.Equal(t, err, nil)

    err = m.DeleteDB("tmp")

    assert.Equal(t, err, nil)
}

func testDelDBTable(t *testing.T, m postgres.Model) {
    _, err := m.Db.Exec("create database tmp")
    assert.Equal(t, err, nil)

    func (t *testing.T) {
        uri := "postgresql://crolbar:aoeu@localhost:5432/tmp?sslmode=disable"
        tmpm := postgres.Model{Uri: uri}
        tmpm.SetDb(tmpm.Open()().(dbms.DbConnectMsg).Db) // such beatiful syntax

        _, err := tmpm.Db.Exec("create table t1 (id int, name varchar(20))")
        assert.Equal(t, err, nil)

        assert.Equal(t, tmpm.DeleteDBTable("", "t1"), nil)
        assert.EqualError(t, tmpm.DeleteDBTable("", "t1"), "pq: table \"t1\" does not exist")
    }(t)

    _, err = m.Db.Exec("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'tmp'")
    assert.Equal(t, err, nil)
    assert.Equal(t, m.DeleteDB("tmp"), nil)
}

func testDelRow(t *testing.T, m postgres.Model) {
    _, err := m.Db.Exec("create database tmp")
    assert.Equal(t, err, nil)

    func (t *testing.T) {
        uri := "postgresql://crolbar:aoeu@localhost:5432/tmp?sslmode=disable"
        tmpm := postgres.Model{Uri: uri}
        tmpm.SetDb(tmpm.Open()().(dbms.DbConnectMsg).Db) // such beatiful syntax

        _, err := tmpm.Db.Exec("create table t1 (id int, name varchar(20))")
        assert.Equal(t, err, nil)


        { // deleting row ==================================
            _, err := tmpm.Db.Exec("insert into t1 values (2, 'name')")
            assert.Equal(t, err, nil)

            fmt.Println(tmpm.GetTable("", "t1", ""))

            assert.Equal(t,
            tmpm.DeleteRow("", "t1",
                table.Row{"2", "name"},
                []table.Column{
                    {Title: "id", Width: 0},
                    {Title: "name", Width: 0},
                },
            ),
            nil)
        }


        assert.Equal(t, tmpm.DeleteDBTable("", "t1"), nil)
        assert.EqualError(t, tmpm.DeleteDBTable("", "t1"), "pq: table \"t1\" does not exist")
    }(t)

    _, err = m.Db.Exec("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'tmp'")
    assert.Equal(t, err, nil)
    assert.Equal(t, m.DeleteDB("tmp"), nil)
}


func testUpdateCell(t *testing.T, m postgres.Model) {
    _, err := m.Db.Exec("create database tmp")
    assert.Equal(t, err, nil)

    func (t *testing.T) {
        uri := "postgresql://crolbar:aoeu@localhost:5432/tmp?sslmode=disable"
        tmpm := postgres.Model{Uri: uri}
        tmpm.SetDb(tmpm.Open()().(dbms.DbConnectMsg).Db) // such beatiful syntax

        _, err := tmpm.Db.Exec("create table t1 (id int, name varchar(20))")
        assert.Equal(t, err, nil)


        {
            _, err := tmpm.Db.Exec("insert into t1 values (2, 'name')")
            _, err = tmpm.Db.Exec("insert into t1 values (8, 'othername')")
            assert.Equal(t, err, nil)

            _, oldRows, _ := tmpm.GetTable("", "t1", "")
            fmt.Print("old rows: ")
            fmt.Println(oldRows)

            { // updating cell ==================================
                assert.Equal(t,
                    tmpm.UpdateCell(
                        "",
                        "t1",
                        table.Row{"2", "name"},
                        []table.Column{
                            {Title: "id", Width: 0},
                            {Title: "name", Width: 0},
                        },
                        0,
                        "3",
                    ),
                    nil,
                )

                _, updatedRows, _ := tmpm.GetTable("", "t1", "")
                fmt.Print("updated rows: ")
                fmt.Println(updatedRows)

                assert.Equal(t, []table.Row{{"8", "othername"}, {"3", "name"}}, updatedRows)
            }

            assert.Equal(t,
                tmpm.DeleteRow("", "t1",
                    table.Row{"2", "name"},
                    []table.Column{
                        {Title: "id", Width: 0},
                        {Title: "name", Width: 0},
                    },
                ),
                nil,
            )
        }


        assert.Equal(t, tmpm.DeleteDBTable("", "t1"), nil)
        assert.EqualError(t, tmpm.DeleteDBTable("", "t1"), "pq: table \"t1\" does not exist")
    }(t)

    _, err = m.Db.Exec("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'tmp'")
    assert.Equal(t, err, nil)
    assert.Equal(t, m.DeleteDB("tmp"), nil)
}


func testChangeDbTableName(t *testing.T, m postgres.Model) {
    _, err := m.Db.Exec("create database tmp")
    assert.Equal(t, err, nil)

    func (t *testing.T) {
        uri := "postgresql://crolbar:aoeu@localhost:5432/tmp?sslmode=disable"
        tmpm := postgres.Model{Uri: uri}
        tmpm.SetDb(tmpm.Open()().(dbms.DbConnectMsg).Db) // such beatiful syntax

        _, err := tmpm.Db.Exec("create table t1 (id int, name varchar(20))")
        assert.Equal(t, err, nil)

        {
            fmt.Print("old tables: ")
            fmt.Println(tmpm.GetDBTables("tmp", ""))

            assert.Equal(t, tmpm.ChangeDbTableName("", "t1", "table1"), nil)

            fmt.Print("new tables: ")
            fmt.Println(tmpm.GetDBTables("tmp", ""))

            assert.Equal(t, tmpm.ChangeDbTableName("", "table1", "t1"), nil)
        }


        assert.Equal(t, tmpm.DeleteDBTable("", "t1"), nil)
        assert.EqualError(t, tmpm.DeleteDBTable("", "t1"), "pq: table \"t1\" does not exist")
    }(t)

    _, err = m.Db.Exec("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'tmp'")
    assert.Equal(t, err, nil)
    assert.Equal(t, m.DeleteDB("tmp"), nil)
}

func testSendQuery(t *testing.T, m postgres.Model) {
    assert.Equal(t, m.SendQuery("select now()"), nil)
}
