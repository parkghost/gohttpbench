package main

import (
	"testing"
	"time"
)

func TestStopWatch(t *testing.T) {
	testData := []int{100, 200} //Millisecond
	for _, value := range testData {
		sw := &StopWatch{}

		sw.Start()
		time.Sleep(time.Duration(value) * time.Millisecond)
		sw.Stop()

		if int64(sw.Elapsed)/1000000 != int64(time.Duration(value)) {
			t.Errorf("expected %d, got %d", time.Duration(value)*time.Millisecond, sw.Elapsed)
		}
	}

}
