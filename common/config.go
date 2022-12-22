package common

import (
	"time"
)

var (
	DefaultConfig = "app:\n  host: 127.0.0.1\n  port: 6800\n  token:\nlog:\n  LogFile: bt-updater.log\n  LogWay: console\n  LogLevel: info\n  LogMaxDays: 0\n  DisableLogColor: false\nbtTrackerUrl: https://cdn.staticaly.com/gh/XIU2/TrackersListCollection/master/best_aria2.txt\nhttpProxy:\nfrequency: 60\n"
)

type Config struct {
	App          *App          `yaml:"app"`
	Log          *Log          `yaml:"log"`
	HttpProxy    string        `yaml:"http_proxy,inline"`
	BtTrackerUrl string        `yaml:"bt_tracker_url,inline"`
	Frequency    time.Duration `yaml:"frequency,inline"`
}

type App struct {
	Host  string `yaml:"host"`
	Port  int    `yaml:"port"`
	Token string `yaml:"token"`
}

type Log struct {
	// LogFile specifies a file where logs will be written to. This value will
	// only be used if LogWay is set appropriately. By default, this value is
	// "console".
	LogFile string `yaml:"log_file"`
	// LogWay specifies the way logging is managed. Valid values are "console"
	// or "file". If "console" is used, logs will be printed to stdout. If
	// "file" is used, logs will be printed to LogFile. By default, this value
	// is "console".
	LogWay string `yaml:"frequency,log_way"`
	// LogLevel specifies the minimum log level. Valid values are "trace",
	// "debug", "info", "warn", and "error". By default, this value is "info".
	LogLevel string `yaml:"frequency,log_level"`
	// LogMaxDays specifies the maximum number of days to store log information
	// before deletion. This is only used if LogWay == "file". By default, this
	// value is 0.
	LogMaxDays int64 `yaml:"frequency,log_max_days"`
	// DisableLogColor disables log colors when LogWay == "console" when set to
	// true. By default, this value is false.
	DisableLogColor bool `yaml:"disable_log_color"`
}
