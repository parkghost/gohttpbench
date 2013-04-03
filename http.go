package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type HttpWorker struct {
	config    *Config
	client    *http.Client
	jobs      chan *http.Request
	collector chan *Record
	start     chan bool
	stop      chan bool
}

func NewHttpWorker(config *Config, jobs chan *http.Request, collector chan *Record, start chan bool, stop chan bool) *HttpWorker {
	return &HttpWorker{config, NewClient(config), jobs, collector, start, stop}
}

var (
	invalidContnetSize = errors.New("invalid content size")
)

func (h *HttpWorker) Run() {
	<-h.start

	for job := range h.jobs {

		executionResult := make(chan *Record)

		go h.send(job, executionResult)

		select {
		case record := <-executionResult:
			h.collector <- record

		case <-time.After(time.Duration(MAX_RESPONSE_TIMEOUT) * time.Second):
			h.collector <- &Record{Error: errors.New("execution timeout")}

		case <-h.stop:
			return
		}
	}

}

func (h *HttpWorker) send(request *http.Request, executionResult chan<- *Record) {

	record := &Record{}
	sw := &StopWatch{}
	sw.Start()

	var contentSize int

	defer func() {
		if r := recover(); r != nil {
			record.Error = &ExceptionError{errors.New(fmt.Sprint(r))}
		} else {
			record.contentSize = contentSize
			record.responseTime = sw.Elapsed
		}

		if record.Error != nil {
			TraceException(record.Error)
		}

		executionResult <- record
	}()

	resp, err := h.client.Do(request)
	if err != nil {
		record.Error = &ConnectError{err}
		return
	} else {
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode > 300 {
			record.Error = &ResponseError{err}
			return
		}

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			record.Error = &ReceiveError{err}
			return
		}
		contentSize = len(body)

		if contentSize != h.config.contentSize {
			record.Error = &LengthError{invalidContnetSize}
			return
		}

	}

	sw.Stop()

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

		//TODO: another variable context
		config.serverName = resp.Header.Get("Server")
		config.contentSize = len(body)
	}

	return nil
}

func NewClient(config *Config) *http.Client {

	//skip certification check for self-signed certificates
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	//TODO: timeout control
	//TODO: tcp options
	//TODO: monitor tcp metrics
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
		//THINK: cache small file
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

type LengthError struct {
	err error
}

func (e *LengthError) Error() string {
	return e.err.Error()
}

type ConnectError struct {
	err error
}

func (e *ConnectError) Error() string {
	return e.err.Error()
}

type ReceiveError struct {
	err error
}

func (e *ReceiveError) Error() string {
	return e.err.Error()
}

type ExceptionError struct {
	err error
}

func (e *ExceptionError) Error() string {
	return e.err.Error()
}

type ResponseError struct {
	err error
}

func (e *ResponseError) Error() string {
	return e.err.Error()
}
