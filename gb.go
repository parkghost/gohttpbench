package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
)

const (
	GB_VERSION           = "0.1.3"
	MAX_RESPONSE_TIMEOUT = 30
	MAX_REQUESTS         = 50000 // if enable timelimit and without setting reqeusts
)

var (
	Verbosity       = 0
	GoMaxProcs      int
	ContinueOnError bool
)

func main() {
	if config, err := loadConfig(); err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(-1)
	} else {
		if err := detectHost(config); err != nil {
			log.Fatal(err)
		} else {
			runtime.GOMAXPROCS(GoMaxProcs)
			startBenchmark(config)
		}
	}
}

func startBenchmark(config *Config) {
	printHeader()

	start := &sync.WaitGroup{}
	start.Add(config.concurrency)
	stop := make(chan bool)

	benchmark := NewBenchmark(config, start, stop)
	monitor := NewMonitor(config, start, stop, benchmark)
	go monitor.Run()
	go benchmark.Run()

	printReport(config, <-monitor.output)
}
