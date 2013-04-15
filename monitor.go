package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

type Monitor struct {
	c         *Context
	collector chan *Record
	output    chan *Stats
}

type Stats struct {
	responseTimeData []time.Duration

	totalRequests       int
	totalExecutionTime  time.Duration
	totalResponseTime   time.Duration
	totalReceived       int64
	totalFailedReqeusts int

	errLength    int
	errConnect   int
	errReceive   int
	errException int
	errResponse  int
}

func NewMonitor(context *Context, benchmark *Benchmark) *Monitor {
	return &Monitor{context, benchmark.collector, make(chan *Stats)}
}

func (m *Monitor) Run() {

	// catch interrupt signal
	userInterrupt := make(chan os.Signal, 1)
	signal.Notify(userInterrupt, os.Interrupt)

	stats := &Stats{}
	stats.responseTimeData = make([]time.Duration, 0, m.c.config.requests)

	var timelimit time.Timer
	if m.c.config.timelimit > 0 {
		timelimit = *time.NewTimer(time.Duration(m.c.config.timelimit) * time.Second)
	}
	defer timelimit.Stop()

	// waiting for all of http workers to start
	m.c.start.Wait()

	fmt.Printf("Benchmarking %s (be patient)\n", m.c.config.host)
	sw := &StopWatch{}
	sw.Start()

loop:
	for {
		select {
		case record := <-m.collector:

			updateStats(stats, record)

			if record.Error != nil && !ContinueOnError {
				break loop
			}

			if stats.totalRequests >= 10 && stats.totalRequests%(m.c.config.requests/10) == 0 {
				fmt.Printf("Completed %d requests\n", stats.totalRequests)
			}

			if stats.totalRequests == m.c.config.requests {
				fmt.Printf("Finished %d requests\n", stats.totalRequests)
				break loop
			}

		case <-timelimit.C:
			break loop
		case <-userInterrupt:
			break loop
		}
	}

	sw.Stop()
	stats.totalExecutionTime = sw.Elapsed

	// shutdown benchmark and all of httpworkers to stop
	close(m.c.stop)
	m.output <- stats
}

func updateStats(stats *Stats, record *Record) {
	stats.totalRequests += 1

	if record.Error != nil {
		stats.totalFailedReqeusts += 1

		switch record.Error.(type) {
		case *ConnectError:
			stats.errConnect += 1
		case *ExceptionError:
			stats.errException += 1
		case *LengthError:
			stats.errLength += 1
		case *ReceiveError:
			stats.errReceive += 1
		case *ResponseError:
			stats.errResponse += 1
		default:
			stats.errException += 1
		}

	} else {
		stats.totalResponseTime += record.responseTime
		stats.totalReceived += record.contentSize
		stats.responseTimeData = append(stats.responseTimeData, record.responseTime)
	}

}
