package main

import (
	"net/url"
	"testing"
	"time"
)

func TestStopWatch(t *testing.T) {
	testData := []int{1, 2}
	for _, value := range testData {
		sw := &StopWatch{}

		sw.Start()
		time.Sleep(time.Duration(value) * time.Second)
		sw.Stop()

		if int64(sw.Elapsed)/1000000 == int64(time.Duration(value)*time.Second) {
			t.Errorf("expected %d, get", time.Duration(value)*time.Second, sw.Elapsed)
		}
	}

}

func TestExtractHostAndPort(t *testing.T) {

	type Pair struct {
		host string
		port int
	}

	testData := map[string]Pair{
		"http://localhost:8080/":  Pair{"localhost", 8080},
		"https://www.google.com/": Pair{"www.google.com", 443},
		"http://localhost/":       Pair{"localhost", 80},
	}

	for testingData, expectedData := range testData {
		URL, _ := url.Parse(testingData)
		host, port := extractHostAndPort(URL)

		if host != expectedData.host && port != expectedData.port {
			t.Errorf("expected host:%s and port:%d, get", host, port)
		}
	}
}
