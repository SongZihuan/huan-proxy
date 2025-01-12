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

const XHuanProxyHeaer = "X-Huan-Proxy"
const ViaHeader = "Via"

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

type HeaderConfig struct {
	Header string `yaml:"header"`
	Value  string `yaml:"value"`
}

type QueryConfig struct {
	Query string `yaml:"query"`
	Value string `yaml:"value"`
}

var WarningHeader = []string{
	"Host",
	"Referer",
	"User-Agent",
	"Forwarded",
	"Content-Length",
	"Transfer-Encoding",
	"Upgrade",
	"Connection",
	"X-Forwarded-For",
	"X-Forwarded-Host",
	"X-Forwarded-Proto",
	"X-Real-Ip",
	"X-Real-Port",
}

type ProxyAPIConfig struct {
	Address       string         `yaml:"address"`
	AddPrefixPath string         `yaml:"addprefixpath"`
	SubPrefixPath string         `yaml:"subprefixpath"`
	RewriteReg    string         `yaml:"rewritereg"`
	RewriteTarget string         `yaml:"rewritetarget"`
	Header        []HeaderConfig `yaml:"header"`
	HeaderAdd     []HeaderConfig `yaml:"headeradd"`
	HeaderDel     []string       `yaml:"headerdel"`
	Query         []QueryConfig  `yaml:"query"`
	QueryAdd      []QueryConfig  `yaml:"queryadd"`
	QueryDel      []string       `yaml:"querydel"`
	Via           string         `yaml:"via"`
}

const defaultVia = "huan-proxy"

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

		if p.Via == "" {
			p.Via = defaultVia
		}
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

		for _, h := range p.Header {
			if h.Header == "" {
				return NewConfigError("header name is empty")
			}

			if h.Header == ViaHeader || h.Header == XHuanProxyHeaer {
				return NewConfigError(fmt.Sprintf("header %s use by http system", h.Header))
			}

			if !utils.IsValidHTTPHeaderKey(h.Header) {
				return NewConfigError(fmt.Sprintf("header %s is not valid", h.Header))
			}

			if isNotGoodHeader(h.Header) {
				_ = NewConfigWarning(fmt.Sprintf("header %s use by http system", h.Header))
			}

			if h.Value == "" {
				_ = NewConfigWarning(fmt.Sprintf("the value of header %s is empty, but maybe it is not delete from requests", h.Header))
			}
		}

		for _, h := range p.HeaderAdd {
			if h.Header == "" {
				return NewConfigError("header name is empty")
			}

			if h.Header == ViaHeader || h.Header == XHuanProxyHeaer {
				return NewConfigError(fmt.Sprintf("header %s use by http system", h.Header))
			}

			if !utils.IsValidHTTPHeaderKey(h.Header) {
				return NewConfigError(fmt.Sprintf("header %s is not valid", h.Header))
			}

			if isNotGoodHeader(h.Header) {
				_ = NewConfigWarning(fmt.Sprintf("header %s use by http system", h.Header))
			}

			if h.Value == "" {
				_ = NewConfigWarning(fmt.Sprintf("the value of header %s is empty, but maybe it is not delete from requests", h.Header))
			}
		}

		for _, h := range p.HeaderDel {
			if h == "" {
				return NewConfigError("header name is empty")
			}

			if h == ViaHeader || h == XHuanProxyHeaer {
				return NewConfigError(fmt.Sprintf("header %s use by http system", h))
			}

			if !utils.IsValidHTTPHeaderKey(h) {
				return NewConfigError(fmt.Sprintf("header %s is not valid", h))
			}

			if isNotGoodHeader(h) {
				_ = NewConfigWarning(fmt.Sprintf("header %s use by http system", h))
			}
		}

		for _, q := range p.Query {
			if q.Query == "" {
				return NewConfigError("query key is empty")
			}

			if !utils.IsGoodQueryKey(q.Query) {
				_ = NewConfigWarning(fmt.Sprintf("query %s is not good", q.Query))
			}

			if q.Value == "" {
				_ = NewConfigWarning(fmt.Sprintf("the value of query %s is empty, but maybe it is not delete from requests", q.Query))
			}
		}

		for _, q := range p.QueryAdd {
			if q.Query == "" {
				return NewConfigError("query key is empty")
			}

			if !utils.IsGoodQueryKey(q.Query) {
				_ = NewConfigWarning(fmt.Sprintf("query %s is not good", q.Query))
			}

			if q.Value == "" {
				_ = NewConfigWarning(fmt.Sprintf("the value of query %s is empty, but maybe it is not delete from requests", q.Query))
			}
		}

		for _, q := range p.QueryDel {
			if q == "" {
				return NewConfigError("query key is empty")
			}

			if !utils.IsGoodQueryKey(q) {
				_ = NewConfigWarning(fmt.Sprintf("query %s is not good", q))
			}
		}

	} else {
		return NewConfigError("proxy type must be file or dir or api")
	}
	return nil
}

func isNotGoodHeader(header string) bool {
	for _, h := range WarningHeader {
		if h == header {
			return true
		}
	}

	return false
}
