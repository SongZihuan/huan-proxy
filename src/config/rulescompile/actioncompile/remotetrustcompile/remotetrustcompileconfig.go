package remotetrustcompile

import "github.com/SongZihuan/huan-proxy/src/config/rules/action/remotetrust"

type RemoteTrustCompileConfig struct {
	UseTrustedIPs bool
	TrustedIPs    []string
}

func NewRemoteTrustCompileConfig(r *remotetrust.RemoteTrustConfig) (*RemoteTrustCompileConfig, error) {
	if r.RemoteTrust.IsDisable(false) {
		return &RemoteTrustCompileConfig{
			UseTrustedIPs: false,
			TrustedIPs:    make([]string, 0),
		}, nil
	} else {
		trustedIPs := make([]string, len(r.TrustedIPs))
		copy(trustedIPs, r.TrustedIPs)
		return &RemoteTrustCompileConfig{
			UseTrustedIPs: true,
			TrustedIPs:    trustedIPs,
		}, nil
	}
}
