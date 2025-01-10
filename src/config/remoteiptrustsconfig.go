package config

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/utils"
)

type RemoteTrustConfig struct {
	RemoteTrust utils.StringBool `json:"remotetrust"`
	TrustedIPs  []string         `json:"trustedips"`
}

func (p *RemoteTrustConfig) setDefault(global *GlobalConfig) {
	if global.IsDebug() || global.IsTest() {
		p.RemoteTrust.SetDefaultEanble()
	} else {
		p.RemoteTrust.SetDefaultDisable()
	}

	if p.RemoteTrust.IsEnable() && len(p.TrustedIPs) == 0 {
		p.TrustedIPs = []string{"127.0.0.0/8", "::1"}
	}
}

func (p *RemoteTrustConfig) check() ConfigError {
	if p.RemoteTrust.IsEnable() {
		if len(p.TrustedIPs) == 0 {
			_ = NewConfigWarning("proxy trusts ips will be ignore because proxy is disabled")
		} else {
			for _, ip := range p.TrustedIPs {
				if !utils.ValidIPv4(ip) && !utils.ValidIPv6(ip) && !utils.IsValidIPv4CIDR(ip) && !utils.IsValidIPv6CIDR(ip) {
					return NewConfigError(fmt.Sprintf("bad proxy trusts ip address: %s", ip))
				}
			}
		}
	} else {
		_ = NewConfigWarning("You trusted all proxies, this is NOT safe. We recommend you to set a value.")
	}

	return nil
}

func (p *RemoteTrustConfig) Enable() bool {
	return p.RemoteTrust.IsEnable()
}
