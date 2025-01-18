package server

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/rewritecompile"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"net"
	"strings"
)

func (s *HuanProxyServer) apiServer(ctx *Context) {
	proxy := ctx.Rule.Api.Server
	if proxy == nil {
		s.abortServerError(ctx)
		return
	}

	targetURL := ctx.Rule.Api.TargetURL
	ctx.ProxyRequest.URL.Scheme = targetURL.Scheme
	ctx.ProxyRequest.URL.Host = targetURL.Host

	s.processProxyHeader(ctx)

	ctx.ProxyRequest.URL.Path = s.apiRewrite(utils.ProcessURLPath(ctx.ProxyRequest.URL.Path), ctx.Rule.Api.AddPath, ctx.Rule.Api.SubPath, ctx.Rule.Api.Rewrite)

	for _, h := range ctx.Rule.Api.HeaderSet {
		ctx.ProxyRequest.Header.Set(h.Header, h.Value)
	}

	for _, h := range ctx.Rule.Api.HeaderAdd {
		ctx.ProxyRequest.Header.Add(h.Header, h.Value)
	}

	for _, h := range ctx.Rule.Api.HeaderDel {
		ctx.ProxyRequest.Header.Del(h.Header)
	}

	query := ctx.ProxyRequest.URL.Query()

	for _, q := range ctx.Rule.Api.QuerySet {
		query.Set(q.Query, q.Value)
	}

	for _, q := range ctx.Rule.Api.QueryAdd {
		query.Add(q.Query, q.Value)
	}

	for _, q := range ctx.Rule.Api.QueryDel {
		query.Del(q.Query)
	}

	ctx.ProxyRequest.URL.RawQuery = query.Encode()

	s.writeViaHeader(ctx)

	req, err := ctx.ProxyWriteToHttpRRequest()
	if err != nil {
		s.abortServerError(ctx)
		return
	}

	proxy.ServeHTTP(ctx.Writer, req) // 反向代理
}

func (s *HuanProxyServer) apiRewrite(srcpath string, prefix string, suffix string, rewrite *rewritecompile.RewriteCompileConfig) string {
	srcpath = utils.ProcessURLPath(srcpath)
	prefix = utils.ProcessURLPath(prefix)
	suffix = utils.ProcessURLPath(suffix)

	if strings.HasPrefix(srcpath, suffix) {
		srcpath = srcpath[len(suffix):]
	}

	srcpath = prefix + srcpath

	if rewrite.Use && rewrite.Regex != nil {
		rewrite.Regex.ReplaceAllString(srcpath, rewrite.Target)
	}

	return srcpath
}

func (s *HuanProxyServer) processProxyHeader(ctx *Context) {
	if ctx.Request.RemoteAddr() == "" {
		return
	}

	remoteIPStr, _, err := net.SplitHostPort(ctx.Request.RemoteAddr())
	if err != nil {
		return
	}

	remoteIP := net.ParseIP(remoteIPStr)

	var ProxyList, ForwardedList []string
	var host, proto string

	if ctx.Request.Header().Get("Forwarded") != "" {
		ProxyList, ForwardedList, host, proto = s.getProxyListForwarder(remoteIP, ctx.Request)
	} else if ctx.Request.Header().Get("X-Forwarded-For") != "" {
		ProxyList, ForwardedList, host, proto = s.getProxyListFromXForwardedFor(remoteIP, ctx.Request)
	} else {
		host = ctx.Request.Header().Get("X-Forwarded-Host")
		proto = ctx.Request.Header().Get("X-Forwarded-Proto")

		if host == "" {
			host = ctx.Request.Host()
		}

		host, _ = utils.SplitHostPort(host) // 去除host中的端口号

		if proto == "http" || proto == "https" {
			if ctx.Request.IsTLS() {
				proto = "https"
			} else {
				proto = "http"
			}
		}

		ProxyList = append(make([]string, 0, 1), remoteIP.String())
		ForwardedList = append(make([]string, 0, 1),
			fmt.Sprintf("for=%s", remoteIP.String()),
			fmt.Sprintf("host=%s", host),
			fmt.Sprintf("proto=%s", proto))
	}

	ctx.ProxyRequest.Header.Set("Forwarded", strings.Join(ForwardedList, ", "))
	ctx.ProxyRequest.Header.Set("X-Forwarded-For", strings.Join(ProxyList, ", "))
	ctx.ProxyRequest.Header.Set("X-Forwarded-Host", host)
	ctx.ProxyRequest.Header.Set("X-Forwarded-Proto", proto)
}

func (s *HuanProxyServer) getProxyListForwarder(remoteIP net.IP, r *ReadOnlyRequest) ([]string, []string, string, string) {
	ForwardedList := strings.Split(r.Header().Get("Forwarded"), ",")
	ProxyList := make([]string, 0, len(ForwardedList)+1)
	NewForwardedList := make([]string, 0, len(ForwardedList)+1)

	host, _ := utils.SplitHostPort(r.Host()) // 去除host中的端口号
	proto := "http"
	if r.IsTLS() {
		proto = "https"
	}

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

func (s *HuanProxyServer) getProxyListFromXForwardedFor(remoteIP net.IP, r *ReadOnlyRequest) ([]string, []string, string, string) {
	XFroWardedForList := strings.Split(r.Header().Get("X-Forwarded-For"), ",")
	ProxyList := make([]string, 0, len(XFroWardedForList)+1)
	NewForwardedList := make([]string, 0, len(XFroWardedForList)+1)

	for _, forIPStr := range XFroWardedForList {
		forIP := net.ParseIP(strings.TrimSpace(forIPStr))
		if forIP != nil {
			ProxyList = append(ProxyList, forIP.String())
		}
	}

	host := r.Header().Get("X-Forwarded-Host")
	proto := r.Header().Get("X-Forwarded-Proto")

	if host == "" {
		host = r.Host()
	}

	host, _ = utils.SplitHostPort(host) // 去除host中的端口号

	if proto == "http" || proto == "https" {
		if r.IsTLS() {
			proto = "https"
		} else {
			proto = "http"
		}
	}

	ProxyList = append(ProxyList, remoteIP.String())

	for _, ip := range ProxyList {
		NewForwardedList = append(NewForwardedList, fmt.Sprintf("for=%s", ip))
	}
	NewForwardedList = append(NewForwardedList, fmt.Sprintf("host=%s", host))
	NewForwardedList = append(NewForwardedList, fmt.Sprintf("proto=%s", proto))

	return ProxyList, NewForwardedList, host, proto
}
