package main

import (
	"log"
	"sync"

	"github.com/kxiong0/bigbro/internal/config"
	"github.com/kxiong0/bigbro/internal/scanner"
)

type BigBro struct {
	wgInputScanners sync.WaitGroup
	inputScanners   []scanner.InputScanner
	ConfigFilePath  string
}

func (bb *BigBro) AddInputScanner(scanner scanner.InputScanner) {
	bb.inputScanners = append(bb.inputScanners, scanner)
}

func (bb *BigBro) GetInputScanners() []scanner.InputScanner {
	return bb.inputScanners
}

func (bb *BigBro) GetScannerChans() []chan string {
	chans := []chan string{}
	for _, value := range bb.inputScanners {
		chans = append(chans, value.GetOutputChan())
	}
	return chans
}

func (bb *BigBro) Init() error {
	// Load InputScanners according to config file at ConfigFilePath
	c := config.Config{}
	err := c.LoadConfigFile(bb.ConfigFilePath)
	if err != nil {
		return err
	}
	scanners := c.GetInputScanners()
	for i, scanner := range scanners {
		bb.AddInputScanner(scanner)
		scanner.Init()
		scanner.SetID(i)
	}
	return nil
}

func (bb *BigBro) startScanner(scanner scanner.InputScanner) {
	defer bb.wgInputScanners.Done()
	err := scanner.Start()
	if err != nil {
		log.Fatalf("Failed to run scanner: %s\n", err.Error())
	}
}

// Start all input scanners and block until all scanners stop
func (bb *BigBro) Start() {
	// for running scanners
	for _, scanner := range bb.inputScanners {
		bb.wgInputScanners.Add(1)
		go bb.startScanner(scanner)
	}
	log.Println("Scanners started")
}

func (bb *BigBro) Stop() {
	// wait for scanners to complete
	bb.wgInputScanners.Wait()
	for _, scanner := range bb.inputScanners {
		scanner.Close()
	}
	log.Println("Scanners stopped")
}
