package config

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"net/url"
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
	Dir string `yaml:"dir"`
}

type ProxyAPIConfig struct {
	Address    string           `yaml:"address"`
	PrefixPath string           `yaml:"prefixpath"`
	EnableSSL  utils.StringBool `yaml:"enablessl"`
}

func (p *ProxyConfig) setDefault() {
	p.BasePath = utils.ProcessPath(p.BasePath)
	fmt.Printf("TAG NB [%s]\n", p.BasePath)

	if p.Type == ProxyTypeAPI {
		p.PrefixPath = utils.ProcessPath(p.PrefixPath)
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
	} else if p.Type == ProxyTypeAPI {
		_, err := url.Parse(p.Address)
		if err != nil {
			return NewConfigError(fmt.Sprintf("Failed to parse target URL: %v", err))
		}
	} else {
		return NewConfigError("proxy type must be file or dir or api")
	}
	fmt.Printf("TAG CC [%s]\n", p.BasePath)
	return nil
}
