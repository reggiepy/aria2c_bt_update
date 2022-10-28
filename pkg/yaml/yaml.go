package yaml

import (
	"aria2_cbt_tracker_updater/pkg/util"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func InitConfigByViper(filePath string, config interface{}) {
	if !util.FileExists(filePath) {
		_ = fmt.Errorf("%v is not a valid file path", filePath)
	}
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filePath)

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func CreateConfig(filePath string, config interface{}) (err error) {
	if util.FileExists(filePath) {
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
