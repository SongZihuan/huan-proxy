package rewrite

import "github.com/SongZihuan/huan-proxy/src/config/configerr"

type RewriteConfig struct {
	Regex  string `yaml:"regex"`
	Target string `yaml:"target"`
}

func (r *RewriteConfig) SetDefault() {

}

func (r *RewriteConfig) Check() configerr.ConfigError {
	if len(r.Target) != 0 && len(r.Regex) == 0 {
		return configerr.NewConfigError("rewrite reg is empty")
	}

	return nil
}
