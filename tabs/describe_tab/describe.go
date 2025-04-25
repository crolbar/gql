package describe_tab

import (
	"gql/table"

	tea "github.com/charmbracelet/bubbletea"
)

type DescribeTab struct {
	Table table.Table
}

func New() DescribeTab {
	t := table.New(nil, nil, 100, 100)
	t.Focus()
	return DescribeTab{
		Table: t,
	}
}

func (t DescribeTab) Update(msg tea.Msg) (DescribeTab, tea.Cmd) {
	var cmd tea.Cmd
	t.Table, cmd = t.Table.Update(msg)
	return t, cmd
}

func (t DescribeTab) View() string {
	return t.Table.View()
}

func (t *DescribeTab) OnWindowResize(msg tea.WindowSizeMsg, isConnected bool) {
	width := msg.Width
	height := perc(80, msg.Height)

	t.Table.SetMaxWidth(width)
	t.Table.SetMaxHeight(height)

	t.Table.UpdateOffset()
}

func perc(per, num int) int {
	return int(float32(num) * (float32(per) / float32(100)))
}
