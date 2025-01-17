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

func (s *HuanProxyServer) writeHuanProxyHeader(r *http.Request) {
	version := strings.TrimSpace(utils.StringToOnlyPrint(resource.Version))
	h := r.Header.Get(XHuanProxyHeaer)
	if h == "" {
		h = version
	} else {
		h = fmt.Sprintf("%s, %s", h, version)
	}

	r.Header.Set(XHuanProxyHeaer, h)
}

func (s *HuanProxyServer) writeViaHeader(rule *rulescompile.RuleCompileConfig, r *http.Request) {
	info := fmt.Sprintf("%s %s", r.Proto, rule.Api.Via)

	h := r.Header.Get(ViaHeader)
	if h == "" {
		h = info
	} else {
		h = fmt.Sprintf("%s, %s", h, info)
	}

	r.Header.Set(ViaHeader, h)
}

func (s *HuanProxyServer) abortForbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
}

func (s *HuanProxyServer) abortNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func (s *HuanProxyServer) abortNotAcceptable(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotAcceptable)
}

func (s *HuanProxyServer) abortMethodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *HuanProxyServer) abortServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func (s *HuanProxyServer) abortNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func (s *HuanProxyServer) statusOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func (s *HuanProxyServer) statusRedirect(w http.ResponseWriter, r *http.Request, url string, code int) {
	http.Redirect(w, r, url, code)
}
