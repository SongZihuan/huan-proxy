package server

import (
	"errors"
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile"
	"github.com/SongZihuan/huan-proxy/src/flagparser"
	"github.com/SongZihuan/huan-proxy/src/logger"
	"net/http"
)

var ServerStop = fmt.Errorf("server stop")

type HTTPServer struct {
	address string
	skip    map[string]struct{}
	isTerm  bool
	writer  func(msg string)
}

func NewServer() *HTTPServer {
	if !flagparser.IsReady() || !config.IsReady() {
		panic("not ready")
	}

	var skip = make(map[string]struct{}, 10)

	return &HTTPServer{
		address: config.GetConfig().Http.Address,
		skip:    skip,
		isTerm:  logger.IsInfoTermNotDumb(),
		writer:  logger.InfoWrite,
	}
}

func (s *HTTPServer) GetConfig() *config.YamlConfig {
	// 不用检查Ready，因为在NewServer的时候已经检查过了
	return config.GetConfig()
}

func (s *HTTPServer) GetRules() *rulescompile.RuleListCompileConfig {
	// 不用检查Ready，因为在NewServer的时候已经检查过了
	return config.GetRules()
}

func (s *HTTPServer) GetRulesList() []*rulescompile.RuleCompileConfig {
	// 不用检查Ready，因为在NewServer的时候已经检查过了
	return s.GetRules().Rules
}

func (s *HTTPServer) RunHttp() error {
	err := s.runHttp()
	if errors.Is(err, http.ErrServerClosed) {
		return ServerStop
	} else if err != nil {
		return err
	}

	return nil
}

func (s *HTTPServer) runHttp() error {
	logger.Infof("start server in %s", s.address)
	return http.ListenAndServe(s.address, s)
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.LoggerServerHTTP(w, r, s.NormalServeHTTP)
}
