package sub

import (
	"fmt"
	"os"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/jsonutil"
	"github.com/reggiepy/aria2c_bt_updater/pkg/goutils/signailUtils"

	"github.com/reggiepy/aria2c_bt_updater/boot"
	"github.com/reggiepy/aria2c_bt_updater/global"
	"github.com/reggiepy/aria2c_bt_updater/pkg/server"
	"github.com/reggiepy/aria2c_bt_updater/pkg/version"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type GlobalConfig struct {
	ShowVersion bool
}

var globalConfig = &GlobalConfig{}

func init() {
	cobra.OnInitialize(initConfig)
	// 设置全局标志
	rootCmd.PersistentFlags().BoolVarP(&globalConfig.ShowVersion, "version", "v", false, "show version information")
	rootCmd.PersistentFlags().StringP("config", "c", "./config.yaml", "config file")

	// 绑定命令行参数到Viper
	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
}

var rootCmd = &cobra.Command{
	Use:   "aria2c_bt_update",
	Short: "A generator for windows service script",
	Long:  `aria2c_bt_update is a CLI for update aria2c bt script`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalConfig.ShowVersion {
			fmt.Println(version.Full())
			return nil
		}
		var err error
		global.Viper, err = boot.Viper()
		if err != nil {
			return err
		}
		// data, _ := jsonutil.EncodeString(global.Config)
		// fmt.Println("Config: ", data)
		global.Logger, global.LoggerCleanup = boot.Logger()
		signailUtils.OnExit(func() {
			global.LoggerCleanup() // 确保在程序退出时刷新日志缓冲区
		})

		configFile := viper.GetString("config")
		if !fsutil.FileExist(configFile) {
			if err := global.Viper.WriteConfig(); err != nil {
				return err
			}
		}

		clientMap := map[string]struct{}{}
		for idx, aria2cConfig := range global.Config.Aria2c {
			name := fmt.Sprintf("%s:%d", aria2cConfig.Host, aria2cConfig.Port)
			fmt.Printf("%d. Start %s\n", idx+1, name)
			configBytes, _ := jsonutil.Encode(aria2cConfig)
			var config server.Aria2c
			_ = jsonutil.Decode(configBytes, &config)
			cfg := server.NewConfig(
				server.WithHttpProxy(global.Config.System.HttpProxy),
				server.WithBtTrackerUrl(global.Config.System.BtTrackerUrl),
				server.WithFrequency(global.Config.System.Frequency),
				server.WithAria2c(config),
			)
			s := server.NewServer(*cfg)
			clientName := s.ClientName()
			if _, ok := clientMap[clientName]; ok {
				fmt.Println("client already exists")
				continue
			}
			clientMap[clientName] = struct{}{}
			go func() {
				s.Run()
			}()
		}
		boot.Boot()
		return nil
	},
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initConfig() {
}
