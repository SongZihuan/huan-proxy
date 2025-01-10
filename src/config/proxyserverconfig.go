package config

import (
	"fmt"
	"net/http/httputil"
	"net/url"
)

const defaultServerPort = 10

type ProxyServerConfig struct {
	Server map[int]*httputil.ReverseProxy
}

func (p *ProxyServerConfig) init() error {
	p.Server = make(map[int]*httputil.ReverseProxy, defaultServerPort)
	return nil
}

func (p *ProxyServerConfig) Add(index int, rule *ProxyConfig) error {
	if rule.Type != ProxyTypeAPI {
		return nil
	}

	if _, ok := p.Server[index]; ok {
		return fmt.Errorf("proxy server %d already exists", index)
	}

	targetURL, err := url.Parse(rule.Address)
	if err != nil {
		return err
	}

	p.Server[index] = httputil.NewSingleHostReverseProxy(targetURL)

	return nil
}

func (p *ProxyServerConfig) Get(index int) *httputil.ReverseProxy {
	if proxy, ok := p.Server[index]; ok {
		return proxy
	}
	return nil
}
