package dircompile

import (
	"github.com/SongZihuan/huan-proxy/src/config/rules/action/dir"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/corscompile"
	"github.com/SongZihuan/huan-proxy/src/config/rulescompile/actioncompile/rewritecompile"
)

type RuleDirCompileConfig struct {
	BasePath   string
	IndexFile  []*IndexFileCompileConfig
	IgnoreFile []*IgnoreFileCompileConfig
	AddPath    string
	SubPath    string
	Rewrite    *rewritecompile.RewriteCompileConfig
	Cors       *corscompile.CorsCompileConfig
}

func NewRuleDirCompileConfig(r *dir.RuleDirConfig) (*RuleDirCompileConfig, error) {
	Index := make([]*IndexFileCompileConfig, 0, len(r.IndexFile))
	for _, i := range r.IndexFile {
		file, err := NewIndexFileCompileConfig(i)
		if err != nil {
			return nil, err
		}
		Index = append(Index, file)
	}

	Ignore := make([]*IgnoreFileCompileConfig, 0, len(r.IgnoreFile))
	for _, i := range r.IgnoreFile {
		file, err := NewIgnoreFileCompileConfig(i)
		if err != nil {
			return nil, err
		}
		Ignore = append(Ignore, file)
	}

	rewrite, err := rewritecompile.NewRewriteCompileConfig(&r.Rewrite)
	if err != nil {
		return nil, err
	}

	cors, err := corscompile.NewCorsCompileConfig(&r.Cors)
	if err != nil {
		return nil, err
	}

	return &RuleDirCompileConfig{
		BasePath:   r.BasePath,
		IndexFile:  Index,
		IgnoreFile: Ignore,
		AddPath:    r.AddPath,
		SubPath:    r.SubPath,
		Rewrite:    rewrite,
		Cors:       cors,
	}, nil
}

func (i *IgnoreFileCompileConfig) CheckName(name string) bool {
	if i.IsRegex {
		return i.Regex.MatchString(name)
	} else {
		return name == i.File
	}
}
