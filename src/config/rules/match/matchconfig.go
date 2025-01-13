package match

import (
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/utils"
)

const (
	PrefixMatch    = "prefix"    // 前缀匹配
	RegexMatch     = "regex"     // 正则匹配
	PrecisionMatch = "precision" // 精准匹配
)

type MatchConfig struct {
	MatchType string `yaml:"matchtype"`
	Path      string `yaml:"path"`
}

func (m *MatchConfig) SetDefault() {
	m.Path = utils.ProcessPath(m.Path)
}

func (m *MatchConfig) Check() configerr.ConfigError {
	if m.MatchType != PrefixMatch && m.MatchType != RegexMatch && m.MatchType != PrecisionMatch {
		return configerr.NewConfigError("proxy mutch type must be prefix or regex or precision")
	}
	return nil
}
