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

func (s *HuanProxyServer) dirServer(ctx *Context) {
	if !s.cors(ctx.Rule.Dir.Cors, ctx) {
		return
	}

	if ctx.Request.Method() != http.MethodGet {
		s.abortMethodNotAllowed(ctx)
		return
	}

	dirBasePath := ctx.Rule.Dir.BasePath // 根部目录
	fileAccess := ""                     // 访问目录
	filePath := ""                       // 根部目录+访问目录=实际目录

	url := utils.ProcessURLPath(ctx.Request.URL().Path)
	if ctx.Rule.MatchType == matchcompile.RegexMatch {
		fileAccess = s.dirRewrite("", ctx.Rule.Dir.AddPath, ctx.Rule.Dir.SubPath, ctx.Rule.Dir.Rewrite)
		filePath = path.Join(dirBasePath, fileAccess)
	} else {
		if url == ctx.Rule.MatchPath {
			fileAccess = s.dirRewrite("", ctx.Rule.Dir.AddPath, ctx.Rule.Dir.SubPath, ctx.Rule.Dir.Rewrite)
			filePath = path.Join(dirBasePath, fileAccess)
		} else if strings.HasPrefix(url, ctx.Rule.MatchPath+"/") {
			fileAccess = s.dirRewrite(url[len(ctx.Rule.MatchPath+"/"):], ctx.Rule.Dir.AddPath, ctx.Rule.Dir.SubPath, ctx.Rule.Dir.Rewrite)
			filePath = path.Join(dirBasePath, fileAccess)
		} else {
			s.abortNotFound(ctx)
			return
		}
	}

	if filePath == "" {
		s.abortNotFound(ctx) // 正常清空不会走到这个流程
		return
	}

	if utils.IsFile(filePath) {
		// 判断这个文件是否被ignore，因为ignore是从dirBasePath写起，也可以是完整路径，因此filePath和fileAccess都要判断
		for _, ignore := range ctx.Rule.Dir.IgnoreFile {
			if ignore.CheckName(fileAccess) || ignore.CheckName(filePath) {
				s.abortNotFound(ctx)
				return
			}
		}
	} else {
		filePath = s.getIndexFile(ctx.Rule, filePath)
		if filePath == "" || !utils.IsFile(filePath) {
			s.abortNotFound(ctx)
			return
		}
	}

	if !utils.CheckIfSubPath(dirBasePath, filePath) {
		s.abortForbidden(ctx)
		return
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		s.abortNotFound(ctx)
		return
	}

	mimeType := mimetype.Detect(file)
	accept := ctx.Request.Header().Get("Accept")
	if !utils.AcceptMimeType(accept, mimeType.String()) {
		s.abortNotAcceptable(ctx)
		return
	}

	_, err = ctx.Writer.Write(file)
	if err != nil {
		s.abortServerError(ctx)
		return
	}
	ctx.Writer.Header().Set("Content-Type", mimeType.String())
	s.statusOK(ctx)
}

func (s *HuanProxyServer) dirRewrite(srcpath string, prefix string, suffix string, rewrite *rewritecompile.RewriteCompileConfig) string {
	if strings.HasPrefix(srcpath, suffix) {
		srcpath = srcpath[len(suffix):]
	}

	srcpath = path.Join(prefix, srcpath)

	if rewrite.Use && rewrite.Regex != nil {
		rewrite.Regex.ReplaceAllString(srcpath, rewrite.Target)
	}

	return srcpath
}

func (s *HuanProxyServer) getIndexFile(rule *rulescompile.RuleCompileConfig, dir string) string {
	return s._getIndexFile(rule, dir, "", IndexMaxDeep)
}

func (s *HuanProxyServer) _getIndexFile(rule *rulescompile.RuleCompileConfig, baseDir string, nextDir string, deep int) string {
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
