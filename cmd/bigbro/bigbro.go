package main

import (
	"log"
	"sync"

	"github.com/kxiong0/bigbro/internal/scanner"
)

type BigBro struct {
	wgInputScanners sync.WaitGroup
	inputScanners   []scanner.InputScanner
}

func (bb *BigBro) Init() error {
	return nil
}

func (bb *BigBro) AddInputScanner(scanner scanner.InputScanner) {
	bb.inputScanners = append(bb.inputScanners, scanner)
}

func (bb *BigBro) startScanner(scanner scanner.InputScanner) {
	defer bb.wgInputScanners.Done()
	err := scanner.Start()
	if err != nil {
		log.Fatalf("Failed to start scanner: %s\n", err.Error())
	}
}

// Start all input scanners and block until all scanners stop
func (bb *BigBro) Start() {
	for _, scanner := range bb.inputScanners {
		bb.wgInputScanners.Add(1)
		go bb.startScanner(scanner)
	}
	bb.wgInputScanners.Wait()
}
