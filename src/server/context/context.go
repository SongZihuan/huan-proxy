package context

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile"
	"github.com/SongZihuan/huan-proxy/src/server/request/proxyrequest"
	"github.com/SongZihuan/huan-proxy/src/server/request/readonlyrequest"
	"github.com/SongZihuan/huan-proxy/src/server/responsewriter"
	"net/http"
)

type Context struct {
	Abort        bool
	resp         http.ResponseWriter
	req          *http.Request
	Writer       *responsewriter.ResponseWriter
	Request      *readonlyrequest.ReadOnlyRequest
	ProxyRequest *proxyrequest.ProxyRequest
	Rule         *rulescompile.RuleCompileConfig
}

func NewContext(rule *rulescompile.RuleCompileConfig, w http.ResponseWriter, r *http.Request) *Context {
	var proxyRequest *proxyrequest.ProxyRequest = nil
	if rule.Type == rulescompile.ProxyTypeAPI {
		proxyRequest = proxyrequest.NewRequest(r)
	}

	return &Context{
		resp:         w,
		req:          r,
		Writer:       responsewriter.NewResponseWriter(w),
		Request:      readonlyrequest.NewReadOnlyRequest(r),
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
	fmt.Println("TAG 1")
	err := ctx.Writer.WriteToResponse()
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Context) MustWriteToResponse() {
	fmt.Println("TAG 5")
	ctx.Writer.MustWriteToResponse()
}

func (ctx *Context) Reset() error {
	ctx.Abort = false

	_ = ctx.Writer.Reset()

	if ctx.ProxyRequest != nil {
		_ = ctx.ProxyRequest.Reset()
	}

	return nil
}

func (ctx *Context) StatusOK() {
	if ctx.Abort {
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
}

func (ctx *Context) Redirect(target string, code int) {
	if ctx.Abort {
		return
	}

	http.Redirect(ctx.Writer, ctx.req, target, code)
}
