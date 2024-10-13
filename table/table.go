package table

import (
	"fmt"
	"log"
	"math"
)

type Cursor struct {
    X int
    Y int
}

type Table struct {
    cols []Column
    rows []Row

	Cursor Cursor

    Height int
    Width  int

    KeyMap KeyMap

    XOffset int
    YOffset int

    renderedColumns int

    Dbg string

    columnSelect   bool
    rowSelect      bool
    selectionStart int
}


type Row []string

type Column struct {
	Title string
	Width int
}

func New(cols []Column, rows []Row, height int, width int) Table {
    return Table {
        cols:   cols,
        rows:   rows,
        Height: height - 2,
        Width:  width,

        XOffset: 0,
        YOffset: 0,

        renderedColumns: len(cols),

        KeyMap: DefaultKeyMap(),
        Cursor: Cursor{0, 0},

        columnSelect: false,
        rowSelect:    false,
        selectionStart:  -1,
    }
}

func (t *Table) SetColumns(cols []Column) {
    t.cols = cols
    t.UpdateOffset()
}

func (t *Table) SetRows(rows []Row) {
    t.rows = rows
    t.UpdateOffset()
}

func (t *Table) GetCursor() Cursor {
    return t.Cursor
}

func (t *Table) GetSelectedRow() Row {
    return t.rows[t.Cursor.Y]
}

// returns nil if we aren't in selection
func (t *Table) GetSelectedRows() []Row {
    if (t.rowSelect && t.selectionStart >= 0) {
        start := min(t.selectionStart, t.Cursor.Y)
        end := max(t.selectionStart, t.Cursor.Y)
        return t.rows[start:end+1]
    }
    return nil
}

// returns nil if we aren't in selection
func (t *Table) GetSelectedColumns() []Column {
    if (t.columnSelect && t.selectionStart >= 0) {
        start := min(t.selectionStart, t.Cursor.X)
        end := max(t.selectionStart, t.Cursor.X)
        return t.cols[start:end+1]
    }
    return nil
}

// Must be called when the width of the terminal changes 
// or there is an update to cursor.x
func (t *Table) UpdateRenderedColums() {
    width := 0
    for i := t.XOffset; i < len(t.cols); i++ {
        currColWidth := t.cols[i].Width + 1
        width += currColWidth

        if currColWidth > t.Width {
            log.Fatal(fmt.Sprintf("Column %d's width is bigger than the table's width %d", i, t.Width))
        }

        if width > t.Width {
            t.renderedColumns = i - t.XOffset
            break;
        }
    }
}

// Must be called when the width of the terminal changes 
// or there is an update to the cursor
func (t *Table) UpdateOffset()  {
    lines_till_sow := t.Cursor.Y - t.YOffset
    lines_till_eow := ((t.Height / 2) - 1) - lines_till_sow

    if (lines_till_eow < 0) {
        t.YOffset += int(math.Abs(float64(lines_till_eow)))
    } 

    if (lines_till_sow < 0) {
        t.YOffset -= int(math.Abs(float64(lines_till_sow)))
    }

    cols_till_sow := t.Cursor.X - t.XOffset
    cols_till_eow := (t.renderedColumns - 1) - cols_till_sow

    //t.Dbg = fmt.Sprintf("cols_t_sow: %d, cols_t_eow: %d, rc: %d", cols_till_sow, cols_till_eow, t.renderedColumns)

    for cols_till_eow < 0 || cols_till_sow < 0 {
        if (cols_till_eow < 0) {
            t.XOffset += int(math.Abs(float64(cols_till_eow)))
        }

        if (cols_till_sow < 0) {
            t.XOffset -= int(math.Abs(float64(cols_till_sow)))
        }

        t.UpdateRenderedColums()

        cols_till_sow = t.Cursor.X - t.XOffset
        cols_till_eow = (t.renderedColumns - 1) - cols_till_sow
    }

    t.UpdateRenderedColums()
}
