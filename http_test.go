package main

import (
	"testing"
)

var DefaultConfig = &Config{
	method: "GET",
	url:    "http://localhost/",
}

func BenchmarkNewHttpRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewHttpRequest(DefaultConfig)
		if err != nil {
			b.Fail()
		}
	}
	b.ReportAllocs()
}

func BenchmarkCopyHttpRequest(b *testing.B) {
	base, err := NewHttpRequest(DefaultConfig)
	if err != nil {
		b.Fail()
	}
	for i := 0; i < b.N; i++ {
		CopyHttpRequest(DefaultConfig, base)
	}
	b.ReportAllocs()
}
