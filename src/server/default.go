package server

import "net/http"

func (s *HuanProxyServer) defaultResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}
