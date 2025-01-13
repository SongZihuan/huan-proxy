package filecompile

import (
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/file"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/corscompile"
)

type RuleFileCompileConfig struct {
	File string
	Cors *corscompile.CorsCompileConfig
}

func NewRuleFileCompileConfig(f *file.RuleFileConfig) (*RuleFileCompileConfig, error) {
	cors, err := corscompile.NewCorsCompileConfig(&f.FileCors)
	if err != nil {
		return nil, err
	}

	return &RuleFileCompileConfig{
		File: f.File,
		Cors: cors,
	}, nil
}
