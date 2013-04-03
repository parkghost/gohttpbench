package main

import (
	"testing"
)

func TestStdDev(t *testing.T) {
	testData := map[float64][]int64{
		2.0:    []int64{2, 4, 4, 4, 5, 5, 7, 9},
		1.5811: []int64{5, 6, 8, 9},
	}

	for expectedData, testingData := range testData {

		if result := stdDev(testingData); int(result*1000) != int(expectedData*1000) {
			t.Errorf("expected %d, get", expectedData, result)
		}
	}
}
