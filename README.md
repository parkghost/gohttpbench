Go-HttpBench

[![Build Status](https://drone.io/github.com/parkghost/gohttpbench/status.png)](https://drone.io/github.com/parkghost/gohttpbench/latest)

====

*an ab-like benchmark tool run on multi-core cpu*

Installation
--------------
1. install [Go](http://golang.org/doc/install) into your environment
2. download and build Go-HttpBench

```
go get github.com/parkghost/gohttpbench
go build -o gb github.com/parkghost/gohttpbench
```

Usage
-----------

```
Usage: gb [options] http[s]://hostname[:port]/path
Options are:
  -A="": Add Basic WWW Authentication, the attributes are a colon separated username and password.
  -C=[]: Add cookie, eg. 'Apache=1234. (repeatable)
  -G=2: Number of CPU
  -H=[]: Add Arbitrary header line, eg. 'Accept-Encoding: gzip' Inserted after all normal header lines. (repeatable)
  -T="text/plain": Content-type header for POSTing, eg. 'application/x-www-form-urlencoded' Default is 'text/plain'
  -c=1: Number of multiple requests to make
  -h=false: Display usage information (this message)
  -i=false: Use HEAD instead of GET
  -k=false: Use HTTP KeepAlive feature
  -n=1: Number of requests to perform
  -p="": File containing data to POST. Remember also to set -T
  -r=false: Don't exit when errors
  -t=0: Seconds to max. wait for responses
  -u="": File containing data to PUT. Remember also to set -T
  -v=0: How much troubleshooting info to print
  -z=false: Use HTTP Gzip feature
```

### Example:
	$ gb -c 100 -n 100000 -k http://localhost/10k.dat

	This is GoHttpBench, Version 0.1.9, https://github.com/parkghost/gohttpbench
	Author: Brandon Chen, Email: parkghost@gmail.com
	Licensed under the MIT license

	Benchmarking localhost (be patient)
	Completed 10000 requests
	Completed 20000 requests
	Completed 30000 requests
	Completed 40000 requests
	Completed 50000 requests
	Completed 60000 requests
	Completed 70000 requests
	Completed 80000 requests
	Completed 90000 requests
	Completed 100000 requests
	Finished 100000 requests


	Server Software:        nginx/1.2.6 (Ubuntu)
	Server Hostname:        localhost
	Server Port:            80

	Document Path:          /10k.dat
	Document Length:        10240 bytes

	Concurrency Level:      100
	Time taken for tests:   5.36 seconds
	Complete requests:      100000
	Failed requests:        0
	HTML transferred:       1024000000 bytes
	Requests per second:    18652.41 [#/sec] (mean)
	Time per request:       5.361 [ms] (mean)
	Time per request:       0.054 [ms] (mean, across all concurrent requests)
	HTML Transfer rate:     186524.05 [Kbytes/sec] received

	Connection Times (ms)
	              min	mean[+/-sd]	median	max
	Total:        0     	0   3.03 	4 	32

	Percentage of the requests served within a certain time (ms)
	 50%	 4
	 66%	 5
	 75%	 6
	 80%	 7
	 90%	 8
	 95%	 10
	 98%	 12
	 99%	 14
	 100%	 32 (longest request)


Author
-------

**Brandon Chen**

+ http://brandonc.me
+ http://github.com/parkghost


License
---------------------

This project is licensed under the MIT license