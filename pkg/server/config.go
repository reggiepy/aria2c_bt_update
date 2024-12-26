package server

import (
	"time"
)

type Config struct {
	HttpProxy    string        `yaml:"http_proxy,inline" json:"http_proxy"`
	BtTrackerUrl []string      `yaml:"bt_tracker_url,inline" json:"bt_tracker_url"`
	Frequency    time.Duration `yaml:"frequency,inline" json:"frequency"`
	Aria2c       Aria2c        `yaml:"aria2c" json:"aria2c"` // Aria2c 配置项
}

type Aria2c struct {
	Host  string `yaml:"host" json:"host"`   // Aria2c 服务的主机地址
	Port  int    `yaml:"port" json:"port"`   // Aria2c 服务的端口
	Token string `yaml:"token" json:"token"` // Aria2c 服务的 Token 用于鉴权
}

func (l *Config) clone() *Config {
	clone := *l
	return &clone
}

func (l *Config) WithOptions(options ...Option) *Config {
	c := l.clone()
	for _, opt := range options {
		opt.apply(c)
	}
	return c
}

func NewConfig(opts ...Option) *Config {
	config := &Config{
		BtTrackerUrl: []string{},
		Frequency:    1,
		Aria2c: Aria2c{
			Host:  "127.0.0.1", // 默认绑定到本地地址
			Port:  6800,        // 默认端口号
			Token: "",          // 为空，未配置 token
		},
	}
	return config.WithOptions(opts...)
}
