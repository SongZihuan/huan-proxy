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
	Address       string `yaml:"address"`
	AddPrefixPath string `yaml:"addprefixpath"`
	SubPrefixPath string `yaml:"subprefixpath"`
	RewriteReg    string `yaml:"rewritereg"`
	RewriteTarget string `yaml:"rewritetarget"`
}

func (p *ProxyConfig) setDefault() {
	p.BasePath = utils.ProcessPath(p.BasePath)

	if p.Type == ProxyTypeDir {
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

		for _, i := range p.IndexFile {
			i.setDefault()
		}

		for _, i := range p.IgnoreFile {
			i.setDefault()
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

		for _, i := range p.IndexFile {
			err := i.check()
			if err != nil && err.IsError() {
				return err
			}
		}

		for _, i := range p.IgnoreFile {
			err := i.check()
			if err != nil && err.IsError() {
				return err
			}
		}
	} else if p.Type == ProxyTypeAPI {
		_, err := url.Parse(p.Address)
		if err != nil {
			return NewConfigError(fmt.Sprintf("Failed to parse target URL: %v", err))
		}

		if p.BasePath != p.SubPrefixPath && !strings.HasPrefix(p.BasePath, p.SubPrefixPath+"/") {
			return NewConfigError("sub prefix path error")
		}

		if len(p.RewriteTarget) != 0 && len(p.RewriteReg) == 0 {
			return NewConfigError("rewrite reg is empty")
		}
	} else {
		return NewConfigError("proxy type must be file or dir or api")
	}
	return nil
}
