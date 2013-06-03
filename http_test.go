package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var getRequestConfig = &Config{
	url:              "http://localhost/",
	method:           "GET",
	executionTimeout: time.Duration(100) * time.Millisecond,
}
var postRequestConfig = &Config{
	url:              "http://localhost/",
	method:           "POST",
	contentType:      "application/x-www-form-urlencoded",
	executionTimeout: time.Duration(100) * time.Millisecond,
}

func init() {
	loadFile(postRequestConfig, "testdata/postfile.txt")
}

func TestHTTPWithGet(t *testing.T) {

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
		executionTimeout: MaxExecutionTimeout,
		url:              ts.URL,
	}

	context := NewContext(config)
	context.SetInt(FieldContentSize, len(responseStr))
	jobs := make(chan *http.Request)
	collector := make(chan *Record)

	worker := NewHTTPWorker(context, jobs, collector)

	go worker.Run()

	request, err := NewHTTPRequest(config)
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

func TestHTTPWithPost(t *testing.T) {

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
		executionTimeout: MaxExecutionTimeout,
		url:              ts.URL,
	}
	loadFile(config, "testdata/postfile.txt")

	context := NewContext(config)
	context.SetInt(FieldContentSize, len(responseStr))
	jobs := make(chan *http.Request)
	collector := make(chan *Record)

	worker := NewHTTPWorker(context, jobs, collector)

	go worker.Run()

	request, err := NewHTTPRequest(config)

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

func TestHTTPWorkerWithTimeout(t *testing.T) {

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
	context.SetInt(FieldContentSize, len(responseStr))
	jobs := make(chan *http.Request)
	collector := make(chan *Record)

	worker := NewHTTPWorker(context, jobs, collector)

	go worker.Run()

	request, err := NewHTTPRequest(config)
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

func BenchmarkNewHTTPRequestWithGet(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		NewHTTPRequest(getRequestConfig)
	}
}

func BenchmarkNewHTTPRequestWithPost(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		NewHTTPRequest(postRequestConfig)
	}
}

func BenchmarkCopyHTTPRequestWithGet(b *testing.B) {
	b.ReportAllocs()
	base, _ := NewHTTPRequest(getRequestConfig)
	for i := 0; i < b.N; i++ {
		CopyHTTPRequest(getRequestConfig, base)
	}
}

func BenchmarkCopyHTTPRequestWithPost(b *testing.B) {
	b.ReportAllocs()
	base, _ := NewHTTPRequest(postRequestConfig)
	for i := 0; i < b.N; i++ {
		CopyHTTPRequest(postRequestConfig, base)
	}
}
