package server

import (
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/rewritecompile"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/matchcompile"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"github.com/gabriel-vasile/mimetype"
	"net/http"
	"os"
	"path"
	"strings"
)

const IndexMaxDeep = 5
const DefaultIgnoreFileMap = 20

func (s *HTTPServer) dirServer(rule *rulescompile.RuleCompileConfig, w http.ResponseWriter, r *http.Request) {
	if !s.cors(rule.File.Cors, w, r) {
		return
	}

	if r.Method != http.MethodGet {
		s.abortMethodNotAllowed(w)
		return
	}

	dirBasePath := rule.Dir.Dir // 根部目录
	fileAccess := ""            // 访问目录
	filePath := ""              // 根部目录+访问目录=实际目录

	url := utils.ProcessPath(r.URL.Path)
	if rule.MatchType == matchcompile.RegexMatch {
		fileAccess = s.rewrite("", rule.Dir.AddPrefixPath, rule.Dir.SubPrefixPath, rule.Dir.Rewrite)
		filePath = path.Join(dirBasePath, fileAccess)
	} else {
		if url == rule.MatchPath {
			fileAccess = s.rewrite("", rule.Dir.AddPrefixPath, rule.Dir.SubPrefixPath, rule.Dir.Rewrite)
			filePath = path.Join(dirBasePath, fileAccess)
		} else if strings.HasPrefix(url, rule.MatchPath+"/") {
			fileAccess = s.rewrite(url[len(rule.MatchPath+"/"):], rule.Dir.AddPrefixPath, rule.Dir.SubPrefixPath, rule.Dir.Rewrite)
			filePath = path.Join(dirBasePath, fileAccess)
		} else {
			s.abortNotFound(w)
			return
		}
	}

	if filePath == "" {
		s.abortNotFound(w) // 正常清空不会走到这个流程
		return
	}

	if utils.IsFile(filePath) {
		// 判断这个文件是否被ignore，因为ignore是从dirBasePath写起，也可以是完整路径，因此filePath和fileAccess都要判断
		for _, ignore := range rule.Dir.IgnoreFile {
			if ignore.CheckName(fileAccess) || ignore.CheckName(filePath) {
				s.abortNotFound(w)
				return
			}
		}
	} else {
		filePath = s.getIndexFile(rule, filePath)
		if filePath == "" || !utils.IsFile(filePath) {
			s.abortNotFound(w)
			return
		}
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		s.abortNotFound(w)
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

func (s *HTTPServer) rewrite(srcpath string, prefix string, suffix string, rewrite *rewritecompile.RewriteCompileConfig) string {
	srcpath = utils.ProcessPath(srcpath)
	prefix = utils.ProcessPath(prefix)
	suffix = utils.ProcessPath(suffix)

	if strings.HasPrefix(srcpath, suffix) {
		srcpath = srcpath[len(suffix):]
	}

	srcpath = prefix + srcpath

	if rewrite.Use && rewrite.Regex != nil {
		rewrite.Regex.ReplaceAllString(srcpath, rewrite.Target)
	}

	return srcpath
}

func (s *HTTPServer) getIndexFile(rule *rulescompile.RuleCompileConfig, dir string) string {
	return s._getIndexFile(rule, dir, "", IndexMaxDeep)
}

func (s *HTTPServer) _getIndexFile(rule *rulescompile.RuleCompileConfig, baseDir string, nextDir string, deep int) string {
	if deep == 0 {
		return ""
	}

	dir := path.Join(baseDir, nextDir)
	if !utils.IsDir(dir) {
		return ""
	}

	fileList, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	var ignoreFileMap = make(map[string]bool, DefaultIgnoreFileMap)
	for _, ignore := range rule.Dir.IgnoreFile {
		for _, file := range fileList {
			fileName := path.Join(nextDir, file.Name())
			if ignore.CheckName(fileName) {
				ignoreFileMap[fileName] = true
			}
		}
	}

	var indexDirNum = -1
	var indexDir os.DirEntry = nil

	var indexFileNum = -1
	var indexFile os.DirEntry = nil

	for indexID, index := range rule.Dir.IndexFile {
		for _, file := range fileList {
			fileName := path.Join(nextDir, file.Name())

			if _, ok := ignoreFileMap[fileName]; ok {
				continue
			}

			if index.CheckName(file.Name()) {
				if file.IsDir() {
					if indexDirNum == -1 || indexID < indexDirNum {
						indexDirNum = indexID
						indexDir = file
					}
				} else {
					if indexFileNum == -1 || indexID < indexFileNum {
						indexFileNum = indexID
						indexFile = file
					}
				}
			}
		}
	}

	if indexFile != nil {
		return path.Join(dir, indexFile.Name())
	} else if indexDir != nil {
		return s._getIndexFile(rule, dir, indexDir.Name(), deep-1)
	} else {
		return ""
	}
}
