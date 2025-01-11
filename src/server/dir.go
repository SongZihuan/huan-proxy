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

func (s *HTTPServer) dirServer(ruleIndex int, rule *config.ProxyConfig, w http.ResponseWriter, r *http.Request) {
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
		filePath = s.getIndexFile(ruleIndex, filePath)
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

func (s *HTTPServer) getIndexFile(ruleIndex int, dir string) string {
	return s._getIndexFile(ruleIndex, dir, IndexMaxDeep)
}

func (s *HTTPServer) _getIndexFile(ruleIndex int, dir string, deep int) string {
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

	var indexDirNum = -1
	var indexDir os.DirEntry = nil

	var indexFileNum = -1
	var indexFile os.DirEntry = nil

	_, err = s.cfg.IndexFile.ForEach(ruleIndex, func(file *config.IndexFileCompile) (any, error) {
		for _, i := range lst {
			if file.CheckDirEntry(i) {
				if i.IsDir() {
					if indexDirNum == -1 || file.Index < indexDirNum {
						indexDirNum = file.Index
						indexDir = i
					}
				} else {
					if indexFileNum == -1 || file.Index < indexFileNum {
						indexFileNum = file.Index
						indexFile = i
					}
				}
			}
		}
		return nil, nil
	})
	if err != nil {
		return ""
	}

	if indexFile != nil {
		return path.Join(dir, indexFile.Name())
	} else if indexDir != nil {
		return s._getIndexFile(ruleIndex, path.Join(dir, indexDir.Name()), deep-1)
	} else {
		return ""
	}
}
