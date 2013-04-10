package main

import (
	"net/http"
	"sync"
	"time"
)

type Benchmark struct {
	config *Config
	start  *sync.WaitGroup
	stop   chan bool

	collector chan *Record
}

type Record struct {
	responseTime time.Duration
	contentSize  int64
	Error        error
}

func NewBenchmark(config *Config, start *sync.WaitGroup, stop chan bool) *Benchmark {
	collector := make(chan *Record, config.requests)
	return &Benchmark{config, start, stop, collector}
}

func (b *Benchmark) Run() {

	jobs := make(chan *http.Request, b.config.concurrency*GoMaxProcs)

	for i := 0; i < b.config.concurrency; i++ {
		go NewHttpWorker(b.config, b.start, b.stop, jobs, b.collector).Run()
	}

	for i := 0; i < b.config.requests; i++ {
		request, _ := NewHttpRequest(b.config)
		jobs <- request
	}

	<-b.stop
}
