package sub

import (
	"fmt"
	"github.com/gookit/goutil/fsutil"
	"github.com/reggiepy/aria2c_bt_updater/pkg/goutils/signailUtils"
	"os"

	"github.com/reggiepy/aria2c_bt_updater/boot"
	"github.com/reggiepy/aria2c_bt_updater/global"
	"github.com/reggiepy/aria2c_bt_updater/pkg/server"
	"github.com/reggiepy/aria2c_bt_updater/pkg/version"

	"github.com/reggiepy/aria2c_bt_updater/aria2c"
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
		fmt.Println("Config: ", global.Config.ToJson())
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
		jsonRpcOption := aria2c.JsonRpcOption{
			ProxyUrl: global.Config.System.HttpProxy,
		}
		jsonRpc := aria2c.NewJsonRpc(
			global.Config.Aria2c.Host,
			global.Config.Aria2c.Port,
			global.Config.Aria2c.Token,
			jsonRpcOption,
		)
		cfg := server.Config{
			HttpProxy:    global.Config.System.HttpProxy,
			BtTrackerUrl: global.Config.System.BtTrackerUrl,
			Frequency:    global.Config.System.Frequency,
		}
		s := server.NewServer(jsonRpc, cfg)
		go func() {
			s.Run()
		}()
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
