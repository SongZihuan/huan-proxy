package config

import (
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"os"
)

const (
	EnvAliyunKey    = "HP_ALIYUN_ACCESS_KEY"
	EnvAliyunSecret = "HP_ALIYUN_ACCESS_SECRET"
)

type HttpsConfig struct {
	Address               string `yaml:"address"`
	SSLEmail              string `json:"sslemail"`
	SSLDomain             string `yaml:"ssldomain"`
	SSLCertDir            string `yaml:"sslcertdir"`
	AliyunDNSAccessKey    string `yaml:"aliyundnsaccesskey"`
	AliyunDNSAccessSecret string `yaml:"aliyundnsaccesssecret"`
	StopWaitSecond        int    `yaml:"stopwaitsecond"`
}

func (h *HttpsConfig) SetDefault() {
	if h.Address != "" {
		if h.SSLCertDir == "" {
			h.SSLCertDir = "./ssl-certs"
		}

		if h.AliyunDNSAccessKey == "" {
			h.AliyunDNSAccessKey = os.Getenv(EnvAliyunKey)
		}

		if h.AliyunDNSAccessSecret == "" {
			h.AliyunDNSAccessKey = os.Getenv(EnvAliyunSecret)
		}
	}

	if h.StopWaitSecond <= 0 {
		h.StopWaitSecond = 10
	}
}

func (h *HttpsConfig) Check() configerr.ConfigError {
	if h.Address != "" {
		if h.SSLEmail == "" || !utils.IsValidEmail(h.SSLEmail) {
			return configerr.NewConfigError("http ssl must has a valid email")
		}

		if h.SSLDomain == "" {
			return configerr.NewConfigError("http ssl must has a domain")
		}

		if h.AliyunDNSAccessKey == "" || h.AliyunDNSAccessSecret == "" {
			return configerr.NewConfigError("http ssl must has a aliyun access key or secret")
		}
	}

	return nil
}
