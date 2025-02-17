package huanproxy

import (
	"errors"
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config"
	"github.com/SongZihuan/huan-proxy/src/config/configwatcher"
	"github.com/SongZihuan/huan-proxy/src/flagparser"
	"github.com/SongZihuan/huan-proxy/src/logger"
	"github.com/SongZihuan/huan-proxy/src/server"
	"github.com/SongZihuan/huan-proxy/src/server/httpserver"
	"github.com/SongZihuan/huan-proxy/src/server/httpsserver"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"os"
	"time"
)

func MainV1() int {
	var err error

	err = flagparser.InitFlag()
	if errors.Is(err, flagparser.StopFlag) {
		return 0
	} else if err != nil {
		return utils.ExitByError(err)
	}

	if !flagparser.IsReady() {
		return utils.ExitByErrorMsg("flag parser unknown error")
	}

	utils.SayHellof("%s", "The backend service program starts normally, thank you.")
	defer func() {
		utils.SayGoodByef("%s", "The backend service program is offline/shutdown normally, thank you.")
	}()

	cfgErr := config.InitConfig(flagparser.ConfigFile())
	if cfgErr != nil && cfgErr.IsError() {
		return utils.ExitByError(cfgErr)
	}

	if !config.IsReady() {
		return utils.ExitByErrorMsg("config parser unknown error")
	}

	err = logger.InitLogger(os.Stdout, os.Stderr)
	if err != nil {
		return utils.ExitByError(err)
	}

	if !logger.IsReady() {
		return utils.ExitByErrorMsg("logger unknown error")
	}

	if flagparser.RunAutoReload() {
		err = configwatcher.WatcherConfigFile()
		if err != nil {
			return utils.ExitByError(err)
		}
		defer configwatcher.CloseNotifyConfigFile()

		logger.Infof("Auto reload enable.")
	} else {
		logger.Infof("Auto reload disable.")
	}

	logger.Executablef("%s", "ready")
	logger.Infof("run mode: %s", config.GetConfig().GlobalConfig.GetRunMode())

	ser := server.NewHuanProxyServer()

	httpErrorChan := make(chan error)
	httpsErrorChan := make(chan error)

	err = ser.Run(httpErrorChan, httpsErrorChan)
	if err != nil {
		return utils.ExitByErrorMsg(fmt.Sprintf("run http/https error: %s", err.Error()))
	}
	defer func() {
		_ = ser.Stop()
		time.Sleep(1 * time.Second)
	}()

	select {
	case <-config.GetSignalChan():
		return 0
	case err := <-httpErrorChan:
		if errors.Is(err, httpserver.ServerStop) {
			return 0
		} else if err != nil {
			return utils.ExitByError(err)
		} else {
			return 0
		}
	case err := <-httpsErrorChan:
		if errors.Is(err, httpsserver.ServerStop) {
			return 0
		} else if err != nil {
			return utils.ExitByError(err)
		} else {
			return 0
		}
	}
}
