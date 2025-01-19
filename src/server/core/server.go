package core

import (
	"errors"
	"fmt"
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

type TestWriter struct {
	http.ResponseWriter
}

func (t *TestWriter) WriteHeader(statusCode int) {
	fmt.Println("WRITE HEADER CALL")
	t.ResponseWriter.WriteHeader(statusCode)
}

func (c *CoreServer) ServeHTTP(_w http.ResponseWriter, r *http.Request) {
	w := &TestWriter{
		ResponseWriter: _w,
	}

	writer := responsewriter.NewResponseWriter(w)
	writer.WriteHeader(http.StatusOK)
	//
	//c.logger.ServeHTTP(writer, r, c.CoreServeHTTP)
	//
	//fmt.Println("TAG 2")
	err := writer.WriteToResponse()
	if err != nil && !errors.Is(err, responsewriter.ErrHasWriter) {
		fmt.Printf("Err: %s", err.Error())
		//writer.ServerError()
	}
}
