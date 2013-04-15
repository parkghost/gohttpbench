package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
)

const (
	GB_VERSION           = "0.1.6"
	MAX_RESPONSE_TIMEOUT = 30
	MAX_REQUESTS         = 50000 // if enable timelimit and without setting reqeusts
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
	PrintGCSummary()
}

func startBenchmark(context *Context) {
	PrintHeader()

	benchmark := NewBenchmark(context)
	monitor := NewMonitor(context, benchmark)
	go monitor.Run()
	go benchmark.Run()

	PrintReport(context, <-monitor.output)
}
