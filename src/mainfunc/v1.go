package mainfunc

import (
	"errors"
	"github.com/SongZihuan/huan-proxy/src/config"
	"github.com/SongZihuan/huan-proxy/src/config/configwatcher"
	"github.com/SongZihuan/huan-proxy/src/flagparser"
	"github.com/SongZihuan/huan-proxy/src/logger"
	"github.com/SongZihuan/huan-proxy/src/server"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"os"
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

		err = logger.InitLogger(os.Stdout, os.Stderr)
		if err != nil {
			return utils.ExitByError(err)
		}

		logger.Infof("Auto reload enable.")
	} else {
		logger.Infof("Auto reload disable.")
	}

	logger.Executablef("%s", "ready")
	logger.Infof("run mode: %s", config.GetConfig().GlobalConfig.GetRunMode())

	ser := server.NewServer()

	serstop := make(chan bool)
	sererror := make(chan error)

	go func() {
		err := ser.Run()
		if errors.Is(err, server.ServerStop) {
			serstop <- true
		} else if err != nil {
			sererror <- err
		} else {
			serstop <- false
		}
	}()

	select {
	case <-config.GetSignalChan():
		break
	case err := <-sererror:
		return utils.ExitByError(err)
	case <-serstop:
		break
	}

	return 0
}
