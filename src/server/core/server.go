package core

import (
	"errors"
	"github.com/SongZihuan/huan-proxy/src/config"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile"
	"github.com/SongZihuan/huan-proxy/src/flagparser"
	"github.com/SongZihuan/huan-proxy/src/server/responsewriter"
	"net/http"
)

type Middleware interface {
	ServeHTTP(http.ResponseWriter, *http.Request, http.HandlerFunc)
}

type CoreServer struct {
	logger Middleware
}

func NewCoreServer(logger Middleware) *CoreServer {
	if !flagparser.IsReady() || !config.IsReady() {
		panic("not ready")
	}
	return &CoreServer{
		logger: logger,
	}
}

func (c *CoreServer) GetConfig() *config.YamlConfig {
	// 不用检查Ready，因为在NewServer的时候已经检查过了
	return config.GetConfig()
}

func (c *CoreServer) GetRules() *rulescompile.RuleListCompileConfig {
	// 不用检查Ready，因为在NewServer的时候已经检查过了
	return config.GetRules()
}

func (c *CoreServer) GetRulesList() []*rulescompile.RuleCompileConfig {
	// 不用检查Ready，因为在NewServer的时候已经检查过了
	return c.GetRules().Rules
}

func (c *CoreServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writer := responsewriter.NewResponseWriter(w)

	c.logger.ServeHTTP(writer, r, c.CoreServeHTTP)

	err := writer.WriteToResponse()
	if err != nil && !errors.Is(err, responsewriter.ErrHasWriter) {
		writer.ServerError()
	}
}
