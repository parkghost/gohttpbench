package main

import (
	"net/http"
	"time"
)

type Benchmark struct {
	config    *Config
	collector chan *Record
	monitor   *Monitor
	start     chan bool
	stop      chan bool
}

type Record struct {
	responseTime time.Duration
	contentSize  int
	Error        error
}

func NewBenchmark(config *Config) *Benchmark {
	start := make(chan bool)
	stop := make(chan bool)

	collector := make(chan *Record, config.requests)
	monitor := NewMonitor(config, collector, start, stop)

	return &Benchmark{config, collector, monitor, start, stop}
}

func (b *Benchmark) Run() {

	go b.monitor.Run()

	jobs := make(chan *http.Request, b.config.requests)

	for i := 0; i < b.config.concurrency; i++ {
		go NewHttpWorker(b.config, jobs, b.collector, b.start, b.stop).Run()
	}

	for i := 0; i < b.config.requests; i++ {
		request, _ := NewHttpRequest(b.config)
		jobs <- request
	}

	close(b.start)
	<-b.stop
}
