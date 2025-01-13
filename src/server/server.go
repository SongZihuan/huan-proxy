package server

import (
	"errors"
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config"
	"github.com/SongZihuan/huan-proxy/src/flagparser"
	"github.com/SongZihuan/huan-proxy/src/logger"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"github.com/mattn/go-isatty"
	"net/http"
	"os"
	"strings"
)

var ServerStop = fmt.Errorf("server stop")

type HTTPServer struct {
	address string
	cfg     *config.ConfigStruct
	skip    map[string]struct{}
	isTerm  bool
}

func NewServer() *HTTPServer {
	if !flagparser.IsReady() || !config.IsReady() {
		panic("not ready")
	}

	var skip = make(map[string]struct{}, 10)
	var isTerm = true
	var out = logger.InfoWriter()

	w, ok := out.(*os.File)
	if !ok {
		isTerm = false
	} else if !isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd()) { // 非终端
		isTerm = false
	} else if os.Getenv("TERM") == "dumb" {
		// TERM为dump表示终端为基础模式，不支持高级显示
		isTerm = false
	}

	return &HTTPServer{
		address: config.Config().Yaml.Http.Address,
		cfg:     config.Config(),
		skip:    skip,
		isTerm:  isTerm,
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

func (s *HTTPServer) NormalServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.writeHuanProxyHeader(r)
	if !s.checkProxyTrust(w, r) {
		return
	}

	func() {
		for ruleIndex, rule := range s.cfg.Yaml.Rules.Rules {
			if rule.Type == config.ProxyTypeFile {
				url := utils.ProcessPath(r.URL.Path)
				if url == rule.BasePath {
					if s.corsHandler(w, r) {
						s.fileServer(rule, w, r)
					}
					return
				}
			} else if rule.Type == config.ProxyTypeDir {
				if r.Method == http.MethodGet {
					urlpath := utils.ProcessPath(r.URL.Path)
					if urlpath == rule.BasePath || strings.HasPrefix(urlpath, rule.BasePath+"/") {
						if s.corsHandler(w, r) {
							s.dirServer(ruleIndex, rule, w, r)
						}
						return
					}
				}
			} else if rule.Type == config.ProxyTypeAPI {
				urlpath := utils.ProcessPath(r.URL.Path)
				if urlpath == rule.BasePath || strings.HasPrefix(urlpath, rule.BasePath+"/") {
					s.apiServer(ruleIndex, rule, w, r)
					return
				}
			}
		}

		s.abortNotFound(w)
	}()
}
