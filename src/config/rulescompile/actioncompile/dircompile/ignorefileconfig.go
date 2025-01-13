package dircompile

import (
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/dir"
	"regexp"
)

type IgnoreFileCompileConfig struct {
	IsRegex bool
	File    string
	Regex   *regexp.Regexp
}

func NewIgnoreFileCompileConfig(i *dir.IgnoreFileConfig) (*IgnoreFileCompileConfig, error) {
	if i.Regex.IsEnable(true) {
		reg, err := regexp.Compile(i.File)
		if err != nil {
			return nil, err
		}
		return &IgnoreFileCompileConfig{
			IsRegex: true,
			File:    "",
			Regex:   reg,
		}, nil
	} else {
		return &IgnoreFileCompileConfig{
			IsRegex: false,
			File:    i.File,
			Regex:   nil,
		}, nil
	}
}
