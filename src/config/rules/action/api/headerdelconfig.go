package api

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/utils"
)

type HeaderDelConfig struct {
	Header string `yaml:"header"`
}

func (h *HeaderDelConfig) SetDefault() {

}

func (h *HeaderDelConfig) Check() configerr.ConfigError {
	if h.Header == "" {
		return configerr.NewConfigError("header name is empty")
	}

	if h.Header == ViaHeader || h.Header == XHuanProxyHeaer {
		return configerr.NewConfigError(fmt.Sprintf("header %s use by http system", h.Header))
	}

	if !utils.IsValidHTTPHeaderKey(h.Header) {
		return configerr.NewConfigError(fmt.Sprintf("header %s is not valid", h.Header))
	}

	if isNotGoodHeader(h.Header) {
		_ = configerr.NewConfigWarning(fmt.Sprintf("header %s use by http system", h.Header))
	}

	return nil
}