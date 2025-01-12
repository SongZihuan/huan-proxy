package server

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	Status int
	Size   int64
}

func NewWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		Status:         0,
	}
}

func (r *ResponseWriter) Write(p []byte) (int, error) {
	n, err := r.ResponseWriter.Write(p)
	if err != nil {
		return n, err
	}
	r.Size += int64(n)
	return n, nil
}

func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.Status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
