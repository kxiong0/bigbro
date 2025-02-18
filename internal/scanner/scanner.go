package scanner

import (
	"bufio"
	"errors"
	"log"
	"os/exec"

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
	SetOutputColor(string)
	Start() error
}

type BaseInputScanner struct {
	Name        string `json:"name"`
	OutputColor string `json:"color"`
}

func (bis *BaseInputScanner) SetName(name string) {
	bis.Name = name
}

func (bis *BaseInputScanner) SetOutputColor(color string) {
	bis.OutputColor = color
}

type CmdInputScanner struct {
	Command string `json:"command"`

	BaseInputScanner
}

func (cis *CmdInputScanner) SetCmd(command string) {
	cis.Command = command
}

// Need to use a pointer receiver to modify the actual struct,
// Otherwise, a copy of the struct is passed.
func (cis *CmdInputScanner) Start() error {
	cmd := exec.Command("bash", "-c", cis.Command)

	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	outputColorAttr, ok := colorMap[cis.OutputColor]
	var outputColor color.Color
	if ok {
		outputColor = *color.New(outputColorAttr)
	} else {
		outputColor = *color.New(color.FgBlack)
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		outputColor.Println(m)
		// fmt.Fprintln(os.Stdout, colorRed, m, colorNone)
	}
	cmd.Wait()
	return nil
}

type K8sInputScanner struct {
	UseK8sTimestamp bool
	Pod             Pod

	BaseInputScanner
}

type Pod struct {
	Name        string
	Namespace   string
	PodSelector map[string]string
	Container   string
}

func (kis *K8sInputScanner) Start() error {
	_, err := exec.LookPath("kubectl")
	if err != nil {
		return err
	}

	if kis.Pod.Name != "" {
		log.Printf("111")
	} else if kis.Pod.PodSelector != nil {
		log.Printf("2222")
	} else {
		return errors.New("must provide one of: podName, podSelector")
	}

	return nil
}
