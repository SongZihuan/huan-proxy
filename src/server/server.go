package server

import (
	"errors"
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config"
	"github.com/SongZihuan/huan-proxy/src/flagparser"
	"github.com/SongZihuan/huan-proxy/src/logger"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"net/http"
	"strings"
)

var ServerStop = fmt.Errorf("server stop")

type HTTPServer struct {
	address string
	cfg     *config.ConfigStruct
}

func NewServer() *HTTPServer {
	if !flagparser.IsReady() || !config.IsReady() {
		panic("not ready")
	}

	return &HTTPServer{
		address: config.Config().Yaml.Http.Address,
		cfg:     config.Config(),
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
	s.writeHuanProxyHeader(w)
	if !s.checkProxyTrust(w, r) {
		return
	}

	func() {
		for index, rule := range s.cfg.Yaml.Rules.Rules {
			if rule.Type == config.ProxyTypeFile {
				url := utils.ProcessPath(r.URL.Path)
				logger.Tagf("A [%s] [%s]", url, rule.BasePath)
				if url == rule.BasePath {
					if s.corsHandler(w, r) {
						s.fileServer(rule, w, r)
					}
					return
				}
			} else if rule.Type == config.ProxyTypeFile {
				if r.Method == http.MethodGet {
					urlpath := utils.ProcessPath(r.URL.Path)
					if urlpath == rule.BasePath || strings.HasPrefix(urlpath, rule.BasePath+"/") {
						if s.corsHandler(w, r) {
							s.dirServer(rule, w, r)
						}
						return
					}
				}
			} else if rule.File == config.ProxyTypeAPI {
				urlpath := utils.ProcessPath(r.URL.Path)
				if urlpath == rule.BasePath || strings.HasPrefix(urlpath, rule.BasePath+"/") {
					s.apiServer(index, rule, w, r)
					return
				}
			} else {
				s.abortNotFound(w)
			}
		}

		s.abortNotFound(w)
	}()
}
