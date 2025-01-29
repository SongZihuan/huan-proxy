package core

import (
	"github.com/SongZihuan/huan-proxy/src/server/context"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"github.com/gabriel-vasile/mimetype"
	"net/http"
	"os"
)

func (c *CoreServer) fileServer(ctx *context.Context) {
	if !c.cors(ctx.Rule.File.Cors, ctx) {
		return
	}

	if ctx.Request.Method() != http.MethodGet {
		c.abortMethodNotAllowed(ctx)
		return
	}

	file, err := os.ReadFile(ctx.Rule.File.Path)
	if err != nil {
		c.abortServerError(ctx)
		return
	}

	mimeType := mimetype.Detect(file)
	accept := ctx.Request.Header().Get("Accept")
	if !utils.AcceptMimeType(accept, mimeType.String()) {
		c.abortNotAcceptable(ctx)
		return
	}

	_, err = ctx.Writer.Write(file)
	if err != nil {
		c.abortServerError(ctx)
		return
	}
	ctx.Writer.Header().Set("Content-Type", mimeType.String())
	c.statusOK(ctx)
}
