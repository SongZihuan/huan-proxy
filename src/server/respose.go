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

func (s *HTTPServer) writeHuanProxyHeader(r *http.Request) {
	version := strings.TrimSpace(utils.StringToOnlyPrint(resource.Version))
	h := r.Header.Get(XHuanProxyHeaer)
	if h == "" {
		h = version
	} else {
		h = fmt.Sprintf("%s, %s", h, version)
	}

	r.Header.Set(XHuanProxyHeaer, h)
}

func (s *HTTPServer) writeViaHeader(rule *rulescompile.RuleCompileConfig, r *http.Request) {
	info := fmt.Sprintf("%s %s", r.Proto, rule.Api.Via)

	h := r.Header.Get(ViaHeader)
	if h == "" {
		h = info
	} else {
		h = fmt.Sprintf("%s, %s", h, info)
	}

	r.Header.Set(ViaHeader, h)
}

func (s *HTTPServer) abortForbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte("Forbidden"))
}

func (s *HTTPServer) abortNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("NotFound"))
}

func (s *HTTPServer) abortNotAcceptable(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotAcceptable)
	_, _ = w.Write([]byte("NotAcceptable"))
}

func (s *HTTPServer) abortMethodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	_, _ = w.Write([]byte("MethodNotAllowed"))
}

func (s *HTTPServer) abortServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte("ServerError"))
}

func (s *HTTPServer) abortNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func (s *HTTPServer) statusOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}
