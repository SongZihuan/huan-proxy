package cors

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/utils"
)

const CorsMaxAgeSec = 86400
const CorsDefaultMaxAgeSec = CorsMaxAgeSec

type CorsConfig struct {
	AllowCors      utils.StringBool `yaml:"allowcors"`
	AllowOrigin    []string         `yaml:"alloworigin"`
	AllowOriginReg []string         `yaml:"alloworiginres"`
	MaxAgeSec      int              `yaml:"maxagesec"`
}

func (c *CorsConfig) SetDefault() {
	c.AllowCors.SetDefaultDisable()
	if c.AllowCors.IsEnable() && c.MaxAgeSec == 0 {
		c.MaxAgeSec = CorsDefaultMaxAgeSec
	}
}

func (c *CorsConfig) Check() configerr.ConfigError {
	if c.AllowCors.IsEnable() {
		if c.MaxAgeSec <= 0 || c.MaxAgeSec > CorsMaxAgeSec {
			return configerr.NewConfigError(fmt.Sprintf("cors maxagesec %d is invalid", c.MaxAgeSec))
		}
	}
	return nil
}
