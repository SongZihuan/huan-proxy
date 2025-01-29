package api

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/utils"
)

type QueryConfig struct {
	Query string `yaml:"query"`
	Value string `yaml:"value"`
}

func (q *QueryConfig) SetDefault() {

}

func (q *QueryConfig) Check() configerr.ConfigError {
	if q.Query == "" {
		return configerr.NewConfigError("query key is empty")
	}

	if !utils.IsGoodQueryKey(q.Query) {
		_ = configerr.NewConfigWarning(fmt.Sprintf("query %s is not good", q.Query))
	}

	if q.Value == "" {
		_ = configerr.NewConfigWarning(fmt.Sprintf("the value of query %s is empty, but maybe it is not delete from requests", q.Query))
	}

	return nil
}
