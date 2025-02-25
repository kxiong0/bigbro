package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type logMsg struct {
	timestamp time.Time
	line      string
}

// Simulate a process that sends events at an irregular interval in real time.
// In this case, we'll send events on the channel at a random interval between
// 100 to 1000 milliseconds. As a command, Bubble Tea will run this
// asynchronously.
func listenForActivity(subs []chan string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("bash", "-c", "kubectl logs -l app=kindnet  -n kube-system -f --since=1s")
		out, _ := cmd.StdoutPipe()

		// Run the command.
		if err := cmd.Start(); err != nil {
			return nil
		}

		// Read command output as it arrives.
		buf := bufio.NewReader(out)
		for {
			line, err := buf.ReadString('\n')
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return nil
			}
			// Send output to our program.
			subs[0] <- fmt.Sprintf("Current time: %s - echo: %s", time.Now(), line)
		}
	}
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
		listenForActivity(m.logInputChans),  // generate activity
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
