module aria2_cbt_tracker_updater

go 1.16

require (
	github.com/fatedier/beego v1.7.2
	github.com/reggie/aria2c v0.0.0
	github.com/spf13/viper v1.13.0
)

replace github.com/reggie/aria2c => ./aria2c
