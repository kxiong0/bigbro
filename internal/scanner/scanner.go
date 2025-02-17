package scanner

import (
	"bufio"
	"os/exec"

	"github.com/fatih/color"
)

type InputScanner interface {
	SetOutputColor(*color.Color)
	Start() error
}

type BaseInputScanner struct {
	name        string
	outputColor *color.Color
}

func (bis *BaseInputScanner) SetName(name string) {
	bis.name = name
}

func (bis *BaseInputScanner) SetOutputColor(color *color.Color) {
	bis.outputColor = color
}

type CmdInputScanner struct {
	cmd string

	BaseInputScanner
}

func (cis *CmdInputScanner) SetCmd(cmd string) {
	cis.cmd = cmd
}

// Need to use a pointer receiver to modify the actual struct,
// Otherwise, a copy of the struct is passed.
func (cis *CmdInputScanner) Start() error {
	cmd := exec.Command("bash", "-c", cis.cmd)

	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	var outputColor color.Color
	if cis.outputColor != nil {
		outputColor = *cis.outputColor
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
