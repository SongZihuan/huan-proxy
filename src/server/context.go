package server

import (
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile"
	"net/http"
)

type Context struct {
	Abort   bool
	Writer  writer
	Request *http.Request
	Rule    *rulescompile.RuleCompileConfig
}

func NewContext(rule *rulescompile.RuleCompileConfig, w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  NewWriter(w),
		Request: r,
		Rule:    rule,
	}
}
