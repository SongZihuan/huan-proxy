package server

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func (s *HTTPServer) apiServer(ruleIndex int, rule *config.ProxyConfig, w http.ResponseWriter, r *http.Request) {
	proxy := s.cfg.ProxyServer.Get(ruleIndex)
	if proxy == nil {
		s.abortServerError(w)
		return
	}

	targetURL, err := url.Parse(rule.Address)
	if err != nil {
		s.abortServerError(w)
		return
	}

	s.processProxyHeader(r)

	r.URL.Scheme = targetURL.Scheme
	r.URL.Host = targetURL.Host

	path := r.URL.Path

	if strings.HasPrefix(path, rule.SubPrefixPath) {
		path = path[len(rule.SubPrefixPath):]
	}

	path = rule.AddPrefixPath + path
	r.URL.Path = path

	proxy.ServeHTTP(w, r) // 反向代理
}

func (s *HTTPServer) processProxyHeader(r *http.Request) {
	if r.RemoteAddr == "" {
		return
	}

	remoteIPStr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return
	}

	remoteIP := net.ParseIP(remoteIPStr)

	var ProxyList, ForwardedList []string
	var host, proto string

	if r.Header.Get("Forwarded") != "" {
		ProxyList, ForwardedList, host, proto = s.getProxyListForwarder(remoteIP, r)
	} else if r.Header.Get("X-Forwarded-For") != "" {
		ProxyList, ForwardedList, host, proto = s.getProxyListFromXForwardedFor(remoteIP, r)
	} else {
		host = r.Header.Get("X-Forwarded-Host")
		proto = r.Header.Get("X-Forwarded-Proto")

		if host == "" {
			host = r.URL.Host
		}

		if proto == "" {
			proto = r.URL.Scheme
		}

		ProxyList = append(make([]string, 0, 1), remoteIP.String())
		ForwardedList = append(make([]string, 0, 1),
			fmt.Sprintf("for=%s", remoteIP.String()),
			fmt.Sprintf("host=%s", host),
			fmt.Sprintf("proto=%s", proto))
	}

	r.Header.Set("Forwarded", strings.Join(ForwardedList, ","))
	r.Header.Set("X-Forwarded-For", strings.Join(ProxyList, ","))
	r.Header.Set("X-Forwarded-Host", host)
	r.Header.Set("X-Forwarded-Proto", proto)
}

func (s *HTTPServer) getProxyListForwarder(remoteIP net.IP, r *http.Request) ([]string, []string, string, string) {
	ForwardedList := strings.Split(r.Header.Get("Forwarded"), ",")
	ProxyList := make([]string, 0, len(ForwardedList)+1)
	NewForwardedList := make([]string, 0, len(ForwardedList)+1)

	host := r.URL.Host
	proto := r.URL.Scheme

	for _, keyStr := range ForwardedList {
		kv := strings.Split(strings.ReplaceAll(keyStr, " ", ""), "=")
		if len(kv) != 2 {
			continue
		}

		if kv[0] == "for" {
			forIP := net.ParseIP(strings.TrimSpace(kv[1]))
			if forIP != nil {
				NewForwardedList = append(NewForwardedList, keyStr)
				ProxyList = append(ProxyList, forIP.String())
			} else if kv[1] == "_hidden" || kv[1] == "_secret" || kv[1] == "unknown" {
				NewForwardedList = append(NewForwardedList, keyStr)
			}
		} else if kv[0] == "by" {
			byIP := net.ParseIP(strings.TrimSpace(kv[1]))
			if byIP != nil || kv[1] == "_hidden" || kv[1] == "_secret" || kv[1] == "unknown" {
				NewForwardedList = append(NewForwardedList, keyStr)
			}
		} else if kv[0] == "host" {
			host = kv[1]
		} else if kv[0] == "proto" {
			proto = kv[1]
		}
	}

	ProxyList = append(ProxyList, remoteIP.String())
	NewForwardedList = append(NewForwardedList, fmt.Sprintf("for=%s", remoteIP.String()))
	NewForwardedList = append(NewForwardedList, fmt.Sprintf("host=%s", host))
	NewForwardedList = append(NewForwardedList, fmt.Sprintf("proto=%s", proto))
	return ProxyList, NewForwardedList, host, proto
}

func (s *HTTPServer) getProxyListFromXForwardedFor(remoteIP net.IP, r *http.Request) ([]string, []string, string, string) {
	XFroWardedForList := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
	ProxyList := make([]string, 0, len(XFroWardedForList)+1)
	NewForwardedList := make([]string, 0, len(XFroWardedForList)+1)

	for _, forIPStr := range XFroWardedForList {
		forIP := net.ParseIP(strings.TrimSpace(forIPStr))
		if forIP != nil {
			ProxyList = append(ProxyList, forIP.String())
		}
	}

	host := r.Header.Get("X-Forwarded-Host")
	proto := r.Header.Get("X-Forwarded-Proto")

	if host == "" {
		host = r.URL.Host
	}

	if proto == "" {
		proto = r.URL.Scheme
	}

	ProxyList = append(ProxyList, remoteIP.String())

	for _, ip := range ProxyList {
		NewForwardedList = append(NewForwardedList, fmt.Sprintf("for=%s", ip))
	}
	NewForwardedList = append(NewForwardedList, fmt.Sprintf("host=%s", host))
	NewForwardedList = append(NewForwardedList, fmt.Sprintf("proto=%s", proto))

	return ProxyList, NewForwardedList, host, proto
}
