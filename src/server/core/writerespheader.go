package core

import "github.com/SongZihuan/huan-proxy/src/server/context"

func (c *CoreServer) WriteRespHeader(ctx *context.Context) {
	writerHeader := ctx.Writer.Header()

	for _, h := range ctx.Rule.RespHeader.HeaderSet {
		writerHeader.Set(h.Header, h.Value)
	}

	for _, h := range ctx.Rule.RespHeader.HeaderAdd {
		writerHeader.Add(h.Header, h.Value)
	}

	for _, h := range ctx.Rule.RespHeader.HeaderDel {
		writerHeader.Del(h.Header)
	}
}
