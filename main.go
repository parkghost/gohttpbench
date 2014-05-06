package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

const (
	GBVersion           = "0.1.9"
	MaxExecutionTimeout = time.Duration(30) * time.Second
	MaxRequests         = 50000 // for timelimit
)

var (
	Verbosity       = 0
	GoMaxProcs      int
	ContinueOnError bool
)

func main() {
	if config, err := LoadConfig(); err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(-1)
	} else {
		context := NewContext(config)
		if err := DetectHost(context); err != nil {
			log.Fatal(err)
		} else {
			runtime.GOMAXPROCS(GoMaxProcs)
			startBenchmark(context)
		}
	}
}

func startBenchmark(context *Context) {
	PrintHeader()

	benchmark := NewBenchmark(context)
	monitor := NewMonitor(context, benchmark.collector)
	go monitor.Run()
	go benchmark.Run()

	PrintReport(context, <-monitor.output)
}
