package yaml

import (
	"aria2_cbt_tracker_updater/pkg/util"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func InitConfig(filePath string, config interface{}) {
	if !util.FileExists(filePath) {
		_ = fmt.Errorf("%v is not a valid file path", filePath)
	}
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Println(err.Error())
	}
}

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

func CreateConfig(filePath string, config interface{}) {

}
