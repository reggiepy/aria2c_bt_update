package boot

import (
	"encoding/json"
	"fmt"
	"github.com/reggiepy/aria2c_bt_updater/global"
	"github.com/reggiepy/aria2c_bt_updater/pkg/goutils/logutil/zapLogger"
	"go.uber.org/zap"
)

func Logger() (*zap.Logger, func()) {
	logConfig := zapLogger.NewLoggerConfig(
		zapLogger.WithInConsole(true),
	)
	jsonBytes, _ := json.Marshal(global.Config.Log)
	err := logConfig.LoadJSON(string(jsonBytes))
	if err != nil {
		fmt.Printf("Error marshalling log config, use default config: %v\n", err)
	}
	//fmt.Println("Log Config: ", logConfig.ToJSON())
	logger, cleanup := zapLogger.NewLogger(logConfig)
	return logger, cleanup
}
