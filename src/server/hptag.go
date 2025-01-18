package server

import (
	"fmt"
	resource "github.com/SongZihuan/huan-proxy"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/apicompile"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"net/http"
	"strings"
)

const XHuanProxyHeaer = apicompile.XHuanProxyHeaer
const ViaHeader = apicompile.ViaHeader

func (s *HuanProxyServer) writeHuanProxyHeader(ctx *Context) {
	version := strings.TrimSpace(utils.StringToOnlyPrint(resource.Version))
	ctx.Writer.Header().Set(XHuanProxyHeaer, version)
	if ctx.ProxyRequest != nil {
		ctx.ProxyRequest.Header.Set(XHuanProxyHeaer, version)
	}
}

func (s *HuanProxyServer) writeViaHeader(ctx *Context) {
	info := fmt.Sprintf("%s %s", ctx.Request.MustProto(), ctx.Rule.Api.Via)

	reqHeader := ctx.Request.Header().Get(ViaHeader)
	if reqHeader == "" {
		reqHeader = info
	} else {
		reqHeader = fmt.Sprintf("%s, %s", reqHeader, info)
	}
	ctx.Request.Header().Set(ViaHeader, reqHeader)

	respHeader := ctx.Writer.Header().Get(ViaHeader)
	if respHeader == "" {
		respHeader = info
	} else if !strings.Contains(respHeader, info) {
		respHeader = fmt.Sprintf("%s, %s", respHeader, info)
	}
	ctx.Writer.Header().Set(ViaHeader, respHeader)
}

func (s *HuanProxyServer) statusOK(ctx *Context) {
	ctx.Writer.WriteHeader(http.StatusOK)
}

func (s *HuanProxyServer) statusRedirect(ctx *Context, url string, code int) {
	ctx.Redirect(url, code)
}
