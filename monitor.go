package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

type Monitor struct {
	config    *Config
	collector chan *Record
	output    chan *Stats
	start     chan bool
	stop      chan bool
}

type Stats struct {
	responseTimeData    []int64
	responseTimeDataIdx int

	totalRequests      int
	totalExecutionTime time.Duration
	totalResponseTime  time.Duration
	totalReceived      int

	totalFailedReqeusts int
	errLength           int
	errConnect          int
	errReceive          int
	errException        int
	errResponse         int
}

func NewMonitor(config *Config, collector chan *Record, start chan bool, stop chan bool) *Monitor {
	return &Monitor{config, collector, make(chan *Stats), start, stop}
}

func (m *Monitor) Run() {

	userInterrupt := make(chan os.Signal, 1)
	signal.Notify(userInterrupt, os.Interrupt)

	<-m.start

	stats := &Stats{}
	stats.responseTimeData = make([]int64, m.config.requests)
	stats.responseTimeDataIdx = 0

	sw := &StopWatch{}
	sw.Start()

	var timelimit <-chan time.Time
	if m.config.timelimit > 0 {
		timelimit = time.After(time.Duration(m.config.timelimit) * time.Second)
	}

	fmt.Printf("Benchmarking %s (be patient)\n", m.config.host)

loop:
	for {
		select {
		case record := <-m.collector:

			updateStats(stats, record)

			if record.Error != nil && !ContinueOnError {
				break loop
			}

			if stats.totalRequests >= 10 && stats.totalRequests%(m.config.requests/10) == 0 {
				fmt.Printf("Completed %d requests\n", stats.totalRequests)
			}

			if stats.totalRequests == m.config.requests {
				fmt.Printf("Finished %d requests\n", stats.totalRequests)
				break loop
			}

		case <-timelimit:
			break loop
		case <-userInterrupt:
			break loop
		}
	}

	sw.Stop()
	stats.totalExecutionTime = sw.Elapsed

	//to notify benchmark and all of httpworkers stop running
	close(m.stop)
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

		stats.responseTimeData[stats.responseTimeDataIdx] = record.responseTime.Nanoseconds()
		stats.responseTimeDataIdx += 1
	}

}
