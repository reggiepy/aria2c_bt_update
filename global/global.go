package global

import (
	"github.com/reggiepy/aria2c_bt_updater/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	Config = config.DefaultConfig()

	Viper *viper.Viper

	Logger        *zap.Logger
	LoggerCleanup func()
)
