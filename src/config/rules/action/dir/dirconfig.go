package dir

import (
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/cors"
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/rewrite"
	"github.com/SongZihuan/huan-proxy/src/utils"
)

type RuleDirConfig struct {
	Dir              string                `yaml:"dir"`
	IndexFile        []IndexFileConfig     `yaml:"indexfile"`
	IgnoreFile       []IgnoreFileConfig    `yaml:"ignorefile"`
	DirAddPrefixPath string                `yaml:"addprefixpath"` // Dir前缀避免充满（yaaml忽略）
	DirSubPrefixPath string                `yaml:"subprefixpath"` // Dir前缀避免充满（yaaml忽略）
	DirRewrite       rewrite.RewriteConfig `yaml:"rewrite"`       // Dir前缀避免充满（yaaml忽略）
	DirCors          cors.CorsConfig       `yaml:"cors"`          // Dir前缀避免充满（yaaml忽略）
}

func (r *RuleDirConfig) SetDefault() {
	r.Dir = utils.ProcessPath(r.Dir)
	r.DirAddPrefixPath = utils.ProcessPath(r.DirAddPrefixPath)
	r.DirSubPrefixPath = utils.ProcessPath(r.DirSubPrefixPath)

	if len(r.IndexFile) == 0 {
		r.IndexFile = []IndexFileConfig{
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

	r.DirRewrite.SetDefault()
	r.DirCors.SetDefault()
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

	err := r.DirRewrite.Check()
	if err != nil && err.IsError() {
		return err
	}

	err = r.DirCors.Check()
	if err != nil && err.IsError() {
		return err
	}

	return nil
}
