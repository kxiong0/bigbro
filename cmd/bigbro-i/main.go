package main

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kxiong0/bigbro/internal/config"
)

type logMsg struct {
	timestamp time.Time
	line      string
}

// A command that waits for the activity on a channel.
func waitForActivity(sub chan string) tea.Cmd {
	return func() tea.Msg {
		return logMsg{timestamp: time.Now(), line: <-sub}
	}
}

type model struct {
	ready         bool
	viewport      viewport.Model
	content       string
	logInputChans []chan string
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		waitForActivity(m.logInputChans[0]), // wait for activity
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
			m.viewport.SetContent("111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111\n1111\n1111\n1111\n1111\n1111\n1111\n2222\n")
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
		}
	case logMsg:
		m.content = fmt.Sprintf("%s%s", m.content, msg.line)
		m.viewport.SetContent(m.content)
		return m, waitForActivity(m.logInputChans[0]) // wait for next event
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
	c := config.Config{}
	err := c.LoadConfigFile("config/default.json")
	if err != nil {
		log.Fatal(err)
	}

	scanners := c.GetInputScanners()
	bb := BigBro{}
	for _, scanner := range scanners {
		bb.AddInputScanner(scanner)
	}
	bb.Start()

	m := model{}
	m.logInputChans = []chan string{make(chan string)}
	p := tea.NewProgram(
		m,
		tea.WithMouseAllMotion(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
