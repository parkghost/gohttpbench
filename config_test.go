package main

import (
	"net/url"
	"testing"
)

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
			t.Errorf("expected host:%s and port:%d, got", host, port)
		}
	}
}
