package server

import (
	"github.com/SongZihuan/huan-proxy/src/config"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"github.com/gabriel-vasile/mimetype"
	"net/http"
	"os"
	"path"
	"strings"
)

const IndexMaxDeep = 5

func (s *HTTPServer) dirServer(rule *config.ProxyConfig, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.abortMethodNotAllowed(w)
		return
	}

	dirBasePath := rule.Dir
	filePath := ""

	url := utils.ProcessPath(r.URL.Path)
	if url == rule.BasePath {
		filePath = dirBasePath
	} else if strings.HasPrefix(url, rule.BasePath+"/") {
		filePath = path.Join(dirBasePath, url[len(rule.BasePath+"/"):])
	} else {
		s.abortNotFound(w)
		return
	}

	if filePath == "" {
		s.abortNotFound(w)
		return
	}

	if !utils.IsFile(filePath) {
		filePath = s.getIndexFile(filePath)
	}

	if filePath == "" || !utils.IsFile(filePath) {
		s.abortNotFound(w)
		return
	}

	file, err := os.ReadFile(filePath)
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
}

func (s *HTTPServer) getIndexFile(dir string) string {
	return s._getIndexFile(dir, IndexMaxDeep)
}

func (s *HTTPServer) _getIndexFile(dir string, deep int) string {
	if deep == 0 {
		return ""
	}

	if !utils.IsDir(dir) {
		return ""
	}

	lst, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	var nextDir os.DirEntry = nil
	var indexHTML os.DirEntry = nil
	var indexXML os.DirEntry = nil
	var index os.DirEntry = nil

	for _, i := range lst {
		if i.IsDir() {
			nextDir = i
		} else if i.Name() == "index.html" {
			indexHTML = i
		} else if i.Name() == "index.xml" {
			indexXML = i
		} else if strings.HasPrefix(i.Name(), "index.") {
			index = i
		}
	}

	if indexHTML != nil {
		return path.Join(dir, indexHTML.Name())
	} else if indexXML != nil {
		return path.Join(dir, indexXML.Name())
	} else if index != nil {
		return path.Join(dir, index.Name())
	} else if nextDir != nil {
		return s._getIndexFile(path.Join(dir, nextDir.Name()), deep-1)
	} else {
		return ""
	}
}
