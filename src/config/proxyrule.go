package config

import "fmt"

type ProxyRuleConfig struct {
	Rules []*ProxyConfig `yaml:"rules"`
}

func (r *ProxyRuleConfig) setDefault() {
	for _, rule := range r.Rules {
		rule.setDefault()
		fmt.Printf("TAG DD [%s]\n", rule.BasePath)
	}
	fmt.Printf("TAG DDAA [%s]\n", r.Rules[0].BasePath)
}

func (r *ProxyRuleConfig) check(ps *ProxyServerConfig) ConfigError {
	if len(r.Rules) == 0 {
		return NewConfigError("proxy rule is empty")
	}

	for index, rule := range r.Rules {
		err := rule.check()
		if err != nil && err.IsError() {
			return err
		}

		if rule.Type == ProxyTypeAPI {
			err := ps.Add(index, rule)
			if err != nil {
				return NewConfigError(fmt.Sprintf("proxy server can not create: %s", err.Error()))
			}
		}
	}

	return nil
}
