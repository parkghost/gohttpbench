package main

import (
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
			t.Errorf("expected %d, get %d", time.Duration(value)*time.Second, sw.Elapsed)
		}
	}

}
