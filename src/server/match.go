package server

import (
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/matchcompile"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"net/http"
	"strings"
)

func (s *HuanProxyServer) matchURL(rule *rulescompile.RuleCompileConfig, r *http.Request) bool {
	url := utils.ProcessURLPath(r.URL.Path)
	if rule.MatchType == matchcompile.RegexMatch {
		if rule.MatchRegex.MatchString(url) || rule.MatchRegex.MatchString(url+"/") {
			return true
		}
	} else if rule.MatchType == matchcompile.PrefixMatch {
		path := utils.ProcessURLPath(rule.MatchPath)
		if url == path || strings.HasPrefix(url, path+"/") {
			return true
		}
	} else if rule.MatchType == matchcompile.PrecisionMatch {
		path := utils.ProcessURLPath(rule.MatchPath)
		if url == path {
			return true
		}
	}
	return false
}
