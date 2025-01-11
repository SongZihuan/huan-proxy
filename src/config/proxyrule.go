package config

import (
	"fmt"
)

type ProxyRuleConfig struct {
	Rules []*ProxyConfig `yaml:"rules"`
}

func (r *ProxyRuleConfig) setDefault() {
	for _, rule := range r.Rules {
		rule.setDefault()
	}
}

func (r *ProxyRuleConfig) check(ps *ProxyServerConfig, ixfile *IndexFileCompileList, igfile *IgnoreFileCompileList) ConfigError {
	if len(r.Rules) == 0 {
		return NewConfigError("proxy rule is empty")
	}

	for ruleIndex, rule := range r.Rules {
		err := rule.check()
		if err != nil && err.IsError() {
			return err
		}

		if rule.Type == ProxyTypeDir {
			for fileIndex, file := range rule.IndexFile {
				err := ixfile.Add(ruleIndex, fileIndex, file)
				if err != nil {
					return NewConfigError(fmt.Sprintf("index file %s error", err.Error()))
				}
			}

			for fileIndex, file := range rule.IgnoreFile {
				err := igfile.Add(ruleIndex, fileIndex, file)
				if err != nil {
					return NewConfigError(fmt.Sprintf("ignore file %s error", err.Error()))
				}
			}
		} else if rule.Type == ProxyTypeAPI {
			err := ps.Add(ruleIndex, rule)
			if err != nil {
				return NewConfigError(fmt.Sprintf("proxy server can not create: %s", err.Error()))
			}
		}
	}

	return nil
}
