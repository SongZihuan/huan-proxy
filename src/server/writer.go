package server

import (
	"bytes"
	"net/http"
)

type writer http.ResponseWriter

type ResponseWriter struct {
	writer
	status  int
	buffer  bytes.Buffer
	written bool
	header  http.Header
}

func NewWriter(w writer) *ResponseWriter {
	res := &ResponseWriter{
		writer: w,
		status: 0,
		header: make(http.Header, 10),
	}

	for n, h := range w.Header() {
		nh := make([]string, 0, len(h))
		copy(nh, h)
		res.header[n] = nh
	}

	return res
}

func (r *ResponseWriter) Size() int64 {
	return int64(r.buffer.Len())
}

func (r *ResponseWriter) Write(p []byte) (int, error) {
	return r.buffer.Write(p)
}

func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.status = statusCode
}

func (r *ResponseWriter) Header() http.Header {
	return r.header
}

func (r *ResponseWriter) Reset() {
	r.status = 0
	r.header = make(http.Header, 10)
	r.buffer.Reset()
	r.written = false
}

func (r *ResponseWriter) WriteToResponse() error {
	if r.written {
		return nil
	}

	_, err := r.writer.Write(r.buffer.Bytes())
	if err != nil {
		return err
	}

	r.writer.WriteHeader(r.status)

	writerHeader := r.writer.Header()
	for n, h := range r.header {
		nh := make([]string, 0, len(h))
		copy(nh, h)
		writerHeader[n] = nh
	}

	delHeader := make([]string, 0, 10)
	for n, _ := range writerHeader {
		if _, ok := r.header[n]; !ok {
			delHeader = append(delHeader, n)
		}
	}

	for _, n := range delHeader {
		delete(writerHeader, n)
	}

	r.written = true
	return nil
}
