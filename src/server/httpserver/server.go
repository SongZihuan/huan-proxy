package httpserver

import (
	"errors"
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config"
	"github.com/SongZihuan/huan-proxy/src/logger"
	"net/http"
)

var ServerStop = fmt.Errorf("https server stop")

type HTTPServer struct {
	cfg     *config.HttpConfig
	server  *http.Server
	handler http.Handler
}

func NewHTTPServer(handler http.Handler) *HTTPServer {
	httpcfg := config.GetConfig().Http

	if httpcfg.Address == "" {
		return nil
	}

	return &HTTPServer{
		cfg:     &httpcfg,
		server:  nil,
		handler: handler,
	}
}

func (s *HTTPServer) LoadHttp() error {
	s.server = &http.Server{
		Addr:    s.cfg.Address,
		Handler: s.handler,
	}
	return nil
}

func (s *HTTPServer) RunHttp(_httpschan chan error) chan error {
	go func(httpschan chan error) {
		logger.Infof("start http server in %s", s.cfg.Address)
		err := s.server.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			httpschan <- ServerStop
			return
		} else if err != nil {
			httpschan <- err
			return
		}
	}(_httpschan)

	return _httpschan
}
