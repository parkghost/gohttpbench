package main

import (
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
