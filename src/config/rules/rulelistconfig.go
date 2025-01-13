package rules

import (
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
)

type RuleListConfig struct {
	Rules []RuleConfig `yaml:"rules"`
}

func (r *RuleListConfig) SetDefault() {
	for _, rule := range r.Rules {
		rule.SetDefault()
	}
}

func (r *RuleListConfig) Check() configerr.ConfigError {
	for _, rule := range r.Rules {
		err := rule.Check()
		if err != nil && err.IsError() {
			return err
		}
	}

	return nil
}
