package tests

import (
	"gql/data"
	"gql/table"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test00(t *testing.T) {
	table := table.New(data.ColumnsBig, data.RowsBig, 32+1+2, 148)
	table.UpdateRenderedColums()

	f, _ := os.ReadFile("./dump-00")
	assert.Equal(t, string(f), table.View()+"\n")
}

func Test01(t *testing.T) {
	table := table.New(data.ColumnsBig, data.RowsBig, 32+1+2, 168)
	table.UpdateRenderedColums()

	f, _ := os.ReadFile("./dump-01")
	assert.Equal(t, string(f), table.View()+"\n")
}

func Test02(t *testing.T) {
	table := table.New(data.ColumnsBig, data.RowsBig, 32+1+2, 168)
	table.UpdateRenderedColums()

	table.GotoBottom()
	table.MoveRight(1)
	table.MoveRight(1)
	table.MoveRight(1)
	table.MoveRight(1)
	table.MoveRight(1)

	table.UpdateRenderedColums()

	//util.Logg(fmt.Sprintf("%v\n%s", table.GetCursor(), table.View()))

	f, _ := os.ReadFile("./dump-02")
	assert.Equal(t, string(f), table.View()+"\n")
}

func Test03(t *testing.T) {
	table := table.New(data.ColumnsBig, data.RowsBig, 32+1+2, 148)
	table.UpdateRenderedColums()

	table.MoveDown(20)
	table.MoveRight(1)
	table.MoveRight(1)
	table.MoveRight(1)
	table.MoveRight(1)
	table.MoveRight(1)
	table.MoveRight(1)
	table.MoveRight(1)

	table.UpdateRenderedColums()

	//util.Logg(fmt.Sprintf("%v\n%s", table.GetCursor(), table.View()))

	f, _ := os.ReadFile("./dump-03")
	assert.Equal(t, string(f), table.View()+"\n")
}

func TestOffset(t *testing.T) {
	table := table.New(data.ColumnsBig, data.RowsBig, 32+1+2, 148)
	table.UpdateRenderedColums()

	table.MoveRight(1)
	table.MoveRight(1)

	assert.Equal(t, table.GetXOffset(), 1)

	table.MoveDown(70)

	assert.Equal(t, table.GetXOffset(), 1)
	assert.Equal(t, table.GetYOffset(), 56)

	table.MoveRight(1)
	assert.Equal(t, table.GetXOffset(), 1)

	table.MoveRight(1)
	assert.Equal(t, table.GetXOffset(), 2)

	table.MoveRight(1)
	assert.Equal(t, table.GetXOffset(), 3)
	assert.Equal(t, table.GetYOffset(), 56)

	table.GotoTop()
	table.MoveLeft(1)
	table.MoveLeft(1)
	table.MoveLeft(1)
	table.MoveLeft(1)
	table.MoveLeft(1)

	assert.Equal(t, table.GetXOffset(), 0)
	assert.Equal(t, table.GetYOffset(), 0)

	table.GotoBottom()
	assert.Equal(t, table.GetYOffset(), 85)
}
