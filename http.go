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
	FieldServerName  = "ServerName"
	FieldContentSize = "ContentSize"
	MaxBufferSize    = 8192
)

var (
	ErrInvalidContnetSize = errors.New("invalid content size")
)

type HTTPWorker struct {
	c         *Context
	client    *http.Client
	jobs      chan *http.Request
	collector chan *Record
	discard   io.ReaderFrom
}

func NewHTTPWorker(context *Context, jobs chan *http.Request, collector chan *Record) *HTTPWorker {

	var buf []byte
	contentSize := context.GetInt(FieldContentSize)
	if contentSize < MaxBufferSize {
		buf = make([]byte, contentSize)
	} else {
		buf = make([]byte, MaxBufferSize)
	}

	return &HTTPWorker{
		context,
		NewClient(context.config),
		jobs,
		collector,
		&Discard{buf},
	}
}

func (h *HTTPWorker) Run() {
	h.c.start.Done()
	h.c.start.Wait()

	timer := time.NewTimer(h.c.config.executionTimeout)

	for job := range h.jobs {

		timer.Reset(h.c.config.executionTimeout)
		asyncResult := h.send(job)

		select {
		case record := <-asyncResult:
			h.collector <- record

		case <-timer.C:
			h.collector <- &Record{Error: &ResponseTimeoutError{errors.New("execution timeout")}}
			h.client.Transport.(*http.Transport).CancelRequest(job)

		case <-h.c.stop:
			h.client.Transport.(*http.Transport).CancelRequest(job)
			timer.Stop()
			return
		}
	}
	timer.Stop()
}

func (h *HTTPWorker) send(request *http.Request) (asyncResult chan *Record) {

	asyncResult = make(chan *Record, 1)
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
		}

		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode > 300 {
			record.Error = &ResponseError{err}
			return
		}

		contentSize, err = h.discard.ReadFrom(resp.Body)
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				record.Error = &LengthError{ErrInvalidContnetSize}
				return
			}

			record.Error = &ReceiveError{err}
			return
		}

		sw.Stop()
	}()
	return asyncResult
}

type Discard struct {
	blackHole []byte
}

func (d *Discard) ReadFrom(r io.Reader) (n int64, err error) {
	readSize := 0
	for {
		readSize, err = r.Read(d.blackHole)
		n += int64(readSize)
		if err != nil {
			if err == io.EOF {
				return n, nil
			}
			return
		}
	}
}

func DetectHost(context *Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			TraceException(r)
		}
	}()

	client := NewClient(context.config)
	reqeust, err := NewHTTPRequest(context.config)
	if err != nil {
		return
	}

	resp, err := client.Do(reqeust)

	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	context.SetString(FieldServerName, resp.Header.Get("Server"))
	headerContentSize := resp.Header.Get("Content-Length")

	if headerContentSize != "" {
		contentSize, _ := strconv.Atoi(headerContentSize)
		context.SetInt(FieldContentSize, contentSize)
	} else {
		context.SetInt(FieldContentSize, len(body))
	}

	return
}

func NewClient(config *Config) *http.Client {

	// skip certification check for self-signed certificates
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	// TODO: tcp options
	// TODO: monitor tcp metrics
	transport := &http.Transport{
		DisableCompression: !config.gzip,
		DisableKeepAlives:  !config.keepAlive,
		TLSClientConfig:    tlsconfig,
	}

	return &http.Client{Transport: transport}
}

func NewHTTPRequest(config *Config) (request *http.Request, err error) {

	var body io.Reader

	if config.method == "POST" || config.method == "PUT" {
		body = bytes.NewReader(config.bodyContent)
	}

	request, err = http.NewRequest(config.method, config.url, body)

	if err != nil {
		return
	}

	request.Header.Set("Content-Type", config.contentType)
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

func CopyHTTPRequest(config *Config, request *http.Request) *http.Request {
	newRequest := *request
	if request.Body != nil {
		newRequest.Body = ioutil.NopCloser(bytes.NewReader(config.bodyContent))
	}
	return &newRequest
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
