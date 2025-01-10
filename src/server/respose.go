package server

import (
	resource "github.com/SongZihuan/huan-proxy"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"net/http"
)

func (s *HTTPServer) writeHuanProxyHeader(w http.ResponseWriter) {
	w.Header().Set("HuanProxy", utils.StringToOnlyPrint(resource.Version))
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
