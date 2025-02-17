package corscompile

import (
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/cors"
	"regexp"
	"strings"
)

const CorsMaxAgeSec = cors.CorsMaxAgeSec
const CorsDefaultMaxAgeSec = CorsMaxAgeSec
const AllowAllOrigin = "*"

type CorsCompileConfig struct {
	Ignore         bool
	AllowOrigin    []string
	AllowOriginReg []*regexp.Regexp
	MaxAgeSec      int
}

func NewCorsCompileConfig(c *cors.CorsConfig) (*CorsCompileConfig, error) {
	if c.AllowCors.IsDisable(false) {
		return &CorsCompileConfig{
			Ignore:         true,
			AllowOrigin:    make([]string, 0),
			AllowOriginReg: make([]*regexp.Regexp, 0),
			MaxAgeSec:      0,
		}, nil
	}

	regexps := make([]*regexp.Regexp, 0, len(c.AllowOriginRegex))
	for _, v := range c.AllowOriginRegex {
		reg, err := regexp.Compile(v)
		if err != nil {
			return nil, err
		}
		regexps = append(regexps, reg)
	}

	res := &CorsCompileConfig{
		Ignore:         false,
		AllowOrigin:    c.AllowOrigin,
		AllowOriginReg: regexps,
		MaxAgeSec:      c.MaxAgeSec,
	}

	if res.MaxAgeSec >= CorsMaxAgeSec {
		res.MaxAgeSec = CorsMaxAgeSec
	} else if res.MaxAgeSec < 0 {
		res.MaxAgeSec = CorsDefaultMaxAgeSec
	}

	return res, nil
}

func (c *CorsCompileConfig) InOriginList(origin string) bool {
	if len(c.AllowOrigin) == 0 && len(c.AllowOriginReg) == 0 {
		return false
	}

	origin = strings.TrimSpace(origin)

	for _, org := range c.AllowOrigin {
		org = strings.TrimSpace(org)
		if org == AllowAllOrigin {
			return true
		} else if org == origin {
			return true
		}
	}

	for _, reg := range c.AllowOriginReg {
		if reg != nil && reg.MatchString(origin) {
			return true
		}
	}
	return false
}
