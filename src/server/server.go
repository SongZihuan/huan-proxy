package server

import (
	"errors"
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config"
	"github.com/SongZihuan/huan-proxy/src/flagparser"
	"github.com/SongZihuan/huan-proxy/src/logger"
	"net/http"
)

var ServerStop = fmt.Errorf("server stop")

type HTTPServer struct {
	address string
	cfg     *config.ConfigStruct
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
		address: config.Config().Yaml.Http.Address,
		cfg:     config.Config(),
		skip:    skip,
		isTerm:  logger.IsInfoTermNotDumb(),
		writer:  logger.InfoWrite,
	}
}

func (s *HTTPServer) Run() error {
	err := s.run()
	if errors.Is(err, http.ErrServerClosed) {
		return ServerStop
	} else if err != nil {
		return err
	}

	return nil
}

func (s *HTTPServer) run() error {
	logger.Infof("start server in %s", s.address)
	return http.ListenAndServe(s.address, s)
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.LoggerServerHTTP(w, r, s.NormalServeHTTP)
}
