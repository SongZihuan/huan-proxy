package server

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile"
	"net/http"
)

type Context struct {
	Abort        bool
	Writer       http.ResponseWriter
	Request      *ReadOnlyRequest
	ProxyRequest *ProxyRequest
	Rule         *rulescompile.RuleCompileConfig
}

func NewContext(rule *rulescompile.RuleCompileConfig, w http.ResponseWriter, r *http.Request) *Context {
	var proxyRequest *ProxyRequest = nil
	if rule.Type == rulescompile.ProxyTypeAPI {
		proxyRequest = NewRequest(r)
	}

	return &Context{
		Writer:       NewResponseWriter(w),
		Request:      NewReadOnlyRequest(r),
		ProxyRequest: proxyRequest,
		Rule:         rule,
	}
}

func (ctx *Context) ProxyWriteToHttpRRequest() (*http.Request, error) {
	if ctx.ProxyRequest == nil {
		return nil, fmt.Errorf("proxy request is nil")
	}

	req, err := ctx.ProxyRequest.WriteToHttpRRequest()
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (ctx *Context) WriteToResponse() error {
	w, ok := ctx.Writer.(*ResponseWriter)
	if !ok {
		return nil
	}
	err := w.WriteToResponse()
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Context) MustWriteToResponse() {
	if w, ok := ctx.Writer.(*ResponseWriter); ok {
		w.MustWriteToResponse()
	}
}

func (ctx *Context) Reset() error {
	ctx.Abort = false

	if w, ok := ctx.Writer.(*ResponseWriter); ok {
		_ = w.Reset()
	}

	if ctx.ProxyRequest != nil {
		_ = ctx.ProxyRequest.Reset()
	}

	return nil
}

func (ctx *Context) Redirect(target string, code int) {
	http.Redirect(ctx.Writer, ctx.Request.req, target, code)
}
