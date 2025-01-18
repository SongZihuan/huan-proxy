package server

import "net/http"

func (s *HuanProxyServer) abort(ctx *Context, code int) {
	if ctx.Abort {
		return
	}

	ctx.Writer.WriteHeader(code)
	ctx.Abort = true
}

func (s *HuanProxyServer) abortForbidden(ctx *Context) {
	s.abort(ctx, http.StatusForbidden)
}

func (s *HuanProxyServer) abortNotFound(ctx *Context) {
	s.abort(ctx, http.StatusNotFound)
}

func (s *HuanProxyServer) abortNotAcceptable(ctx *Context) {
	s.abort(ctx, http.StatusNotAcceptable)
}

func (s *HuanProxyServer) abortMethodNotAllowed(ctx *Context) {
	s.abort(ctx, http.StatusMethodNotAllowed)
}

func (s *HuanProxyServer) abortServerError(ctx *Context) {
	s.abort(ctx, http.StatusInternalServerError)
}

func (s *HuanProxyServer) abortNoContent(ctx *Context) {
	s.abort(ctx, http.StatusNoContent)
}
