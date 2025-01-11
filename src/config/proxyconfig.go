package config

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"net/url"
	"strings"
)

const (
	ProxyTypeFile = "file"
	ProxyTypeDir  = "dir"
	ProxyTypeAPI  = "api"
)

type ProxyConfig struct {
	Type     string `yaml:"type"`
	BasePath string `yaml:"basepath"`

	ProxyFileConfig `yaml:",inline"`
	ProxyDirConfig  `yaml:",inline"`
	ProxyAPIConfig  `yaml:",inline"`
}

type ProxyFileConfig struct {
	File string `yaml:"file"`
}

type ProxyDirConfig struct {
	Dir        string        `yaml:"dir"`
	IndexFile  []*IndexFile  `yaml:"indexfile"`
	IgnoreFile []*IgnoreFile `yaml:"ignorefile"`
}

type ProxyAPIConfig struct {
	Address       string           `yaml:"address"`
	AddPrefixPath string           `yaml:"addprefixpath"`
	SubPrefixPath string           `yaml:"subprefixpath"`
	EnableSSL     utils.StringBool `yaml:"enablessl"`
}

func (p *ProxyConfig) setDefault() {
	p.BasePath = utils.ProcessPath(p.BasePath)

	if p.Type == ProxyTypeAPI {
		if len(p.IndexFile) == 0 {
			p.IndexFile = []*IndexFile{
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
					File:  "index.txtzzz",
				},
				{
					Regex: "enable",
					File:  `^index\.\S+$`,
				},
			}
		}
	} else if p.Type == ProxyTypeAPI {
		p.AddPrefixPath = utils.ProcessPath(p.AddPrefixPath)
		p.SubPrefixPath = utils.ProcessPath(p.SubPrefixPath)
	}
}

func (p *ProxyConfig) check() ConfigError {
	if p.Type == ProxyTypeFile {
		if !utils.IsFile(p.File) {
			return NewConfigError(fmt.Sprintf("file path %s not exist", p.File))
		}
	} else if p.Type == ProxyTypeDir {
		if !utils.IsDir(p.Dir) {
			return NewConfigError(fmt.Sprintf("dir path %s not exist", p.Dir))
		}

		if len(p.IndexFile) == 0 {
			return NewConfigError("index file is empty")
		}
	} else if p.Type == ProxyTypeAPI {
		_, err := url.Parse(p.Address)
		if err != nil {
			return NewConfigError(fmt.Sprintf("Failed to parse target URL: %v", err))
		}

		if strings.HasPrefix(p.BasePath, p.SubPrefixPath) {
			return NewConfigError("sub prefix path error")
		}

	} else {
		return NewConfigError("proxy type must be file or dir or api")
	}
	return nil
}
