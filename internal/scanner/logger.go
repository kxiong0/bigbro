package scanner

import (
	"fmt"
	"io"
	"time"
)

type TimestampWriter struct {
	Writer io.Writer
}

func (tw *TimestampWriter) Write(p []byte) (n int, err error) {
	timestamp := time.Now().Format(time.RFC3339Nano)
	formatted := fmt.Sprintf("%s %s", timestamp, p)
	return tw.Writer.Write([]byte(formatted))
}
