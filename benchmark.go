package main

import (
	"net/http"
	"time"
)

type Benchmark struct {
	c         *Context
	collector chan *Record
}

type Record struct {
	responseTime time.Duration
	contentSize  int64
	Error        error
}

func NewBenchmark(context *Context) *Benchmark {
	collector := make(chan *Record, context.config.requests)
	return &Benchmark{context, collector}
}

func (b *Benchmark) Run() {

	jobs := make(chan *http.Request, b.c.config.concurrency*GoMaxProcs)

	for i := 0; i < b.c.config.concurrency; i++ {
		go NewHTTPWorker(b.c, jobs, b.collector).Run()
	}

	base, _ := NewHTTPRequest(b.c.config)
	for i := 0; i < b.c.config.requests; i++ {
		jobs <- CopyHTTPRequest(b.c.config, base)
	}
	close(jobs)

	<-b.c.stop
}
