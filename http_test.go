package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

func TestHttpGet(t *testing.T) {

	//fake http server
	responseStr := "hello"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(responseStr))
	}))
	defer ts.Close()

	// http worker

	config := &Config{
		concurrency:      1,
		requests:         1,
		method:           "GET",
		executionTimeout: MAX_EXECUTION_TIMEOUT,
		url:              ts.URL,
	}

	context := NewContext(config)
	context.SetInt(CONTENT_SIZE, len(responseStr))
	jobs := make(chan *http.Request)
	collector := make(chan *Record)

	worker := NewHttpWorker(context, jobs, collector)

	go worker.Run()

	request, err := NewHttpRequest(config)
	if err != nil {
		t.Fatalf("new http request failed: %s", err)
	}

	jobs <- request
	record := <-collector
	close(jobs)
	close(context.stop)

	if record.Error != nil {
		t.Fatalf("sent a http reqeust but was error: %s", record.Error)
	}

	if record.contentSize != int64(len(responseStr)) {
		t.Fatalf("send a http reqeust but content size dismatch")
	}
}

func TestHttpPost(t *testing.T) {

	//fake http server
	responseStr := "hello"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// from values from testdata/postfile.txt
		if r.FormValue("email") == "test" && r.FormValue("password") == "testing" {
			w.Write([]byte(responseStr))
		}
	}))
	defer ts.Close()

	// http worker

	config := &Config{
		concurrency:      1,
		requests:         1,
		method:           "POST",
		contentType:      "application/x-www-form-urlencoded",
		executionTimeout: MAX_EXECUTION_TIMEOUT,
		url:              ts.URL,
	}
	loadFile(config, "testdata/postfile.txt")

	context := NewContext(config)
	context.SetInt(CONTENT_SIZE, len(responseStr))
	jobs := make(chan *http.Request)
	collector := make(chan *Record)

	worker := NewHttpWorker(context, jobs, collector)

	go worker.Run()

	request, err := NewHttpRequest(config)

	if err != nil {
		t.Fatalf("new http request failed: %s", err)
	}

	jobs <- request
	record := <-collector
	close(jobs)
	close(context.stop)

	if record.Error != nil {
		t.Fatalf("sent a http reqeust but was error: %s", record.Error)
	}

	if record.contentSize != int64(len(responseStr)) {
		t.Fatalf("send a http reqeust but content size dismatch")
	}
}

func TestHttpWorkerWithTimeout(t *testing.T) {

	//fake http server
	responseStr := "hello"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(200) * time.Millisecond)
		w.Write([]byte(responseStr))
	}))
	defer ts.Close()

	// http worker

	config := &Config{
		concurrency:      1,
		requests:         1,
		method:           "GET",
		executionTimeout: time.Duration(100) * time.Millisecond,
		url:              ts.URL,
	}

	context := NewContext(config)
	context.SetInt(CONTENT_SIZE, len(responseStr))
	jobs := make(chan *http.Request)
	collector := make(chan *Record)

	worker := NewHttpWorker(context, jobs, collector)

	go worker.Run()

	request, err := NewHttpRequest(config)
	if err != nil {
		t.Fatalf("new http request failed: %s", err)
	}

	jobs <- request
	record := <-collector
	close(jobs)
	close(context.stop)

	if record.Error == nil {

		fmt.Println(record)
		t.Fatal("expected timeout error")
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
