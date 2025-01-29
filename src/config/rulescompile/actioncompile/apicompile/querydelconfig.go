package apicompile

import "github.com/SongZihuan/huan-proxy/src/config/rules/action/api"

type QueryDelCompileConfig struct {
	Query string `yaml:"query"`
}

func NewQueryDelCompileConfig(q *api.QueryDelConfig) (*QueryDelCompileConfig, error) {
	return &QueryDelCompileConfig{
		Query: q.Query,
	}, nil
}
