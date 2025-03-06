package scanner

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os/exec"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

var colorMap = map[string]color.Attribute{
	"black":   color.FgBlack,
	"white":   color.FgWhite,
	"red":     color.FgRed,
	"green":   color.FgGreen,
	"yellow":  color.FgYellow,
	"blue":    color.FgBlue,
	"magenta": color.FgMagenta,
	"cyan":    color.FgCyan,
}

type InputScanner interface {
	SetName(string)
	Init()
	Start() error
	Close()
	GetOutputChan() chan string
}

type BaseInputScanner struct {
	Name        string `json:"name"`
	OutputColor string `json:"color"`
	OutputChan  chan string
}

func (bis *BaseInputScanner) Init() {
	bis.OutputChan = make(chan string, 10)
}

func (bis *BaseInputScanner) Close() {
	close(bis.OutputChan)
}

func (bis *BaseInputScanner) SetName(name string) {
	bis.Name = name
}

func (bis *BaseInputScanner) SetOutputColor(color string) {
	bis.OutputColor = color
}

func (bis *BaseInputScanner) GetCMD() (string, error) {
	return "", nil
}

func (bis *BaseInputScanner) Start() error {
	bis.OutputChan <- "BaseScanner Start"
	return nil
}

func (bis *BaseInputScanner) GetOutputChan() chan string {
	return bis.OutputChan
}

type CmdInputScanner struct {
	Command string `json:"command"`

	BaseInputScanner
}

func (cis *CmdInputScanner) SetCmd(command string) {
	cis.Command = command
}

func (cis *CmdInputScanner) GetCMD() (string, error) {
	return cis.Command, nil
}

func (cis *CmdInputScanner) Start() error {
	cmd, err := cis.GetCMD()
	if err != nil {
		return err
	}

	execCmd := exec.Command("bash", "-c", cmd)
	out, _ := execCmd.StdoutPipe()
	if err := execCmd.Start(); err != nil {
		return nil
	}

	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("13"))

		log.Println(style.Render(scanner.Text()))

		cis.OutputChan <- scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Println("Error scanning output:", err)
		}
	}
	return nil
}

type K8sInputScanner struct {
	UseK8sTimestamp bool
	Pod             Pod

	CmdInputScanner
}

type Pod struct {
	Name        string
	Namespace   string
	PodSelector map[string]string
	Container   string
}

func (kis *K8sInputScanner) GetCMD() (string, error) {
	_, err := exec.LookPath("kubectl")
	if err != nil {
		return "", err
	}

	var cmd string
	if kis.Pod.Name != "" {
		cmd = fmt.Sprintf("kubectl logs %s", kis.Pod.Name)
	} else if kis.Pod.PodSelector != nil {
		cmd = "kubectl logs"
		for key, value := range kis.Pod.PodSelector {
			cmd = fmt.Sprintf("%s -l %s=%s", cmd, key, value)
		}
	} else {
		return "", errors.New("must provide one of: podName, podSelector")
	}

	cmd = fmt.Sprintf("%s -n %s --timestamps=true --since=1s -f", cmd, kis.Pod.Namespace)
	return cmd, nil
}
