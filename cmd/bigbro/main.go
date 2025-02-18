package main

import (
	"log"

	"github.com/kxiong0/bigbro/internal/config"
)

func main() {
	c := config.Config{}
	err := c.LoadConfigFile("config/default.json")
	if err != nil {
		log.Fatal(err)
	}

	scanners := c.GetInputScanners()
	bb := BigBro{}
	for _, scanner := range scanners {
		bb.AddInputScanner(scanner)
	}
	bb.Start()
}
