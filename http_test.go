package main

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

var getRequestConfig = &Config{
	method: "GET",
	url:    "http://localhost/",
}
var postRequestConfig = &Config{
	method: "POST",
	url:    "http://localhost/",
}

func init() {
	loadFile(postRequestConfig, "testdata/postfile.txt")
}

func TestHttpWorker(t *testing.T) {

	requests := 1
	var received int64 = 0

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
		atomic.AddInt64(&received, 1)
	}))
	defer ts.Close()

	config := &Config{
		concurrency: 1,
		requests:    requests,
		method:      "GET",
		url:         ts.URL,
	}

	context := NewContext(config)
	context.SetInt(CONTENT_SIZE, 5)

	jobs := make(chan *http.Request, 1)
	collector := make(chan *Record, 1)

	worker := NewHttpWorker(context, jobs, collector)

	go worker.Run()
	go func() {
		counter := 0
		for record := range collector {
			counter += 1
			if counter == requests || record.Error != nil {
				break
			}
		}
		close(context.stop)
	}()

	request, err := NewHttpRequest(config)
	if err != nil {
		t.Fatalf("new http request failed: %s", err)
	}

	jobs <- request
	<-context.stop

	if actualReceived := atomic.LoadInt64(&received); int64(requests) != actualReceived {
		t.Fatalf("expected to send %d requests and receive %d responses, but get %d responses", requests, requests, actualReceived)
	}
}

func BenchmarkNewHttpRequest_Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewHttpRequest(getRequestConfig)
	}
	b.ReportAllocs()
}

func BenchmarkNewHttpRequest_Post(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewHttpRequest(postRequestConfig)
	}
	b.ReportAllocs()
}

func BenchmarkCopyHttpRequest_Get(b *testing.B) {
	base, _ := NewHttpRequest(getRequestConfig)
	for i := 0; i < b.N; i++ {
		CopyHttpRequest(getRequestConfig, base)
	}
	b.ReportAllocs()
}

func BenchmarkCopyHttpRequest_Post(b *testing.B) {
	base, _ := NewHttpRequest(postRequestConfig)
	for i := 0; i < b.N; i++ {
		CopyHttpRequest(postRequestConfig, base)
	}
	b.ReportAllocs()
}
