package server

import (
	"fmt"
	resource "github.com/SongZihuan/huan-proxy"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/apicompile"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"net/http"
	"strings"
)

const XHuanProxyHeaer = apicompile.XHuanProxyHeaer
const ViaHeader = apicompile.ViaHeader

func (s *HuanProxyServer) writeHuanProxyHeader(w http.ResponseWriter, r *http.Request) {
	version := strings.TrimSpace(utils.StringToOnlyPrint(resource.Version))
	r.Header.Set(XHuanProxyHeaer, version)
	w.Header().Set(XHuanProxyHeaer, version)
}

func (s *HuanProxyServer) writeViaHeader(rule *rulescompile.RuleCompileConfig, w http.ResponseWriter, r *http.Request) {
	info := fmt.Sprintf("%s %s", r.Proto, rule.Api.Via)

	reqHeader := r.Header.Get(ViaHeader)
	if reqHeader == "" {
		reqHeader = info
	} else {
		reqHeader = fmt.Sprintf("%s, %s", reqHeader, info)
	}
	r.Header.Set(ViaHeader, reqHeader)

	respHeader := w.Header().Get(ViaHeader)
	if respHeader == "" {
		respHeader = info
	} else if !strings.Contains(respHeader, info) {
		respHeader = fmt.Sprintf("%s, %s", respHeader, info)
	}
	w.Header().Set(ViaHeader, respHeader)
}

func (s *HuanProxyServer) statusOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func (s *HuanProxyServer) statusRedirect(w http.ResponseWriter, r *http.Request, url string, code int) {
	http.Redirect(w, r, url, code)
}
