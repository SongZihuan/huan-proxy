package server

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/certssl"
	"github.com/SongZihuan/huan-proxy/src/config"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile"
	"github.com/SongZihuan/huan-proxy/src/flagparser"
	"github.com/SongZihuan/huan-proxy/src/logger"
	"net/http"
	"sync"
	"time"
)

var ServerStop = fmt.Errorf("server stop")

type HTTPServer struct {
	cfg     *config.HttpConfig
	server  *http.Server
	handler http.Handler
}

type HTTPSServer struct {
	cfg         *config.HttpsConfig
	reloadMutex sync.Mutex
	key         crypto.PrivateKey
	cert        *x509.Certificate
	cacert      *x509.Certificate
	server      *http.Server
	handler     http.Handler
}

type LogServer struct {
	skip   map[string]struct{}
	isTerm bool
	writer func(msg string)
}

type HuanProxyServer struct {
	http  HTTPServer
	https HTTPSServer
	LogServer
}

func NewHuanProxyServer() *HuanProxyServer {
	if !flagparser.IsReady() || !config.IsReady() {
		panic("not ready")
	}

	skip := make(map[string]struct{}, 10)
	httpcfg := config.GetConfig().Http
	httpscfg := config.GetConfig().Https

	res := &HuanProxyServer{
		http: HTTPServer{
			cfg:    &httpcfg,
			server: nil,
		},

		https: HTTPSServer{
			cfg:    &httpscfg,
			server: nil,
		},

		LogServer: LogServer{
			skip:   skip,
			isTerm: logger.IsInfoTermNotDumb(),
			writer: logger.InfoWrite,
		},
	}

	res.http.handler = res
	res.https.handler = res

	return res
}

func (s *HuanProxyServer) GetConfig() *config.YamlConfig {
	// 不用检查Ready，因为在NewServer的时候已经检查过了
	return config.GetConfig()
}

func (s *HuanProxyServer) GetRules() *rulescompile.RuleListCompileConfig {
	// 不用检查Ready，因为在NewServer的时候已经检查过了
	return config.GetRules()
}

func (s *HuanProxyServer) GetRulesList() []*rulescompile.RuleCompileConfig {
	// 不用检查Ready，因为在NewServer的时候已经检查过了
	return s.GetRules().Rules
}

func (s *HuanProxyServer) Run(httpschan chan error, httpchan chan error) (err error) {
	if s.https.cfg.Address != "" {
		err := s.https.loadHttps()
		if err != nil {
			return err
		}

		s.https.runHttps(httpschan)
	}

	if s.http.cfg.Address != "" {
		err := s.http.loadHttp()
		if err != nil {
			return err
		}

		s.http.runHttp(httpchan)
	}

	return nil
}

func (s *HTTPSServer) loadHttps() error {
	privateKey, certificate, issuerCertificate, err := certssl.GetCertificateAndPrivateKey(s.cfg.SSLCertDir, s.cfg.SSLEmail, s.cfg.AliyunDNSAccessKey, s.cfg.AliyunDNSAccessSecret, s.cfg.SSLDomain)
	if err != nil {
		return fmt.Errorf("init htttps cert ssl server error: %s", err.Error())
	} else if privateKey == nil || certificate == nil || issuerCertificate == nil {
		return fmt.Errorf("init https server error: get key and cert error, return nil, unknown reason")
	}

	s.key = privateKey
	s.cert = certificate
	s.cacert = issuerCertificate

	err = s.reloadHttps()
	if err != nil {
		return err
	}

	return nil
}

func (s *HTTPSServer) reloadHttps() error {
	if s.key == nil || s.cert == nil || s.cacert == nil {
		return fmt.Errorf("init https server error: get key and cert error, return nil, unknown reason")
	}

	if s.cert.Raw == nil || len(s.cert.Raw) == 0 || s.cacert.Raw == nil || len(s.cacert.Raw) == 0 {
		return fmt.Errorf("init https server error: get cert.raw error, return nil, unknown reason")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{s.cert.Raw, s.cacert.Raw}, // Raw包含 DER 编码的证书
			PrivateKey:  s.key,
			Leaf:        s.cert,
		}},
		MinVersion: tls.VersionTLS12,
	}

	s.server = &http.Server{
		Addr:      s.cfg.Address,
		Handler:   s.handler,
		TLSConfig: tlsConfig,
	}

	return nil
}

func (s *HTTPSServer) runHttps(_httpschan chan error) chan error {
	_watchstopchan := make(chan bool)

	s.watchCertificate(_watchstopchan)

	go func(httpschan chan error, watchstopchan chan bool) {
		defer func() {
			watchstopchan <- true
		}()
	ListenCycle:
		for {
			logger.Infof("start https server in %s", s.cfg.Address)
			err := s.server.ListenAndServeTLS("", "")
			if err != nil && errors.Is(err, http.ErrServerClosed) {
				if s.reloadMutex.TryLock() {
					s.reloadMutex.Unlock()
					_httpschan <- ServerStop
					return
				}
				s.reloadMutex.Lock()
				s.reloadMutex.Unlock() // 等待证书更换完毕
				continue ListenCycle
			} else if err != nil {
				_httpschan <- fmt.Errorf("https server error: %s", err.Error())
				return
			}
		}
	}(_httpschan, _watchstopchan)

	return _httpschan
}

func (s *HTTPSServer) watchCertificate(stopchan chan bool) {
	newchan := make(chan certssl.NewCert)

	go func() {
		err := certssl.WatchCertificate(s.cfg.SSLCertDir, s.cfg.SSLEmail, s.cfg.AliyunDNSAccessKey, s.cfg.AliyunDNSAccessSecret, s.cfg.SSLDomain, s.cert, stopchan, newchan)
		if err != nil {
			fmt.Printf("watch https cert server error: %s", err.Error())
		}
	}()

	go func() {
		select {
		case res := <-newchan:
			if res.Certificate == nil && res.PrivateKey == nil && res.Error == nil {
				close(newchan)
				return
			} else if res.Error != nil {
				fmt.Printf("https cert reload server error: %s", res.Error.Error())
			} else if res.PrivateKey != nil && res.Certificate != nil && res.IssuerCertificate != nil {
				func() {
					s.reloadMutex.Lock()
					defer s.reloadMutex.Unlock()

					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()

					err := s.server.Shutdown(ctx)
					if err != nil {
						fmt.Printf("https server reload shutdown error: %s", err.Error())
					}

					s.key = res.PrivateKey
					s.cert = res.Certificate
					s.cacert = res.IssuerCertificate

					err = s.reloadHttps()
					if err != nil {
						fmt.Printf("https server reload init error: %s", err.Error())
					}
				}()
			}
		}
	}()
}

func (s *HTTPServer) loadHttp() error {
	s.server = &http.Server{
		Addr:    s.cfg.Address,
		Handler: s.handler,
	}
	return nil
}

func (s *HTTPServer) runHttp(_httpschan chan error) chan error {
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

func (s *HuanProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.LoggerServerHTTP(w, r, s.NormalServeHTTP)
}
