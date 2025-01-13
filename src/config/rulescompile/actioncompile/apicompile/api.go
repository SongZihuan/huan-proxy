package apicompile

import (
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/api"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/rewritecompile"
	"net/http/httputil"
	"net/url"
)

const XHuanProxyHeaer = api.XHuanProxyHeaer
const ViaHeader = api.ViaHeader

type RuleAPICompileConfig struct {
	Address       string
	TargetURL     *url.URL
	Server        *httputil.ReverseProxy
	AddPrefixPath string
	SubPrefixPath string
	Rewrite       *rewritecompile.RewriteCompileConfig
	HeaderSet     []*HeaderCompileConfig
	HeaderAdd     []*HeaderCompileConfig
	HeaderDel     []*HeaderDelCompileConfig
	QuerySet      []*QueryCompileConfig
	QueryAdd      []*QueryCompileConfig
	QueryDel      []*QueryDelCompileConfig
	Via           string
}

func NewRuleAPICompileConfig(r *api.RuleAPIConfig) (*RuleAPICompileConfig, error) {
	rewrite, err := rewritecompile.NewRewriteCompileConfig(&r.ApiRewrite)
	if err != nil {
		return nil, err
	}

	targetURL, err := url.Parse(r.Address)
	if err != nil {
		return nil, err
	}

	server := httputil.NewSingleHostReverseProxy(targetURL)

	HeaderSet := make([]*HeaderCompileConfig, 0, len(r.HeaderSet))
	for _, v := range r.HeaderSet {
		h, err := NewHeaderCompileConfig(&v)
		if err != nil {
			return nil, err
		}
		HeaderSet = append(HeaderSet, h)
	}

	HeaderAdd := make([]*HeaderCompileConfig, 0, len(r.HeaderAdd))
	for _, v := range r.HeaderAdd {
		h, err := NewHeaderCompileConfig(&v)
		if err != nil {
			return nil, err
		}
		HeaderAdd = append(HeaderAdd, h)
	}

	HeaderDel := make([]*HeaderDelCompileConfig, 0, len(r.HeaderDel))
	for _, v := range r.HeaderDel {
		h, err := NewHeaderDelCompileConfig(&v)
		if err != nil {
			return nil, err
		}
		HeaderDel = append(HeaderDel, h)
	}

	QuerySet := make([]*QueryCompileConfig, 0, len(r.QuerySet))
	for _, v := range r.QuerySet {
		q, err := NewQueryCompileConfig(&v)
		if err != nil {
			return nil, err
		}
		QuerySet = append(QuerySet, q)
	}

	QueryAdd := make([]*QueryCompileConfig, 0, len(r.QueryAdd))
	for _, v := range r.QueryAdd {
		q, err := NewQueryCompileConfig(&v)
		if err != nil {
			return nil, err
		}
		QueryAdd = append(QueryAdd, q)
	}

	QueryDel := make([]*QueryDelCompileConfig, 0, len(r.QueryDel))
	for _, v := range r.QueryDel {
		q, err := NewQueryDelCompileConfig(&v)
		if err != nil {
			return nil, err
		}
		QueryDel = append(QueryDel, q)
	}

	return &RuleAPICompileConfig{
		Address:       r.Address,
		TargetURL:     targetURL,
		Server:        server,
		AddPrefixPath: r.ApiAddPrefixPath,
		SubPrefixPath: r.ApiSubPrefixPath,
		Rewrite:       rewrite,
		HeaderSet:     HeaderSet,
		HeaderAdd:     HeaderAdd,
		HeaderDel:     HeaderDel,
		QuerySet:      QuerySet,
		QueryAdd:      QueryAdd,
		Via:           r.Via,
	}, nil
}
