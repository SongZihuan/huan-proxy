package server

import (
	"github.com/SongZihuan/huan-proxy/src/utils"
	"github.com/gabriel-vasile/mimetype"
	"net/http"
	"os"
)

func (s *HuanProxyServer) fileServer(ctx *Context) {
	if !s.cors(ctx.Rule.File.Cors, ctx) {
		return
	}

	if ctx.Request.Method() != http.MethodGet {
		s.abortMethodNotAllowed(ctx)
		return
	}

	file, err := os.ReadFile(ctx.Rule.File.Path)
	if err != nil {
		s.abortServerError(ctx)
		return
	}

	mimeType := mimetype.Detect(file)
	accept := ctx.Request.Header().Get("Accept")
	if !utils.AcceptMimeType(accept, mimeType.String()) {
		s.abortNotAcceptable(ctx)
		return
	}

	_, err = ctx.Writer.Write(file)
	if err != nil {
		s.abortServerError(ctx)
		return
	}
	ctx.Writer.Header().Set("Content-Type", mimeType.String())
	s.statusOK(ctx)
}
