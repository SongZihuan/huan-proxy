package rulescompile

import "github.com/SongZihuan/huan-proxy/src/config/rules"

type RuleListCompileConfig struct {
	Rules []*RuleCompileConfig `yaml:"rules"`
}

func NewRuleListConfig(rs *rules.RuleListConfig) (*RuleListCompileConfig, error) {
	res := make([]*RuleCompileConfig, 0, len(rs.Rules))
	for _, v := range rs.Rules {
		r, err := NewRuleCompileConfig(v)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}

	return &RuleListCompileConfig{
		Rules: res,
	}, nil
}
