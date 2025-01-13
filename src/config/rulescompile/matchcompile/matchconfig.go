package matchcompile

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config/rules/match"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"regexp"
)

const (
	RulesPrefixMatch    = match.PrefixMatch    // 前缀匹配
	RulesRegexMatch     = match.RegexMatch     // 正则匹配
	RulesPrecisionMatch = match.PrecisionMatch // 精准匹配
)

const (
	PrefixMatch    = iota // 前缀匹配
	RegexMatch            // 正则匹配
	PrecisionMatch        // 精准匹配
)

var MatchTypeMap = map[string]int{
	RulesPrefixMatch:    PrefixMatch,
	RulesRegexMatch:     RegexMatch,
	RulesPrecisionMatch: PrecisionMatch,
}

type MatchCompileConfig struct {
	MatchType  int
	MatchPath  string         // Prefix和Precision使用
	MatchRegex *regexp.Regexp // regex使用
}

func NewMatchConfig(m *match.MatchConfig) (*MatchCompileConfig, error) {
	res := new(MatchCompileConfig)

	matchType, ok := MatchTypeMap[m.MatchType]
	if !ok {
		return nil, fmt.Errorf("bad match type")
	}

	res.MatchType = matchType

	if matchType == RegexMatch {
		reg, err := regexp.Compile(m.Path)
		if err != nil {
			return nil, err
		}
		res.MatchRegex = reg
	} else if matchType == PrefixMatch || matchType == PrecisionMatch {
		if !utils.IsValidURLPath(m.Path) {
			return nil, fmt.Errorf("bad path")
		}

		res.MatchPath = utils.ProcessPath(m.Path)
	} else {
		return nil, fmt.Errorf("bad match type")
	}

	return res, nil
}
