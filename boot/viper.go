package boot

import (
	"fmt"
	"github.com/reggiepy/aria2c_bt_updater/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Viper() (*viper.Viper, error) {
	configFile := viper.GetString("config")
	if configFile == "" {
		configFile = "config.yaml"
	}
	v := viper.New()
	v.SetConfigFile(configFile)
	v.AddConfigPath(".")
	v.SetEnvPrefix("ABU") // 设置环境变量前缀
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Config file not found：%s \n", err.Error())
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed: ", e.String())
		err := BindConfig(v)
		if err != nil {
			fmt.Printf("Config file changed, but failed to bind: %v", err.Error())
		} else {
			fmt.Println("Config file changed: ", global.Config.ToJson())
		}
	})
	err := BindConfig(v)
	if err != nil {
		return nil, err
	}
	allSettings := v.AllSettings()
	fmt.Printf("Current viper settings: %+v\n", allSettings)
	return v, nil
}

func BindConfig(v *viper.Viper) error {
	if err := v.Unmarshal(&global.Config); err != nil {
		return fmt.Errorf("Failed to bind config file：%s \n", err)
	}
	return nil
}
