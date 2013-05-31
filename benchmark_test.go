package main

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestBenchmark(t *testing.T) {

	requests := 100
	var received int64 = 0

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
		atomic.AddInt64(&received, 1)
	}))
	defer ts.Close()

	config := &Config{
		concurrency:      10,
		requests:         requests,
		method:           "GET",
		executionTimeout: MAX_EXECUTION_TIMEOUT,
		url:              ts.URL,
	}

	context := NewContext(config)
	context.SetInt(CONTENT_SIZE, 5)
	benchmark := NewBenchmark(context)

	go benchmark.Run()

	go func() {
		counter := 0
		for record := range benchmark.collector {
			counter += 1
			if counter == requests || record.Error != nil {
				break
			}
		}
		close(context.stop)
	}()

	context.start.Wait()
	<-context.stop

	if actualReceived := atomic.LoadInt64(&received); int64(requests) != actualReceived {
		t.Fatalf("expected to send %d requests and receive %d responses, but got %d responses", requests, requests, actualReceived)
	}
}
