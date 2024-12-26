package sub

import (
	"fmt"
	"os"

	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/jsonutil"
	"github.com/reggiepy/aria2c_bt_updater/boot"
	"github.com/reggiepy/aria2c_bt_updater/config"
	"github.com/reggiepy/aria2c_bt_updater/global"
	"github.com/reggiepy/aria2c_bt_updater/pkg/goutils/enumUtils"
	"github.com/reggiepy/aria2c_bt_updater/pkg/goutils/yamlutil"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

type ConfigConfig struct {
	Format *enumUtils.Enum
	Force  bool
	Config string
}

var configConfig = ConfigConfig{
	Format: enumUtils.NewEnum([]string{"humanReadable", "simple"}, "humanReadable"),
}

func init() {
	configShowCmd.Flags().Var(configConfig.Format, "format", "humanReadable | simple")
	configShowCmd.Flags().StringVarP(&configConfig.Config, "config", "c", "./config.yaml", "config file")
	_ = viper.BindPFlag("config", configGenerateCmd.PersistentFlags().Lookup("config"))
	configCmd.AddCommand(configShowCmd)

	configGenerateCmd.Flags().BoolVarP(&configConfig.Force, "force", "f", false, "Generate configuration forces")
	configGenerateCmd.Flags().StringVarP(&configConfig.Config, "config", "c", "", "config file")
	_ = viper.BindPFlag("config", configGenerateCmd.PersistentFlags().Lookup("config"))
	configCmd.AddCommand(configGenerateCmd)
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "show config",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"json", "simple"}, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		global.Viper, err = boot.Viper()
		if err != nil {
			return err
		}
		configFormat := configConfig.Format.String()

		var data string
		switch configFormat {
		case "humanReadable":
			data, _ = jsonutil.EncodeString(global.Config)
		case "simple":
			dataBytes, _ := jsonutil.Encode(global.Config)
			data = string(dataBytes)
		}
		fmt.Println(data)
		return nil
	},
}

var configGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate default config",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		configFile := viper.GetString("config")
		if configFile == "" {
			return fmt.Errorf("config file not specified")
		}

		configFileExt := fsutil.Extname(configFile)
		defaultConfig := config.DefaultConfig()
		configString := ""
		switch configFileExt {
		case "yaml":
			configString, _ = yamlutil.EncodeString(defaultConfig)
		case "json":
			configString, _ = jsonutil.EncodeString(defaultConfig)
		default:
			return fmt.Errorf("unsupported config file extension: %s", configFileExt)
		}
		flags := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
		if !configConfig.Force {
			flags |= os.O_EXCL
		}
		err = fsutil.WriteFile(configFile, configString, os.ModePerm, flags)
		if err != nil {
			return fmt.Errorf("write config to file failed: %v", err)
		}
		return nil
	},
}
