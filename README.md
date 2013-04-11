Go-HttpBench
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
  -G=2: Number of Goroutine procs
  -H=[]: Add Arbitrary header line, eg. 'Accept-Encoding: gzip' Inserted after all normal header lines. (repeatable)
  -T="text/plain": Content-type header for POSTing, eg. 'application/x-www-form-urlencoded' Default is 'text/plain'
  -c=1: Number of multiple requests to make
  -h=false: Display usage information (this message)
  -k=false: Use HTTP KeepAlive feature
  -n=1: Number of requests to perform
  -p="": File containing data to POST. Remember also to set -T
  -r=false: Don't exit on socket receive errors
  -t=0: Seconds to max. wait for responses
  -u="": File containing data to PUT. Remember also to set -T
  -v=0: How much troubleshooting info to print
  -z=false: Use HTTP Gzip feature
```

### Example:
	$ gb -c 1000 -n 100000 -k http://localhost/10k.dat

	This is GoHttpBench, Version 0.1.4, https://github.com/parkghost/gohttpbench
	Author: Brandon Chen, Email: parkghost@gmail.com
	Licensed under the Apache License, Version 2.0

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


	Server Software:        nginx/1.2.1
	Server Hostname:        localhost
	Server Port:            80

	Document Path:          /10k.dat
	Document Length:        10240 bytes

	Concurrency Level:      1000
	Time taken for tests:   5.58 seconds
	Complete requests:      100000
	Failed requests:        0
	HTML transferred:       1024000000 bytes
	Requests per second:    17912.36 [#/sec] (mean)
	Time per request:       55.827 [ms] (mean)
	Time per request:       0.056 [ms] (mean, across all concurrent requests)
	HTML Transfer rate:     179123.56 [Kbytes/sec] received

	Connection Times (ms)
	              min	mean[+/-sd]	median	max
	Total:        0     	0   19.13 	49 	207

	Percentage of the requests served within a certain time (ms)
	 50%	 49
	 66%	 55
	 75%	 59
	 80%	 62
	 90%	 70
	 95%	 81
	 98%	 104
	 99%	 138
	 100%	 207 (longest request)
	 
Authors
-------

**Brandon Chen**

+ http://brandonc.me
+ http://github.com/parkghost


License
---------------------

Licensed under the Apache License, Version 2.0: http://www.apache.org/licenses/LICENSE-2.0