package main

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"time"
)

type StopWatch struct {
	start   time.Time
	Elapsed time.Duration
}

func (s *StopWatch) Start() {
	s.start = time.Now()
}

func (s *StopWatch) Stop() {
	s.Elapsed = time.Now().Sub(s.start)
}

func TraceException(msg interface{}) {
	switch {
	case Verbosity > 1:
		// print recovered error and stacktrace
		var buffer bytes.Buffer
		buffer.WriteString(fmt.Sprintf("recover: %v\n", msg))
		for skip := 1; ; skip++ {
			pc, file, line, ok := runtime.Caller(skip)
			if !ok {
				break
			}
			f := runtime.FuncForPC(pc)
			buffer.WriteString(fmt.Sprintf("\t%s:%d %s()\n", file, line, f.Name()))
		}
		buffer.WriteString("\n")
		fmt.Fprint(os.Stderr, buffer.String())
	case Verbosity > 0:
		// print recovered error only
		fmt.Fprintf(os.Stderr, "recover: %v\n", msg)
	}
}
