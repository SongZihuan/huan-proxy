package config

import "regexp"

const DefaultOriginListSize = 10
const AllowAllOrigin = "*"

type CorsOrigin struct {
	OriginReg    []*regexp.Regexp
	OriginString []string
}

func (c *CorsOrigin) init() error {
	c.OriginReg = make([]*regexp.Regexp, 0, DefaultOriginListSize)
	c.OriginString = nil
	return nil
}

func (c *CorsOrigin) ApplyReg(origin string) error {
	reg, err := regexp.Compile(origin)
	if err != nil {
		return err
	}
	c.OriginReg = append(c.OriginReg, reg)
	return nil
}

func (c *CorsOrigin) SetString(origins []string) error {
	c.OriginString = origins
	return nil
}

func (c *CorsOrigin) InOriginList(origin string) bool {
	if (c.OriginString == nil || len(c.OriginString) == 0) && len(c.OriginReg) == 0 {
		return true
	}

	for _, org := range c.OriginString {
		if org == AllowAllOrigin {
			return true
		} else if org == origin {
			return true
		}
	}

	for _, reg := range c.OriginReg {
		if reg != nil && reg.MatchString(origin) {
			return true
		}
	}
	return false
}
