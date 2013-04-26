package main

import (
	"errors"
	"os"
	"testing"
	"time"
)

func TestMonitorWithSuccessedResponse(t *testing.T) {

	config := &Config{
		requests: 2,
	}

	collector := make(chan *Record, config.requests)

	context := NewContext(config)
	monitor := NewMonitor(context, collector)

	request1 := &Record{10, 10, nil}
	request2 := &Record{20, 20, nil}

	collector <- request1
	collector <- request2

	devnull, _ := os.Open(os.DevNull)
	defer devnull.Close()

	stdout := os.Stdout
	os.Stdout = devnull

	go monitor.Run()
	stats := <-monitor.output

	os.Stdout = stdout
	if stats.totalRequests != config.requests {
		t.Fatalf("expected %d requests, actual %d requests", config.requests, stats.totalRequests)
	}

	if stats.responseTimeData[0] != request1.responseTime || stats.responseTimeData[1] != request2.responseTime {
		t.Fatalf("expected %s responseTimeData, actual %s responseTimeData", []time.Duration{request1.responseTime, request2.responseTime}, stats.responseTimeData)
	}

	if stats.totalReceived != request1.contentSize+request2.contentSize {
		t.Fatalf("expected %d content received, actual %d content received", request1.contentSize+request2.contentSize, stats.totalReceived)
	}
}

func TestMonitorWithFailedResponse(t *testing.T) {

	ContinueOnError = true

	config := &Config{
		requests: 6,
	}

	collector := make(chan *Record, config.requests)

	context := NewContext(config)
	monitor := NewMonitor(context, collector)

	dummy := errors.New("dummy error")

	records := []*Record{
		&Record{Error: &LengthError{dummy}},
		&Record{Error: &ConnectError{dummy}},
		&Record{Error: &ReceiveError{dummy}},
		&Record{Error: &ExceptionError{dummy}},
		&Record{Error: &ResponseError{dummy}},
		&Record{Error: &ResponseTimeoutError{dummy}},
	}

	expectedStat := &Stats{
		totalRequests:       config.requests,
		totalFailedReqeusts: 6,
		errLength:           1,
		errConnect:          1,
		errReceive:          1,
		errException:        2,
		errResponse:         1,
	}

	for _, record := range records {
		collector <- record
	}

	devnull, _ := os.Open(os.DevNull)
	defer devnull.Close()

	stdout := os.Stdout
	os.Stdout = devnull

	go monitor.Run()
	actualStats := <-monitor.output
	os.Stdout = stdout

	if actualStats.totalRequests != expectedStat.totalRequests ||
		actualStats.totalFailedReqeusts != expectedStat.totalFailedReqeusts ||
		actualStats.errLength != expectedStat.errLength ||
		actualStats.errConnect != expectedStat.errConnect ||
		actualStats.errReceive != expectedStat.errReceive ||
		actualStats.errException != expectedStat.errException ||
		actualStats.errResponse != expectedStat.errResponse {
		t.Fatalf("expected %#+v , actual %#+v", expectedStat, actualStats)
	}

}
