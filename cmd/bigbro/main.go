package main

import (
	"log"
	"os"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"

	"github.com/kxiong0/bigbro/internal/log_collector"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

// A command that waits for the activity on a channel.
func waitForActivity(sub chan log_collector.LogMsg) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

type model struct {
	bb *BigBro // Scanner manager

	table table.Model
	rows  []table.Row

	// Window dimensions
	totalWidth  int
	totalHeight int

	// Table dimensions
	horizontalMargin int
	verticalMargin   int
}

func (m model) rowStyleFunc(input table.RowStyleFuncInput) lipgloss.Style {
	scannerIdx := 0
	switch value := input.Row.Data["source"].(type) {
	case string:
		var err error
		scannerIdx, err = strconv.Atoi(value)
		if err != nil {
			return lipgloss.NewStyle().
				Foreground(
					lipgloss.Color("123"),
				)
		}
	case int:
		scannerIdx = value
	default:
		return lipgloss.NewStyle().
			Foreground(
				lipgloss.Color("123"),
			)
	}

	return lipgloss.NewStyle().
		Foreground(
			lipgloss.Color(m.bb.inputScanners[scannerIdx].GetColor()),
		)
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		waitForActivity(m.bb.LogOutputChan), // wait for activity
	)
}

func (m model) View() string {
	return m.table.View()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// if m.table.Focused() {
			// 	m.table.Blur()
			// } else {
			// 	m.table.Focus()
			// }
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			// return m, tea.Batch(
			// 	tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			// )
		}
	case tea.WindowSizeMsg:
		m.totalWidth = msg.Width
		m.totalHeight = msg.Height
		m.recalculateTable()
	case log_collector.LogMsg:
		log.Println("Add line to rows")
		row := table.NewRow(table.RowData{
			"time":   msg.Timestamp.Format(time.RFC3339),
			"source": msg.ScannerIdx,
			"line":   msg.Line,
		})
		m.rows = append(m.rows, row)
		m.table = m.table.WithRows(m.rows)

		cmds = append(cmds, waitForActivity(m.bb.LogOutputChan))
	}
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
func (m *model) recalculateTable() {
	m.table = m.table.
		WithTargetWidth(m.calculateWidth()).
		WithMinimumHeight(m.calculateHeight())
}

func (m model) calculateWidth() int {
	return m.totalWidth - m.horizontalMargin - 4
}

func (m model) calculateHeight() int {
	return m.totalHeight - m.verticalMargin - 2
}

func main() {
	// log to file
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Set the log output to the file
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	m := model{}
	bb := &BigBro{ConfigFilePath: "config/default.json"}
	err = bb.Init()
	if err != nil {
		log.Fatal(err)
	}
	m.bb = bb

	// Init table
	t := table.New(
		[]table.Column{
			table.NewColumn("time", "Time", 20),
			table.NewFlexColumn("source", "Source", 2),
			table.NewFlexColumn("line", "Line", 10),
		},
	).
		BorderRounded().
		WithBaseStyle(baseStyle).
		WithRowStyleFunc(m.rowStyleFunc)

	m.table = t
	p := tea.NewProgram(
		m,
		tea.WithMouseAllMotion(),
	)

	bb.Start()
	defer bb.Stop()

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
