package sub

import (
	"aria2c_bt_updater/common"
	"aria2c_bt_updater/pkg/util/log"
	"aria2c_bt_updater/pkg/util/version"
	"aria2c_bt_updater/server"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	// Used for flags.
	cfgFile     string
	showVersion bool
	rootCmd     = &cobra.Command{
		Use:   "aria2c_bt_update",
		Short: "A generator for windows service script",
		Long:  `aria2c_bt_update is a CLI for update aria2c bt script`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				fmt.Println(version.Full())
				return nil
			}
			config := &common.Config{}
			err := InitConfigByViper(cfgFile, config)
			if err != nil {
				fmt.Println(err)
				return err
			}
			logConfig := config.Log
			log.InitLog(
				logConfig.LogWay,
				logConfig.LogFile,
				logConfig.LogLevel,
				logConfig.LogMaxDays,
				logConfig.DisableLogColor,
			)
			s := server.NewServer(config)
			s.Run()
			return nil
		},
	}
)

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./config.yaml", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	_ = viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version")
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigName("config.yaml")
	}
	viper.AutomaticEnv() // 读取匹配的环境变量

	// 尝试读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}
}

// InitConfigByViper 从 Viper 初始化配置
func InitConfigByViper(cfgFile string, config *common.Config) error {
	err := viper.Unmarshal(config)
	if err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}
	return nil
}
