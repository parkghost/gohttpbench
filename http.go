package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	SERVER_NAME  = "ServerName"
	CONTENT_SIZE = "ContentSize"
)

var (
	invalidContnetSize = errors.New("invalid content size")
)

type HttpWorker struct {
	c         *Context
	client    *http.Client
	jobs      chan *http.Request
	collector chan *Record
	readBuf   *bytes.Buffer
}

func NewHttpWorker(context *Context, jobs chan *http.Request, collector chan *Record) *HttpWorker {
	return &HttpWorker{
		context,
		NewClient(context.config),
		jobs,
		collector,
		bytes.NewBuffer(make([]byte, 0, context.GetInt(CONTENT_SIZE)+bytes.MinRead)),
	}
}

func (h *HttpWorker) Run() {
	h.c.start.Done()
	h.c.start.Wait()

	for job := range h.jobs {

		asyncResult := h.send(job)
		// TODO:(Go 1.1) use timer.Reset(d) instead of create new timer
		timeout := time.NewTimer(time.Duration(MAX_RESPONSE_TIMEOUT) * time.Second)

		select {
		case record := <-asyncResult:
			h.collector <- record

		case <-timeout.C:
			h.collector <- &Record{Error: &ResponseTimeoutError{errors.New("execution timeout")}}

			// TODO:(Go 1.1) timeout control https://code.google.com/p/go/issues/detail?id=3362
			// h.client.Transport.(*http.Transport).CancelRequest(job)
		case <-h.c.stop:
			timeout.Stop()
			// h.client.Transport.(*http.Transport).CancelRequest(job)
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

		var contentSize int64

		defer func() {
			if r := recover(); r != nil {
				if Err, ok := r.(error); ok {
					record.Error = Err
				} else {
					record.Error = &ExceptionError{errors.New(fmt.Sprint(r))}
				}

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
			record.Error = &ConnectError{err}
			return
		} else {
			defer resp.Body.Close()

			if resp.StatusCode < 200 || resp.StatusCode > 300 {
				record.Error = &ResponseError{err}
				return
			}

			defer h.readBuf.Reset()
			var err error
			contentSize, err = h.readBuf.ReadFrom(resp.Body)

			if err != nil {
				record.Error = &ReceiveError{err}
				return
			}

			expectedContentSize := 0
			headerContentSize := resp.Header.Get("Content-Length")
			if headerContentSize != "" {
				expectedContentSize, _ = strconv.Atoi(headerContentSize)
			} else {
				expectedContentSize = h.c.GetInt(CONTENT_SIZE)
			}

			if h.c.config.method != "HEAD" && int64(expectedContentSize) != contentSize {
				record.Error = &LengthError{invalidContnetSize}
				return
			}

		}

		sw.Stop()
	}()
	return asyncResult
}

func DetectHost(context *Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			TraceException(r)
		}
	}()

	client := NewClient(context.config)
	reqeust, err := NewHttpRequest(context.config)
	if err != nil {
		return
	}

	resp, err := client.Do(reqeust)

	if err != nil {
		return
	} else {

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		context.SetString(SERVER_NAME, resp.Header.Get("Server"))
		headerContentSize := resp.Header.Get("Content-Length")

		if headerContentSize != "" {
			contentSize, _ := strconv.Atoi(headerContentSize)
			context.SetInt(CONTENT_SIZE, contentSize)
		} else {
			context.SetInt(CONTENT_SIZE, len(body))
		}
	}

	return
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

	return &http.Client{Transport: transport}
}

func NewHttpRequest(config *Config) (request *http.Request, err error) {

	var body io.Reader

	if (config.method == "POST" || config.method == "PUT") && config.bodyContent != nil {
		body = bytes.NewReader(config.bodyContent)
	}

	request, err = http.NewRequest(config.method, config.url, body)

	if err != nil {
		return
	}

	if body != nil && config.contentType == "text/plain" {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		request.Header.Set("Content-Type", config.contentType)
	}

	request.Header.Set("User-Agent", config.userAgent)

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

	return
}

func CopyHttpRequest(config *Config, request *http.Request) *http.Request {
	if config.method == "POST" || config.method == "PUT" {
		newRequest := *request
		if newRequest.Body != nil {
			newRequest.Body = ioutil.NopCloser(bytes.NewReader(config.bodyContent))
		}
		return &newRequest
	} else {
		return request
	}
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

type ResponseTimeoutError struct {
	err error
}

func (e *ResponseTimeoutError) Error() string {
	return e.err.Error()
}
