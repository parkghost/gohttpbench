package main

import (
	"testing"
	"time"
)

func TestStdDev(t *testing.T) {
	testData := map[float64][]time.Duration{
		2.0:    []time.Duration{2, 4, 4, 4, 5, 5, 7, 9},
		1.5811: []time.Duration{5, 6, 8, 9},
	}

	for expectedData, testingData := range testData {
		if result := stdDev(testingData); int(result*1000) != int(expectedData*1000) {
			t.Errorf("expected %f, got %f", expectedData, result)
		}
	}
}
