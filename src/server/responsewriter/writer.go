package responsewriter

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

var ErrHasWriter = fmt.Errorf("ResponseWriter has been written")

type ResponseWriter struct {
	writer http.ResponseWriter

	status int
	buffer bytes.Buffer
	header http.Header
	size   int64

	writtenStatus bool
	writtenBody   bool
	writtenHeader bool
	written       bool
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	if writer, ok := w.(*ResponseWriter); ok {
		return writer
	}

	return &ResponseWriter{
		writer: w,

		status: 0,
		header: w.Header().Clone(),

		written: false,
	}
}

func (r *ResponseWriter) Size() int64 {
	return r.size
}

func (r *ResponseWriter) Status() int {
	return r.status
}

func (r *ResponseWriter) Write(p []byte) (int, error) {
	if r.written || r.writtenBody {
		return 0, ErrHasWriter
	}

	return r.buffer.Write(p)
}

func (r *ResponseWriter) WriteHeader(statusCode int) {
	if r.written {
		return
	}

	r.status = statusCode
	fmt.Printf("Set Status is: %d\n", r.status)
}

func (r *ResponseWriter) ServerError() {
	if r.written || r.writtenStatus {
		return
	}

	r.status = http.StatusInternalServerError
	r.writer.WriteHeader(r.status)
	r.writtenStatus = true
	r.written = true

	fmt.Printf("Server Error Status is: %d\n", r.status)
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
	r.header = r.writer.Header().Clone()
	r.written = false

	return nil
}

func (r *ResponseWriter) WriteToResponse() error {
	if r.written {
		return ErrHasWriter
	}

	// status 抓鬼太吗最先写入
	if !r.writtenStatus {
		r.writer.WriteHeader(r.status)
		r.writtenStatus = true
		fmt.Printf("Write Status is: %d\n", r.status)
	}

	if !r.writtenHeader {
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
		r.writtenHeader = false
	}

	if !r.writtenBody {
		_, err := r.writer.Write(r.buffer.Bytes())
		if err != nil {
			return err
		}
		r.size = int64(r.buffer.Len())
		r.buffer.Reset() // 清理
		r.writtenBody = true
	}

	r.written = true
	return nil
}

func (r *ResponseWriter) MustWriteToResponse() {
	err := r.WriteToResponse()
	if err == nil || errors.Is(err, ErrHasWriter) {
		return
	}

	r.ServerError()
}
