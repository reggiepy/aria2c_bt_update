package config

func DefaultConfig() *Config {
	return &Config{
		Aria2c: []Aria2c{Aria2c{
			Host:  "127.0.0.1", // 默认绑定到本地地址
			Port:  6800,        // 默认端口号
			Token: "",          // 为空，未配置 token
		}},
		Log: Log{
			File:       "./logs/aria2c_bt_update.log", // 默认日志文件位置
			MaxSize:    5,                             // 最大日志文件大小限制为 5MB
			MaxBackups: 3,                             // 保留最多 3 个旧日志
			MaxAge:     7,                             // 旧日志文件保留 7 天
			Compress:   true,                          // 启用日志压缩
			Level:      "info",                        // 默认日志级别为 info
			Format:     "json",                        // 默认日志格式为 JSON
		},
		System: System{
			HttpProxy: "", // 默认没有设置 HTTP 代理
			//配置来源：https://github.com/XIU2/TrackersListCollection
			//动漫：https://github.com/DeSireFire/animeTrackerList
			BtTrackerUrl: []string{
				//"https://cdn.staticaly.com/gh/XIU2/TrackersListCollection/master/best_aria2.txt",
				"https://cf.trackerslist.com/all.txt",
			}, // 默认 BT Tracker URL
			Frequency: 1, // 默认刷新频率 1 秒
		},
	}
}
