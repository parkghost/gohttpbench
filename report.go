package main

import (
	"bytes"
	"fmt"
	"math"
	"net/url"
	"sort"
)

func printHeader() {
	fmt.Println(`
This is GoHttpBench, Version ` + GB_VERSION + `, https://github.com/parkghost/gohttpbench
Author: Brandon Chen, Email: parkghost@gmail.com
Licensed under the Apache License, Version 2.0
`)
}

func printReport(config *Config, stats *Stats) {

	var buffer bytes.Buffer

	responseTimeData := stats.responseTimeData
	responseTimeDataIdx := stats.responseTimeDataIdx
	totalFailedReqeusts := stats.totalFailedReqeusts
	totalRequests := stats.totalRequests
	totalExecutionTime := stats.totalExecutionTime
	totalReceived := stats.totalReceived

	URL, _ := url.Parse(config.url)

	fmt.Fprint(&buffer, "\n\n")
	fmt.Fprintf(&buffer, "Server Software:        %s\n", config.serverName)
	fmt.Fprintf(&buffer, "Server Hostname:        %s\n", config.host)
	fmt.Fprintf(&buffer, "Server Port:            %d\n\n", config.port)

	fmt.Fprintf(&buffer, "Document Path:          %s\n", URL.RequestURI())
	fmt.Fprintf(&buffer, "Document Length:        %d bytes\n\n", config.contentSize)

	fmt.Fprintf(&buffer, "Concurrency Level:      %d\n", config.concurrency)
	fmt.Fprintf(&buffer, "Time taken for tests:   %.2f seconds\n", totalExecutionTime.Seconds())
	fmt.Fprintf(&buffer, "Complete requests:      %d\n", totalRequests)
	if totalFailedReqeusts == 0 {
		fmt.Fprintln(&buffer, "Failed requests:        0")
	} else {
		fmt.Fprintf(&buffer, "Failed requests:        %d\n", totalFailedReqeusts)
		fmt.Fprintf(&buffer, "   (Connect: %d, Receive: %d, Length: %d, Exceptions: %d)\n", stats.errConnect, stats.errReceive, stats.errLength, stats.errException)
	}
	if stats.errResponse > 0 {
		fmt.Fprintf(&buffer, "Non-2xx responses:      %d\n", stats.errResponse)
	}
	fmt.Fprintf(&buffer, "HTML transferred:       %d bytes\n", totalReceived)

	if responseTimeDataIdx > 0 && totalExecutionTime > 0 {
		stdDevOfResponseTime := StdDev(responseTimeData[:responseTimeDataIdx]) / 1000000
		sort.Sort(Int64Slice(responseTimeData))

		meanOfResponseTime := int64(totalExecutionTime) / int64(totalRequests) / 1000000
		medianOfResponseTime := responseTimeData[len(responseTimeData)/2] / 1000000
		minResponseTime := responseTimeData[0] / 1000000
		maxResponseTime := responseTimeData[len(responseTimeData)-1] / 1000000

		fmt.Fprintf(&buffer, "Requests per second:    %.2f [#/sec] (mean)\n", float64(totalRequests)/totalExecutionTime.Seconds())
		fmt.Fprintf(&buffer, "Time per request:       %.3f [ms] (mean)\n", float64(config.concurrency)*float64(totalExecutionTime.Nanoseconds())/1000000/float64(totalRequests))
		fmt.Fprintf(&buffer, "Time per request:       %.3f [ms] (mean, across all concurrent requests)\n", float64(totalExecutionTime.Nanoseconds())/1000000/float64(totalRequests))
		fmt.Fprintf(&buffer, "HTML Transfer rate:     %.2f [Kbytes/sec] received\n\n", float64(totalReceived/1024)/totalExecutionTime.Seconds())

		fmt.Fprint(&buffer, "Connection Times (ms)\n")
		fmt.Fprint(&buffer, "              min\tmean[+/-sd]\tmedian\tmax\n")
		fmt.Fprintf(&buffer, "Total:        %d     \t%d   %.2f \t%d \t%d\n\n",
			minResponseTime,
			meanOfResponseTime,
			stdDevOfResponseTime,
			medianOfResponseTime,
			maxResponseTime)

		fmt.Fprintln(&buffer, "Percentage of the requests served within a certain time (ms)")

		percentages := []int{50, 66, 75, 80, 90, 95, 98, 99}

		for _, percentage := range percentages {
			fmt.Fprintf(&buffer, " %d%%\t %d\n", percentage, responseTimeData[percentage*len(responseTimeData)/100]/1000000)
		}
		fmt.Fprintf(&buffer, " %d%%\t %d (longest request)\n", 100, maxResponseTime)
	}
	fmt.Println(buffer.String())
}

// custom sortable []int64 data type
type Int64Slice []int64

func (s Int64Slice) Len() int           { return len(s) }
func (s Int64Slice) Less(i, j int) bool { return s[i] < s[j] }
func (s Int64Slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// calculate standard deviation
func StdDev(data []int64) float64 {
	var sum int64 = 0
	for _, i := range data {
		sum += i
	}
	avg := float64(sum / int64(len(data)))

	sumOfSquares := 0.0
	for _, i := range data {

		sumOfSquares += math.Pow(float64(i)-avg, 2)
	}
	return math.Sqrt(sumOfSquares / float64(len(data)))

}
