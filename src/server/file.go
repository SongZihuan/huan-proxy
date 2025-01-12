package server

import (
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"github.com/gabriel-vasile/mimetype"
	"net/http"
	"os"
)

func (s *HTTPServer) fileServer(rule *rulescompile.RuleCompileConfig, w http.ResponseWriter, r *http.Request) {
	if !s.cors(rule.File.Cors, w, r) {
		return
	}

	if r.Method != http.MethodGet {
		s.abortMethodNotAllowed(w)
		return
	}

	file, err := os.ReadFile(rule.File.Path)
	if err != nil {
		s.abortServerError(w)
		return
	}

	mimeType := mimetype.Detect(file)
	accept := r.Header.Get("Accept")
	if !utils.AcceptMimeType(accept, mimeType.String()) {
		s.abortNotAcceptable(w)
		return
	}

	_, err = w.Write(file)
	if err != nil {
		s.abortServerError(w)
		return
	}
	w.Header().Set("Content-Type", mimeType.String())
	s.statusOK(w)
}
