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
const DefaultIgnoreFileMap = 20

func (s *HTTPServer) dirServer(ruleIndex int, rule *config.ProxyConfig, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.abortMethodNotAllowed(w)
		return
	}

	dirBasePath := rule.Dir
	fileBase := ""
	filePath := ""

	url := utils.ProcessPath(r.URL.Path)
	if url == rule.BasePath {
		filePath = dirBasePath
		fileBase = ""
	} else if strings.HasPrefix(url, rule.BasePath+"/") {
		fileBase = url[len(rule.BasePath+"/"):]
		filePath = path.Join(dirBasePath, fileBase)
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
		fileBase = filePath[len(rule.BasePath+"/"):len(filePath)]
	} else if fileBase != "" {
		ignore, err := s.cfg.IgnoreFile.ForEach(ruleIndex, func(file *config.IgnoreFileCompile) (any, error) {
			if file.CheckName(fileBase) {
				return true, nil
			}
			return nil, nil
		})
		if err != nil {
			s.abortNotFound(w)
			return
		} else if ig, ok := ignore.(bool); ok && ig {
			filePath = s.getIndexFile(ruleIndex, filePath)
			fileBase = filePath[len(rule.BasePath+"/"):len(filePath)]
		}
	}

	// 接下来的部分不在使用fileBase

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
	s.statusOK(w)
}

func (s *HTTPServer) getIndexFile(ruleIndex int, dir string) string {
	return s._getIndexFile(ruleIndex, dir, "", IndexMaxDeep)
}

func (s *HTTPServer) _getIndexFile(ruleIndex int, baseDir string, nextDir string, deep int) string {
	if deep == 0 {
		return ""
	}

	dir := path.Join(baseDir, nextDir)
	if !utils.IsDir(dir) {
		return ""
	}

	lst, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	var ignoreFileMap = make(map[string]bool, DefaultIgnoreFileMap)

	_, err = s.cfg.IgnoreFile.ForEach(ruleIndex, func(file *config.IgnoreFileCompile) (any, error) {
		for _, i := range lst {
			fileName := path.Join(nextDir, i.Name())
			if file.CheckName(fileName) {
				ignoreFileMap[fileName] = true
			}
		}
		return nil, nil
	})
	if err != nil {
		return ""
	}

	var indexDirNum = -1
	var indexDir os.DirEntry = nil

	var indexFileNum = -1
	var indexFile os.DirEntry = nil

	_, err = s.cfg.IndexFile.ForEach(ruleIndex, func(file *config.IndexFileCompile) (any, error) {
		for _, i := range lst {
			fileName := path.Join(nextDir, i.Name())

			if _, ok := ignoreFileMap[fileName]; ok {
				continue
			}

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
		return s._getIndexFile(ruleIndex, dir, indexDir.Name(), deep-1)
	} else {
		return ""
	}
}
