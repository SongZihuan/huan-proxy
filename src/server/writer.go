package server

import (
	"bytes"
	"fmt"
	"net/http"
)

var ErrHasWriter = fmt.Errorf("ResponseWriter has been written")

type ResponseWriter struct {
	writer http.ResponseWriter

	status int
	buffer bytes.Buffer
	header http.Header

	written bool
}

func NewResponseWriter(w http.ResponseWriter) http.ResponseWriter {
	if _, ok := w.(*ResponseWriter); ok {
		return w
	}

	return &ResponseWriter{
		writer: w,

		status: 0,
		header: w.Header().Clone(),

		written: false,
	}
}

func (r *ResponseWriter) Size() int64 {
	return int64(r.buffer.Len())
}

func (r *ResponseWriter) Status() int {
	return r.status
}

func (r *ResponseWriter) Write(p []byte) (int, error) {
	if r.written {
		return 0, ErrHasWriter
	}

	return r.buffer.Write(p)
}

func (r *ResponseWriter) WriteHeader(statusCode int) {
	if r.written {
		return
	}

	r.status = statusCode
}

func (r *ResponseWriter) Header() http.Header {
	if r.written {
		return nil
	}

	return r.header
}

func (r *ResponseWriter) Reset() error {
	if r.written {
		return ErrHasWriter
	}

	r.status = 0
	r.buffer.Reset()
	r.header = r.writer.Header()
	r.written = false

	return nil
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

func (r *ResponseWriter) MustWriteToResponse() {
	err := r.WriteToResponse()
	if err == nil {
		return
	}

	r.writer.WriteHeader(http.StatusInternalServerError)
}
