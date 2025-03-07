package main

import (
	"log"
	"sync"

	"github.com/kxiong0/bigbro/internal/config"
	"github.com/kxiong0/bigbro/internal/log_collector"
)

type BigBro struct {
	ConfigFilePath string
	LogOutputChan  chan log_collector.LogMsg

	wgInputScanners sync.WaitGroup
	inputScanners   []log_collector.InputScanner
}

func (bb *BigBro) AddInputScanner(scanner log_collector.InputScanner) {
	bb.inputScanners = append(bb.inputScanners, scanner)
}

func (bb *BigBro) GetInputScanners() []log_collector.InputScanner {
	return bb.inputScanners
}

func (bb *BigBro) Init() error {
	bb.LogOutputChan = make(chan log_collector.LogMsg, 10)

	// Load InputScanners according to config file at ConfigFilePath
	c := config.Config{}
	err := c.LoadConfigFile(bb.ConfigFilePath)
	if err != nil {
		return err
	}
	scanners := c.GetInputScanners()
	for i, scanner := range scanners {
		bb.AddInputScanner(scanner)
		scanner.Init(bb.LogOutputChan)
		scanner.SetID(i)
	}
	return nil
}

func (bb *BigBro) startScanner(scanner log_collector.InputScanner) {
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
