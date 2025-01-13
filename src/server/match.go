package server

import (
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/matchcompile"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"net/http"
	"strings"
)

func (s *HTTPServer) matchURL(rule *rulescompile.RuleCompileConfig, r *http.Request) bool {
	url := utils.ProcessPath(r.URL.Path)
	if rule.MatchType == matchcompile.RegexMatch {
		if rule.MatchRegex.MatchString(url) || rule.MatchRegex.MatchString(url+"/") {
			return true
		}
	} else if rule.MatchType == matchcompile.PrefixMatch {
		path := utils.ProcessPath(rule.MatchPath)
		if url == path || strings.HasPrefix(url, path+"/") {
			return true
		}
	} else if rule.MatchType == matchcompile.PrecisionMatch {
		path := utils.ProcessPath(rule.MatchPath)
		if url == path {
			return true
		}
	}
	return false
}
