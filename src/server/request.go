package server

import (
	"github.com/SongZihuan/huan-proxy/src/utils"
	"net/http"
	"net/url"
)

type ProxyRequest struct {
	req *http.Request

	Host       string
	Proto      string
	IsTLS      bool
	Method     string
	RemoteAddr string
	URL        *url.URL
	Header     http.Header

	_host       string
	_proto      string
	_method     string
	_remoteAddr string
	_url        *url.URL
	_header     http.Header

	written bool
}

func NewRequest(req *http.Request) *ProxyRequest {
	return &ProxyRequest{
		req: req,

		Host:       req.Host,
		Proto:      req.Proto,
		IsTLS:      req.TLS != nil,
		Method:     req.Method,
		RemoteAddr: req.RemoteAddr,
		Header:     req.Header.Clone(),
		URL:        utils.URLClone(req.URL),

		_host:       req.Host,
		_proto:      req.Proto,
		_method:     req.Method,
		_remoteAddr: req.RemoteAddr,
		_url:        req.URL,
		_header:     req.Header,

		written: false,
	}
}

func (r *ProxyRequest) ResetHttpRequest() error {
	r.req.Host = r._host
	r.req.Proto = r._proto
	r.req.Method = r._method
	r.req.RemoteAddr = r._remoteAddr
	r.req.Header = r._header
	r.req.URL = r._url
	return nil
}

func (r *ProxyRequest) Reset() error {
	err := r.ResetHttpRequest()
	if err != nil {
		return err
	}

	r.Host = r.req.Host
	r.Proto = r.req.Proto
	r.Method = r.req.Method
	r.RemoteAddr = r.req.RemoteAddr
	r.Header = r.req.Header.Clone()
	r.URL = utils.URLClone(r.req.URL)
	r.written = false
	return nil
}

func (r *ProxyRequest) WriteToHttpRRequest() (req *http.Request, err error) {
	if r.written {
		return r.req, nil
	}

	defer func() {
		if err != nil {
			_ = r.ResetHttpRequest() // 复原所有操作
			r.written = false
		}
	}()

	r.req.Host = r.Host
	r.req.Proto = r.Proto
	r.req.Method = r.Method
	r.req.RemoteAddr = r.RemoteAddr
	r.req.URL = utils.URLClone(r.URL)
	r.req.Header = r.Header.Clone()
	r.written = true
	return r.req, nil
}
