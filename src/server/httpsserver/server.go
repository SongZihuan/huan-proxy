package httpsserver

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/certssl"
	"github.com/SongZihuan/huan-proxy/src/config"
	"github.com/SongZihuan/huan-proxy/src/logger"
	"github.com/pires/go-proxyproto"
	"net"
	"net/http"
	"sync"
	"time"
)

var ServerStop = fmt.Errorf("https server stop")

type HTTPSServer struct {
	cfg         *config.HttpsConfig
	reloadMutex sync.Mutex
	key         crypto.PrivateKey
	cert        *x509.Certificate
	cacert      *x509.Certificate
	server      *http.Server
	handler     http.Handler
}

func NewHTTPSServer(handler http.Handler) *HTTPSServer {
	httpsCfg := config.GetConfig().Https

	if httpsCfg.Address == "" {
		return nil
	}

	return &HTTPSServer{
		cfg:     &httpsCfg,
		server:  nil,
		handler: handler,
	}
}

func (s *HTTPSServer) LoadHttps() error {
	privateKey, certificate, issuerCertificate, err := certssl.GetCertificateAndPrivateKey(s.cfg.SSLCertDir, s.cfg.SSLEmail, s.cfg.AliyunDNSAccessKey, s.cfg.AliyunDNSAccessSecret, s.cfg.SSLDomain)
	if err != nil {
		return fmt.Errorf("init https cert ssl server error: %s", err.Error())
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

func (s *HTTPSServer) StopHttps() error {
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

func (s *HTTPSServer) reloadHttps() error {
	if s.key == nil || s.cert == nil || s.cacert == nil {
		return fmt.Errorf("init https server error: get key and cert error, return nil, unknown reason")
	}

	if s.cert.Raw == nil || len(s.cert.Raw) == 0 || s.cacert.Raw == nil || len(s.cacert.Raw) == 0 {
		return fmt.Errorf("init https server error: get cert.raw error, return nil, unknown reason")
	}

	s.server = &http.Server{
		Addr:    s.cfg.Address,
		Handler: s.handler,
	}
	return nil
}

func (s *HTTPSServer) getListener() (net.Listener, error) {
	if s.key == nil || s.cert == nil || s.cacert == nil {
		return nil, fmt.Errorf("init https server error: get key and cert error, return nil, unknown reason")
	}

	if s.cert.Raw == nil || len(s.cert.Raw) == 0 || s.cacert.Raw == nil || len(s.cacert.Raw) == 0 {
		return nil, fmt.Errorf("init https server error: get cert.raw error, return nil, unknown reason")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{s.cert.Raw, s.cacert.Raw}, // Raw包含 DER 编码的证书
			PrivateKey:  s.key,
			Leaf:        s.cert,
		}},
		MinVersion: tls.VersionTLS12,
	}

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

	tlsListener := tls.NewListener(proxyListener, tlsConfig)

	return tlsListener, nil
}

func (s *HTTPSServer) RunHttps(httpsErrorChan chan error) {
	watchStopChan := make(chan bool)

	s.watchCertificate(watchStopChan)

	go func() {
		defer func() {
			close(watchStopChan)
		}()

		defer func() {
			s.server = nil
		}()

		for {
			res := func() bool {
				listener, err := s.getListener()
				if err != nil {
					httpsErrorChan <- fmt.Errorf("listen fail")
					return true
				}
				defer func() {
					_ = listener.Close()
				}()

				logger.Infof("start https server in %s", s.cfg.Address)
				err = s.server.Serve(listener)
				if err != nil && errors.Is(err, http.ErrServerClosed) {
					if s.reloadMutex.TryLock() {
						s.reloadMutex.Unlock()
						httpsErrorChan <- ServerStop
						return true
					}

					s.reloadMutex.Lock()
					s.reloadMutex.Unlock() // 等待证书更换完毕
					return false
				} else if err != nil {
					httpsErrorChan <- fmt.Errorf("https server error: %s", err.Error())
					return true
				}

				return false
			}()
			if res {
				return
			}
		}
	}()
}

func (s *HTTPSServer) watchCertificate(stopChan chan bool) {
	newCertChan := make(chan certssl.NewCert)

	go func() {
		err := certssl.WatchCertificate(s.cfg.SSLCertDir, s.cfg.SSLEmail, s.cfg.AliyunDNSAccessKey, s.cfg.AliyunDNSAccessSecret, s.cfg.SSLDomain, s.cert, stopChan, newCertChan)
		if err != nil {
			logger.Errorf("watch https cert server error: %s", err.Error())
		}
	}()

	go func() {
		defer close(newCertChan)

		for {
			select {
			case <-stopChan:
				return
			case res := <-newCertChan:
				if res.Error != nil {
					logger.Errorf("https cert reload server error: %s", res.Error.Error())
				} else if res.PrivateKey != nil && res.Certificate != nil && res.IssuerCertificate != nil {
					func() {
						s.reloadMutex.Lock()
						defer s.reloadMutex.Unlock()

						ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
						defer cancel()

						err := s.server.Shutdown(ctx)
						if err != nil {
							logger.Errorf("https server reload shutdown error: %s", err.Error())
						}

						s.key = res.PrivateKey
						s.cert = res.Certificate
						s.cacert = res.IssuerCertificate

						err = s.reloadHttps()
						if err != nil {
							logger.Errorf("https server reload init error: %s", err.Error())
						}
					}()
				}
			}
		}
	}()
}
