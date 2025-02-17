package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/kxiong0/bigbro/internal/scanner"
)

func main() {
	//
	fmt.Print("1111")
	input_scanner := &scanner.CmdInputScanner{}
	input_scanner.SetCmd("echo hi; echo hi; echo hi; sleep 1; echo hi; sleep 1; echo hi; sleep 1; echo hi;")
	input_scanner.SetOutputColor(color.New(color.FgBlue))

	input_scanner_2 := &scanner.CmdInputScanner{}
	input_scanner_2.SetCmd("echo yyo; echo yo; echo yo; sleep 1; echo yo; sleep 1; echo yo; sleep 1; echo yo;")
	input_scanner_2.SetOutputColor(color.New(color.FgRed))
	bb := BigBro{}
	bb.AddInputScanner(input_scanner)
	bb.AddInputScanner(input_scanner_2)
	bb.Start()

	time.Sleep(10 * time.Second)
}
