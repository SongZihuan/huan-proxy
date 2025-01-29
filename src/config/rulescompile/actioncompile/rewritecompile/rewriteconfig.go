package rewritecompile

import (
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/rewrite"
	"regexp"
)

type RewriteCompileConfig struct {
	Use    bool
	Regex  *regexp.Regexp
	Target string
}

func NewRewriteCompileConfig(r *rewrite.RewriteConfig) (*RewriteCompileConfig, error) {
	if r.Regex == "" {
		return &RewriteCompileConfig{
			Use:    false,
			Regex:  nil,
			Target: "",
		}, nil
	}

	reg, err := regexp.Compile(r.Regex)
	if err != nil {
		return nil, err
	}

	return &RewriteCompileConfig{
		Use:    true,
		Regex:  reg,
		Target: r.Target,
	}, nil
}
