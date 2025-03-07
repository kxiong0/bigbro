package log_collector

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"time"
)

type InputScanner interface {
	Init(chan LogMsg)
	Start() error
	Close()
	SetName(string)
	SetID(int)
	GetID() int
	GetColor() string
	GetOutputChan() chan LogMsg
}

type BaseInputScanner struct {
	Name       string `json:"name"`
	Color      string `json:"color"`
	OutputChan chan LogMsg
	Id         int
}

func (bis *BaseInputScanner) Init(channel chan LogMsg) {
	bis.OutputChan = channel
}

func (bis *BaseInputScanner) Close() {
}

func (bis *BaseInputScanner) SetName(name string) {
	bis.Name = name
}

func (bis *BaseInputScanner) SetID(id int) {
	bis.Id = id
}

func (bis *BaseInputScanner) GetColor() string {
	return bis.Color
}

func (bis *BaseInputScanner) GetID() int {
	return bis.Id
}

func (bis *BaseInputScanner) GetOutputChan() chan LogMsg {
	return bis.OutputChan
}
func (bis *BaseInputScanner) Start() error {
	bis.OutputChan <- LogMsg{Line: "BaseScanner"}
	return nil
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
		logmsg := LogMsg{Timestamp: time.Now(), Line: scanner.Text(), ScannerIdx: cis.Id}
		cis.OutputChan <- logmsg
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
