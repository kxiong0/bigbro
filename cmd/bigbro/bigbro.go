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
		log.Fatalf("Failed to run scanner: %s\n", err.Error())
	}
	log.Println("Start scanner return")
}

// Start all input scanners and block until all scanners stop
func (bb *BigBro) Start() {
	var listenerwg sync.WaitGroup
	for i, scanner := range bb.inputScanners {
		listenerwg.Add(1)
		go func() {
			for {
				log.Printf("reading...")
				val, ok := <-scanner.GetOutputChan()
				if !ok {
					break
				}
				log.Printf("GORoutine %d - %s", i, val)
			}
			listenerwg.Done()
		}()
	}
	for _, scanner := range bb.inputScanners {
		bb.wgInputScanners.Add(1)
		scanner.Init()
		go bb.startScanner(scanner)
	}
	bb.wgInputScanners.Wait()
	for _, scanner := range bb.inputScanners {
		scanner.Close()
	}
	listenerwg.Wait()
}
