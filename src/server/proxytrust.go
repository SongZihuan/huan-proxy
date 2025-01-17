package server

import (
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"net"
	"net/http"
)

func (s *HuanProxyServer) checkProxyTrust(rule *rulescompile.RuleCompileConfig, w http.ResponseWriter, r *http.Request) bool {
	if !rule.UseTrustedIPs {
		return true
	}

	if r.RemoteAddr == "" {
		s.abortForbidden(w)
		return false
	}

	remoteIPStr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		s.abortForbidden(w)
		return false
	}

	remoteIP := net.ParseIP(remoteIPStr)
	if remoteIP == nil {
		s.abortForbidden(w)
		return false
	}

	for _, t := range rule.TrustedIPs {
		if utils.ValidIPv4(t) || utils.ValidIPv6(t) {
			trustIP := net.ParseIP(t)
			if trustIP == nil {
				continue
			} else if trustIP.Equal(remoteIP) {
				return true
			}
		} else if utils.IsValidIPv4CIDR(t) || utils.IsValidIPv6CIDR(t) {
			_, trustCIDR, err := net.ParseCIDR(t)
			if err != nil || trustCIDR == nil {
				continue
			} else if trustCIDR.Contains(remoteIP) {
				return true
			}
		}
	}

	s.abortForbidden(w)
	return false
}
