package dir

import (
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/cors"
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/rewrite"
)

type RuleDirConfig struct {
	BasePath   string                `yaml:"base-path"`
	IndexFile  []*IndexFileConfig    `yaml:"index-file"`
	IgnoreFile []*IgnoreFileConfig   `yaml:"ignore-file"`
	AddPath    string                `yaml:"add-path"`
	SubPath    string                `yaml:"sub-path"`
	Rewrite    rewrite.RewriteConfig `yaml:"rewrite"`
	Cors       cors.CorsConfig       `yaml:"cors"`
}

func (r *RuleDirConfig) SetDefault() {
	if len(r.IndexFile) == 0 {
		r.IndexFile = []*IndexFileConfig{
			{
				Regex: "disable",
				File:  "index.html",
			},
			{
				Regex: "disable",
				File:  "index.xml",
			},
			{
				Regex: "disable",
				File:  "index",
			},
			{
				Regex: "enable",
				File:  `^index\.\S+$`,
			},
		}
	}

	for _, i := range r.IndexFile {
		i.SetDefault()
	}

	for _, i := range r.IgnoreFile {
		i.SetDefault()
	}

	r.Rewrite.SetDefault()
	r.Cors.SetDefault()
}

func (r *RuleDirConfig) Check() configerr.ConfigError {
	// 不用检查目录是否存在，因为可能被rewrite
	for _, i := range r.IndexFile {
		err := i.Check()
		if err != nil && err.IsError() {
			return err
		}
	}

	for _, i := range r.IgnoreFile {
		err := i.Check()
		if err != nil && err.IsError() {
			return err
		}
	}

	err := r.Rewrite.Check()
	if err != nil && err.IsError() {
		return err
	}

	err = r.Cors.Check()
	if err != nil && err.IsError() {
		return err
	}

	return nil
}
