package server

import (
	"github.com/SongZihuan/huan-proxy/src/config"
	"github.com/SongZihuan/huan-proxy/src/flagparser"
	"github.com/SongZihuan/huan-proxy/src/server/core"
	"github.com/SongZihuan/huan-proxy/src/server/httpserver"
	"github.com/SongZihuan/huan-proxy/src/server/httpsserver"
	"github.com/SongZihuan/huan-proxy/src/server/middleware/loggerserver"
)

type HuanProxyServer struct {
	http   *httpserver.HTTPServer
	https  *httpsserver.HTTPSServer
	logger *loggerserver.LogServer
	server *core.CoreServer
}

func NewHuanProxyServer() *HuanProxyServer {
	if !flagparser.IsReady() || !config.IsReady() {
		panic("not ready")
	}

	logger := loggerserver.NewLogServer()
	server := core.NewCoreServer(logger)

	res := &HuanProxyServer{
		logger: loggerserver.NewLogServer(),
		server: server,
		http:   httpserver.NewHTTPServer(server),
		https:  httpsserver.NewHTTPSServer(server),
	}

	return res
}

func (s *HuanProxyServer) Run(httpErrorChan chan error, httpsErrorChan chan error) (err error) {
	if s.https != nil {
		err := s.https.LoadHttps()
		if err != nil {
			return err
		}

		s.https.RunHttps(httpsErrorChan)
	}

	if s.http != nil {
		err := s.http.LoadHttp()
		if err != nil {
			return err
		}

		s.http.RunHttp(httpErrorChan)
	}

	return nil
}

func (s *HuanProxyServer) Stop() (err error) {
	if s.http != nil {
		_ = s.http.StopHttp()
	}

	if s.https != nil {
		_ = s.https.StopHttps()
	}

	return nil
}
