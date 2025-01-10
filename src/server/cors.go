package server

import (
	"fmt"
	"net/http"
)

func (s *HTTPServer) corsHandler(w http.ResponseWriter, r *http.Request) bool {
	if s.cfg.Yaml.Http.Cors.Disable() {
		if r.Method == http.MethodOptions {
			s.abortNoContent(w)
		}
		return true
	}

	origin := r.Header.Get("Origin")
	if origin == "" {
		s.abortForbidden(w)
		return false
	}

	if !s.cfg.CoreOrigin.InOriginList(origin) {
		s.abortForbidden(w)
		return false
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", s.cfg.Yaml.Http.Cors.MaxAgeSec))

	return true
}
