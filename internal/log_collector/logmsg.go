package log_collector

import "time"

type LogMsg struct {
	Timestamp  time.Time
	Line       string
	ScannerIdx int
}
