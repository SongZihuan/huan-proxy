package core

import (
	"github.com/SongZihuan/huan-proxy/src/server/context"
	"net/http"
)

func (c *CoreServer) abort(ctx *context.Context, code int) {
	if ctx.Abort {
		return
	}

	ctx.Abort = true
	err := ctx.Writer.Reset()
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
	} else {
		ctx.Writer.WriteHeader(code)
	}
}

func (c *CoreServer) abortForbidden(ctx *context.Context) {
	c.abort(ctx, http.StatusForbidden)
}

func (c *CoreServer) abortNotFound(ctx *context.Context) {
	c.abort(ctx, http.StatusNotFound)
}

func (c *CoreServer) abortNotAcceptable(ctx *context.Context) {
	c.abort(ctx, http.StatusNotAcceptable)
}

func (c *CoreServer) abortMethodNotAllowed(ctx *context.Context) {
	c.abort(ctx, http.StatusMethodNotAllowed)
}

func (c *CoreServer) abortServerError(ctx *context.Context) {
	c.abort(ctx, http.StatusInternalServerError)
}

func (c *CoreServer) abortNoContent(ctx *context.Context) {
	c.abort(ctx, http.StatusNoContent)
}
