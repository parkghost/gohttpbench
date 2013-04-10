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
	$ gb -c 1000 -n 100000 -k http://localhost/

	This is GoHttpBench, Version 0.1.2, https://github.com/parkghost/gohttpbench
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

	Document Path:          /
	Document Length:        151 bytes

	Concurrency Level:      1000
	Time taken for tests:   4.94 seconds
	Complete requests:      100000
	Failed requests:        0
	HTML transferred:       0 bytes
	Requests per second:    20230.47 [#/sec] (mean)
	Time per request:       49.430 [ms] (mean)
	Time per request:       0.049 [ms] (mean, across all concurrent requests)
	HTML Transfer rate:     0.00 [Kbytes/sec] received

	Connection Times (ms)
	              min	mean[+/-sd]	median	max
	Total:        3     	0   18.09 	43 	187

	Percentage of the requests served within a certain time (ms)
	 50%	 43
	 66%	 47
	 75%	 51
	 80%	 54
	 90%	 66
	 95%	 78
	 98%	 99
	 99%	 124
	 100%	 187 (longest request)

Authors
-------

**Brandon Chen**

+ http://brandonc.me
+ http://github.com/parkghost


License
---------------------

Licensed under the Apache License, Version 2.0: http://www.apache.org/licenses/LICENSE-2.0