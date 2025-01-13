package apicompile

import "github.com/SongZihuan/huan-proxy/src/config/rules/action/api"

type QueryCompileConfig struct {
	Query string `yaml:"query"`
	Value string `yaml:"value"`
}

func NewQueryCompileConfig(q *api.QueryConfig) (*QueryCompileConfig, error) {
	return &QueryCompileConfig{
		Query: q.Query,
		Value: q.Value,
	}, nil
}
