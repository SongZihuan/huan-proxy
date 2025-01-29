package api

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config/configerr"
	"github.com/SongZihuan/huan-proxy/src/utils"
)

type QueryDelConfig struct {
	Query string `yaml:"query"`
}

func (q *QueryDelConfig) SetDefault() {

}

func (q *QueryDelConfig) Check() configerr.ConfigError {
	if q.Query == "" {
		return configerr.NewConfigError("query key is empty")
	}

	if !utils.IsGoodQueryKey(q.Query) {
		_ = configerr.NewConfigWarning(fmt.Sprintf("query %s is not good", q.Query))
	}

	return nil
}
