package scanner

import (
	"bufio"
	"errors"
	"log"
	"os/exec"
	"time"

	"github.com/fatih/color"
)

type InputScanner interface {
	SetName(string)
	SetOutputColor(*color.Color)
	Start() error
}

type BaseInputScanner struct {
	name            string
	refreshInterval time.Duration
	outputColor     *color.Color
}

func (bis *BaseInputScanner) SetName(name string) {
	bis.name = name
}

func (bis *BaseInputScanner) SetOutputColor(color *color.Color) {
	bis.outputColor = color
}

func (bis *BaseInputScanner) SetRefreshInterval(interval string) {
	// Interval should be a duration string
	duration, err := time.ParseDuration(interval)
	if err != nil {
		log.Fatal("Error parsing duration string:", err)
		return
	}
	bis.refreshInterval = duration
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

type K8sInputScanner struct {
	podName     string
	namespace   string
	podSelector map[string]string
	container   string

	BaseInputScanner
}

func (kis *K8sInputScanner) SetPodName(name string) {
	kis.podName = name
}

func (kis *K8sInputScanner) SetNamespace(namespace string) {
	kis.namespace = namespace
}

func (kis *K8sInputScanner) SetPodSelector(podSelector map[string]string) {
	kis.podSelector = podSelector
}

func (kis *K8sInputScanner) SetContainer(container string) {
	kis.container = container
}

func (kis *K8sInputScanner) Start() error {
	_, err := exec.LookPath("kubectl")
	if err != nil {
		return err
	}

	if kis.podName != "" {
		log.Printf("111")
	} else if kis.podSelector != nil {
		log.Printf("2222")
	} else {
		return errors.New("must provide one of: podName, podSelector")
	}

	return nil
}
