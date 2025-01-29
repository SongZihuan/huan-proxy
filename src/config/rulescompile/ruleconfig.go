package rulescompile

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config/rules"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/apicompile"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/dircompile"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/filecompile"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/redirectcompile"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/remotetrustcompile"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/respheadercompile"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/matchcompile"
)

const (
	RulesProxyTypeFile     = rules.ProxyTypeFile
	RulesProxyTypeDir      = rules.ProxyTypeDir
	RulesProxyTypeAPI      = rules.ProxyTypeAPI
	RulesProxyTypeRedirect = rules.ProxyTypeRedirect
)

const (
	ProxyTypeFile = iota
	ProxyTypeDir
	ProxyTypeAPI
	ProxyTypeRedirect
)

var ProxyTypeMap = map[string]int{
	RulesProxyTypeFile:     ProxyTypeFile,
	RulesProxyTypeDir:      ProxyTypeDir,
	RulesProxyTypeAPI:      ProxyTypeAPI,
	RulesProxyTypeRedirect: ProxyTypeRedirect,
}

type RuleCompileConfig struct {
	Type int

	*matchcompile.MatchCompileConfig
	*remotetrustcompile.RemoteTrustCompileConfig

	File     *filecompile.RuleFileCompileConfig
	Dir      *dircompile.RuleDirCompileConfig
	Api      *apicompile.RuleAPICompileConfig
	Redirect *redirectcompile.RuleRedirectCompileConfig

	RespHeader *respheadercompile.SetRespHeaderCompileConfig
}

func NewRuleCompileConfig(r *rules.RuleConfig) (*RuleCompileConfig, error) {
	typeID, ok := ProxyTypeMap[r.Type]
	if !ok {
		return nil, fmt.Errorf("error rule type")
	}

	match, err := matchcompile.NewMatchConfig(&r.MatchConfig)
	if err != nil {
		return nil, err
	}

	remoteTrusts, err := remotetrustcompile.NewRemoteTrustCompileConfig(&r.RemoteTrustConfig)
	if err != nil {
		return nil, err
	}

	respHeader, err := respheadercompile.NewSetRespHeaderCompileConfig(&r.RespHeader)
	if err != nil {
		return nil, err
	}

	if typeID == ProxyTypeFile {
		file, err := filecompile.NewRuleFileCompileConfig(&r.File)
		if err != nil {
			return nil, err
		}

		return &RuleCompileConfig{
			Type:                     typeID,
			MatchCompileConfig:       match,
			RemoteTrustCompileConfig: remoteTrusts,
			File:                     file,
			RespHeader:               respHeader,
		}, nil
	} else if typeID == ProxyTypeDir {
		dir, err := dircompile.NewRuleDirCompileConfig(&r.Dir)
		if err != nil {
			return nil, err
		}

		return &RuleCompileConfig{
			Type:                     typeID,
			MatchCompileConfig:       match,
			RemoteTrustCompileConfig: remoteTrusts,
			Dir:                      dir,
			RespHeader:               respHeader,
		}, nil
	} else if typeID == ProxyTypeAPI {
		api, err := apicompile.NewRuleAPICompileConfig(&r.Api)
		if err != nil {
			return nil, err
		}

		return &RuleCompileConfig{
			Type:                     typeID,
			MatchCompileConfig:       match,
			RemoteTrustCompileConfig: remoteTrusts,
			Api:                      api,
			RespHeader:               respHeader,
		}, nil
	} else if typeID == ProxyTypeRedirect {
		redirect, err := redirectcompile.NewRuleAPICompileConfig(&r.Redirect)
		if err != nil {
			return nil, err
		}

		return &RuleCompileConfig{
			Type:                     typeID,
			MatchCompileConfig:       match,
			RemoteTrustCompileConfig: remoteTrusts,
			Redirect:                 redirect,
			RespHeader:               respHeader,
		}, nil
	} else {
		return nil, fmt.Errorf("error rule type")
	}
}
