package config

import (
	"fmt"
	"regexp"
)

const DefaultRewriteConfigMapSize = 20

type RewriteConfigCompile struct {
	Index  int
	Reg    *regexp.Regexp
	Target string
}

func NewRewriteConfigCompile(index int, r string, target string) (*RewriteConfigCompile, error) {
	reg, err := regexp.Compile(r)
	if err != nil {
		return nil, err
	}

	return &RewriteConfigCompile{
		Index:  index,
		Reg:    reg,
		Target: target,
	}, nil
}

type RewriteConfigCompileList struct {
	Map map[int]*RewriteConfigCompile
}

func (i *RewriteConfigCompileList) init() error {
	i.Map = make(map[int]*RewriteConfigCompile, DefaultRewriteConfigMapSize)
	return nil
}

func (i *RewriteConfigCompileList) Add(ruleIndex int, r string, target string) error {
	rewrite, err := NewRewriteConfigCompile(ruleIndex, r, target)
	if err != nil {
		return err
	}

	res, ok := i.Map[ruleIndex]
	if ok || res != nil {
		return fmt.Errorf("rule exists")
	}

	i.Map[ruleIndex] = rewrite
	return nil
}

func (i *RewriteConfigCompileList) Rewrite(ruleIndex int, path string) (string, error) {
	rule, ok := i.Map[ruleIndex]
	if !ok || rule == nil {
		return "", fmt.Errorf("rule not found")
	}

	res := rule.Reg.ReplaceAllString(path, rule.Target)
	return res, nil
}
