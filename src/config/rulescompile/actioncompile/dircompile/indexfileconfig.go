package dircompile

import (
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/dir"
	"regexp"
)

type IndexFileCompileConfig struct {
	IsRegex bool
	File    string
	Regex   *regexp.Regexp
}

func NewIndexFileCompileConfig(i *dir.IndexFileConfig) (*IndexFileCompileConfig, error) {
	if i.Regex.IsEnable(true) {
		reg, err := regexp.Compile(i.File)
		if err != nil {
			return nil, err
		}
		return &IndexFileCompileConfig{
			IsRegex: true,
			File:    "",
			Regex:   reg,
		}, nil
	} else {
		return &IndexFileCompileConfig{
			IsRegex: false,
			File:    i.File,
			Regex:   nil,
		}, nil
	}
}

func (i *IndexFileCompileConfig) CheckName(name string) bool {
	if i.IsRegex {
		return i.Regex.MatchString(name)
	} else {
		return name == i.File
	}
}
