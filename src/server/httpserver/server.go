package httpserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config"
	"github.com/SongZihuan/huan-proxy/src/logger"
	"github.com/pires/go-proxyproto"
	"net"
	"net/http"
	"time"
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

func (s *HTTPServer) StopHttp() error {
	if s.server == nil {
		return nil
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancelFunc()

	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *HTTPServer) RunHttp(httpErrorChan chan error) {
	go func() {
		listener, err := s.getListener()
		if err != nil {
			httpErrorChan <- fmt.Errorf("listen fail")
			return
		}
		defer func() {
			_ = listener.Close()
		}()

		logger.Infof("start http server in %s", s.cfg.Address)
		err = s.server.Serve(listener)
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			httpErrorChan <- ServerStop
			return
		} else if err != nil {
			httpErrorChan <- err
			return
		}
	}()
}

func (s *HTTPServer) getListener() (net.Listener, error) {
	tcpListener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return nil, fmt.Errorf("tcp listen on %s: %s\n", s.server.Addr, err.Error())
	}

	var proxyListener net.Listener
	if s.cfg.ProxyProto.IsEnable(true) {
		proxyListener = &proxyproto.Listener{
			Listener:          tcpListener,
			ReadHeaderTimeout: 10 * time.Second,
		}
	} else {
		proxyListener = tcpListener
	}

	return proxyListener, nil
}
