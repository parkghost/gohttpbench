package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
)

const (
	GB_VERSION           = "0.1.0"
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
	benchmark := NewBenchmark(config)
	benchmark.Run()
	report(config, benchmark.monitor.output)
}

func printHeader() {
	fmt.Println(`
This is GoHttpBench, Version ` + GB_VERSION + `, https://github.com/parkghost/gohttpbench
Author: Brandon Chen, Email: parkghost@gmail.com
Licensed under the Apache License, Version 2.0
`)
}
