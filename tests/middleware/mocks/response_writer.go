package mocks

import (
	"bufio"
	"net"
	"net/http"
)

// ResponseWriterHijacker is a mock of http.ResponseWriter that implements Hijacker
type ResponseWriterHijacker struct {
	http.ResponseWriter
	Hijacker
}

// Hijacker is a mock hijacker
type Hijacker struct {
	HijackFunc func() (net.Conn, *bufio.ReadWriter, error)
}

// Hijack implements http.Hijacker
func (h Hijacker) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h.HijackFunc != nil {
		return h.HijackFunc()
	}
	return nil, nil, nil
}

// NewResponseWriterHijacker creates a new ResponseWriterHijacker
func NewResponseWriterHijacker() *ResponseWriterHijacker {
	return &ResponseWriterHijacker{
		ResponseWriter: NewResponseWriter(),
		Hijacker:       Hijacker{},
	}
}

// ResponseWriter is a mock of http.ResponseWriter
type ResponseWriter struct {
	HeaderFunc      func() http.Header
	WriteFunc       func([]byte) (int, error)
	WriteHeaderFunc func(statusCode int)

	header     http.Header
	statusCode int
	body       []byte
}

// NewResponseWriter creates a new ResponseWriter
func NewResponseWriter() *ResponseWriter {
	rw := &ResponseWriter{
		header: make(http.Header),
	}
	rw.HeaderFunc = func() http.Header {
		return rw.header
	}
	rw.WriteFunc = func(b []byte) (int, error) {
		rw.body = append(rw.body, b...)
		return len(b), nil
	}
	rw.WriteHeaderFunc = func(statusCode int) {
		rw.statusCode = statusCode
	}
	return rw
}

// Header implements http.ResponseWriter
func (rw *ResponseWriter) Header() http.Header {
	return rw.HeaderFunc()
}

// Write implements http.ResponseWriter
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	return rw.WriteFunc(b)
}

// WriteHeader implements http.ResponseWriter
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.WriteHeaderFunc(statusCode)
}
