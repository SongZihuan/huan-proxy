package remotetrust

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/utils"
)

type RemoteTrustConfig struct {
	RemoteTrusted utils.StringBool `yaml:"remote-trusted"`
	TrustedIPs    []string         `yaml:"trusted-ips"`
}

func (p *RemoteTrustConfig) SetDefault() {
	p.RemoteTrusted.SetDefaultDisable()

	if p.RemoteTrusted.IsEnable() && len(p.TrustedIPs) == 0 {
		p.TrustedIPs = []string{"127.0.0.0/8", "::1"}
	}
}

func (p *RemoteTrustConfig) Check() configerr.ConfigError {
	if p.RemoteTrusted.IsEnable() {
		for _, ip := range p.TrustedIPs {
			if !utils.ValidIPv4(ip) && !utils.ValidIPv6(ip) && !utils.IsValidIPv4CIDR(ip) && !utils.IsValidIPv6CIDR(ip) {
				return configerr.NewConfigError(fmt.Sprintf("bad proxy trusts ip address: %s", ip))
			}
		}
	}
	return nil
}
