package server

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/corscompile"
	"net/http"
)

func (s *HTTPServer) cors(corsRule *corscompile.CorsCompileConfig, w http.ResponseWriter, r *http.Request) bool {
	if corsRule.Ignore {
		if r.Method == http.MethodOptions {
			s.abortMethodNotAllowed(w)
			return false
		} else {
			return true
		}
	}

	origin := r.Header.Get("Origin")
	if origin == "" {
		s.abortForbidden(w)
		return false
	}

	if !corsRule.InOriginList(origin) {
		s.abortForbidden(w)
		return false
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", corsRule.MaxAgeSec))

	return true
}
