package readonlyrequest

import (
	"github.com/SongZihuan/huan-proxy/src/utils"
	"net/http"
	"net/url"
)

type ReadOnlyRequest struct {
	req    *http.Request
	url    *url.URL
	header http.Header
}

func NewReadOnlyRequest(req *http.Request) *ReadOnlyRequest {
	return &ReadOnlyRequest{
		req:    req,
		url:    utils.URLClone(req.URL),
		header: req.Header.Clone(),
	}
}

func (r *ReadOnlyRequest) Host() string {
	return r.req.Host
}

func (r *ReadOnlyRequest) Method() string {
	return r.req.Method
}

func (r *ReadOnlyRequest) RemoteAddr() string {
	return r.req.Host
}

func (r *ReadOnlyRequest) Proto() string {
	return r.req.Proto
}

func (r *ReadOnlyRequest) MustProto() string {
	proto := r.req.Proto
	if proto == "" {
		if r.IsTLS() {
			return "https"
		} else {
			return "http"
		}
	} else {
		return proto
	}
}

func (r *ReadOnlyRequest) URL() *url.URL {
	return utils.URLClone(r.req.URL)
}

func (r *ReadOnlyRequest) Header() http.Header {
	return r.req.Header.Clone()
}

func (r *ReadOnlyRequest) IsTLS() bool {
	return r.req.TLS != nil
}
