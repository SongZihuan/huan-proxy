package rulescompile

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config/rules"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/apicompile"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/dircompile"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/filecompile"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/remotetrustcompile"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/matchcompile"
)

const (
	RulesProxyTypeFile = rules.ProxyTypeFile
	RulesProxyTypeDir  = rules.ProxyTypeDir
	RulesProxyTypeAPI  = rules.ProxyTypeAPI
)

const (
	ProxyTypeFile = iota
	ProxyTypeDir
	ProxyTypeAPI
)

var ProxyTypeMap = map[string]int{
	RulesProxyTypeFile: ProxyTypeFile,
	RulesProxyTypeDir:  ProxyTypeDir,
	RulesProxyTypeAPI:  ProxyTypeAPI,
}

type RuleCompileConfig struct {
	Type int

	*matchcompile.MatchCompileConfig
	*remotetrustcompile.RemoteTrustCompileConfig

	File *filecompile.RuleFileCompileConfig
	Dir  *dircompile.RuleDirCompileConfig
	Api  *apicompile.RuleAPICompileConfig
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

	if typeID == ProxyTypeFile {
		file, err := filecompile.NewRuleFileCompileConfig(&r.RuleFileConfig)
		if err != nil {
			return nil, err
		}

		return &RuleCompileConfig{
			Type:                     typeID,
			MatchCompileConfig:       match,
			RemoteTrustCompileConfig: remoteTrusts,
			File:                     file,
		}, nil
	} else if typeID == ProxyTypeDir {
		dir, err := dircompile.NewRuleDirCompileConfig(&r.RuleDirConfig)
		if err != nil {
			return nil, err
		}

		return &RuleCompileConfig{
			Type:                     typeID,
			MatchCompileConfig:       match,
			RemoteTrustCompileConfig: remoteTrusts,
			Dir:                      dir,
		}, nil
	} else if typeID == ProxyTypeAPI {
		api, err := apicompile.NewRuleAPICompileConfig(&r.RuleAPIConfig)
		if err != nil {
			return nil, err
		}

		return &RuleCompileConfig{
			Type:                     typeID,
			MatchCompileConfig:       match,
			RemoteTrustCompileConfig: remoteTrusts,
			Api:                      api,
		}, nil
	} else {
		return nil, fmt.Errorf("error rule type")
	}
}
