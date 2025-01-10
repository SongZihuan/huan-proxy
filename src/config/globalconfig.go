package config

import (
	"github.com/SongZihuan/huan-proxy/src/utils"
	"os"
)

const EnvModeName = "HUAN_PROXY_MODE"

const (
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)

type LoggerLevel string

var levelMap = map[string]bool{
	"debug": true,
	"info":  true,
	"warn":  true,
	"error": true,
	"panic": true,
	"none":  true,
}

type GlobalConfig struct {
	Mode     string           `json:"mode"`
	LogLevel string           `json:"loglevel"`
	LogTag   utils.StringBool `json:"logtag"`
}

func (g *GlobalConfig) setDefault() {
	if g.Mode == "" {
		g.Mode = os.Getenv(EnvModeName)
	}

	if g.Mode == "" {
		g.Mode = DebugMode
	}

	_ = os.Setenv(EnvModeName, g.Mode)

	if g.LogLevel == "" && (g.Mode == DebugMode || g.Mode == TestMode) {
		g.LogLevel = "debug"
		g.LogTag.SetDefaultEanble()
	} else if g.LogLevel == "" {
		g.LogLevel = "warn"
		g.LogTag.SetDefaultDisable()
	}

	return
}

func (g *GlobalConfig) check() ConfigError {
	if g.Mode != DebugMode && g.Mode != ReleaseMode && g.Mode != TestMode {
		return NewConfigError("bad mode")
	}

	if _, ok := levelMap[g.LogLevel]; !ok {
		return NewConfigError("log level error")
	}

	return nil
}

func (g *GlobalConfig) GetGinMode() string {
	return g.Mode
}

func (g *GlobalConfig) IsDebug() bool {
	return g.Mode == DebugMode
}

func (g *GlobalConfig) IsRelease() bool {
	return g.Mode == ReleaseMode
}

func (g *GlobalConfig) IsTest() bool {
	return g.Mode == TestMode
}
