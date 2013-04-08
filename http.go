package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	invalidContnetSize = errors.New("invalid content size")
)

type LengthError error
type ConnectError error
type ReceiveError error
type ExceptionError error
type ResponseError error

type HttpWorker struct {
	config *Config
	start  *sync.WaitGroup
	stop   chan bool

	client    *http.Client
	jobs      chan *http.Request
	collector chan *Record
}

func NewHttpWorker(config *Config, start *sync.WaitGroup, stop chan bool, jobs chan *http.Request, collector chan *Record) *HttpWorker {
	return &HttpWorker{config, start, stop, NewClient(config), jobs, collector}
}

func (h *HttpWorker) Run() {
	h.start.Done()
	h.start.Wait()

	for job := range h.jobs {

		asyncResult := h.send(job)
		timeout := time.NewTimer(time.Duration(MAX_RESPONSE_TIMEOUT) * time.Second)

		select {
		case record := <-asyncResult:
			h.collector <- record

		case <-timeout.C:
			h.collector <- &Record{Error: errors.New("execution timeout")}

		case <-h.stop:
			timeout.Stop()
			return
		}
		timeout.Stop()
	}
}

func (h *HttpWorker) send(request *http.Request) (asyncResult chan *Record) {

	asyncResult = make(chan *Record)
	go func() {
		record := &Record{}
		sw := &StopWatch{}
		sw.Start()

		var contentSize int

		defer func() {
			if r := recover(); r != nil {
				record.Error = ExceptionError(errors.New(fmt.Sprint(r)))
			} else {
				record.contentSize = contentSize
				record.responseTime = sw.Elapsed
			}

			if record.Error != nil {
				TraceException(record.Error)
			}

			asyncResult <- record
		}()

		resp, err := h.client.Do(request)
		if err != nil {
			record.Error = ConnectError(err)
			return
		} else {
			defer resp.Body.Close()

			if resp.StatusCode < 200 || resp.StatusCode > 300 {
				record.Error = ResponseError(err)
				return
			}

			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				record.Error = ReceiveError(err)
				return
			}

			expectedContentSize := 0
			headerContentSize := resp.Header.Get("Content-Length")

			if headerContentSize != "" {
				expectedContentSize, _ = strconv.Atoi(headerContentSize)
			} else {
				expectedContentSize = h.config.contentSize
			}

			if expectedContentSize != len(body) {
				record.Error = LengthError(invalidContnetSize)
				return
			}

		}

		sw.Stop()
	}()
	return asyncResult
}

func detectHost(config *Config) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			TraceException(r)
		}
	}()

	client := NewClient(config)
	reqeust, err := NewHttpRequest(config)
	if err != nil {
		return err
	}

	resp, err := client.Do(reqeust)

	if err != nil {
		return err
	} else {

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		// TODO: place on another context
		config.serverName = resp.Header.Get("Server")
		config.contentSize = len(body)
	}

	return nil
}

func NewClient(config *Config) *http.Client {

	// skip certification check for self-signed certificates
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	// TODO: timeout control
	// TODO: tcp options
	// TODO: monitor tcp metrics
	transport := &http.Transport{
		DisableCompression: !config.gzip,
		DisableKeepAlives:  !config.keepAlive,
		TLSClientConfig:    tlsconfig,
	}

	client := &http.Client{Transport: transport}
	return client
}

func NewHttpRequest(config *Config) (*http.Request, error) {

	var body io.Reader
	var err error

	if (config.method == "POST" || config.method == "PUT") && config.bodyFile != "" {
		// THINK: cache small file
		bytes, err := ioutil.ReadFile(config.bodyFile)
		if err != nil {
			return nil, err
		}

		body = strings.NewReader(string(bytes))
	}

	request, err := http.NewRequest(config.method, config.url, body)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", config.contentType)

	if config.keepAlive {
		request.Header.Set("Connection", "keep-alive")
	}

	for _, header := range config.headers {
		pair := strings.Split(header, ":")
		request.Header.Add(pair[0], pair[1])
	}

	for _, cookie := range config.cookies {
		pair := strings.Split(cookie, "=")
		c := &http.Cookie{Name: pair[0], Value: pair[1]}
		request.AddCookie(c)
	}

	if config.basicAuthentication != "" {
		pair := strings.Split(config.basicAuthentication, ":")
		request.SetBasicAuth(pair[0], pair[1])
	}

	request.Header.Add("User-Agent", config.userAgent)
	return request, err
}
