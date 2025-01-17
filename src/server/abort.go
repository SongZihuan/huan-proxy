package server

import "net/http"

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
