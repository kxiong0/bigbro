package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/kxiong0/bigbro/internal/log_collector"
)

// A command that waits for the activity on a channel.
func waitForActivity(sub chan log_collector.LogMsg) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

type model struct {
	ready    bool
	viewport viewport.Model
	content  string // displayed logs
	bb       *BigBro
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		waitForActivity(m.bb.LogOutputChan), // wait for activity
	)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s",
		m.viewport.View(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height)
			// Whether or not to respond to the mouse. The mouse must be enabled in
			// Bubble Tea for this to work. For details, see the Bubble Tea docs.
			m.viewport.MouseWheelEnabled = true
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
		}
	case log_collector.LogMsg:
		logMsgStr := fmt.Sprintf("[%s] [%d] %s", msg.Timestamp.Format(time.RFC3339Nano), msg.ScannerIdx, msg.Line)
		style := lipgloss.NewStyle().
			Foreground(
				lipgloss.Color(m.bb.inputScanners[msg.ScannerIdx].GetColor()),
			)
		log.Println(m.bb.inputScanners[msg.ScannerIdx].GetColor())
		log.Println(m.bb.inputScanners[msg.ScannerIdx])
		m.content = fmt.Sprintf("%s\n%s", m.content, style.Render((logMsgStr)))
		m.viewport.SetContent(m.content)
		return m, waitForActivity(m.bb.LogOutputChan) // wait for next event
	default:
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
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

	bb := &BigBro{ConfigFilePath: "config/default.json"}
	err = bb.Init()
	if err != nil {
		log.Fatal(err)
	}

	m := model{}
	m.bb = bb
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
