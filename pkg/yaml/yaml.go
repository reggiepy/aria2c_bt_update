package yaml

import (
	"aria2c_bt_updater/common"
	util2 "aria2c_bt_updater/pkg/util/util"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func InitConfigByViper(filePath string, config interface{}) (err error) {
	filePath, _ = util2.AbsFilePath(filePath)
	if !util2.FileExists(filePath) {
		fmt.Println(fmt.Errorf("%v is not a valid file path", filePath))
		err = CreateConfig(filePath, common.DefaultConfig)
		if err != nil {
			return fmt.Errorf("create config error: %v", err)
		}
	}
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filePath)

	err = viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read config file error: %v", err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return fmt.Errorf("parse config error: %v", err)
	}
	return
}

func CreateConfig(filePath string, config interface{}) (err error) {
	filePath, _ = util2.AbsFilePath(filePath)
	fmt.Println(filePath)
	if util2.FileExists(filePath) {
		return fmt.Errorf("%v already exists", filePath)
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("open file error: %v", err)
	}
	defer file.Close()
	_, err = file.Write([]byte(config.(string)))
	if err != nil {
		return fmt.Errorf("write file error: %v", err)
	}
	return nil
}
