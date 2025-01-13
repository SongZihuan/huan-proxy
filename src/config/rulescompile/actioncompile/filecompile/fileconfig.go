package filecompile

import (
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/file"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/corscompile"
)

type RuleFileCompileConfig struct {
	Path string
	Cors *corscompile.CorsCompileConfig
}

func NewRuleFileCompileConfig(f *file.RuleFileConfig) (*RuleFileCompileConfig, error) {
	cors, err := corscompile.NewCorsCompileConfig(&f.Cors)
	if err != nil {
		return nil, err
	}

	return &RuleFileCompileConfig{
		Path: f.Path,
		Cors: cors,
	}, nil
}
