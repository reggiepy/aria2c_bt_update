package server

import "time"

// Option 是一个函数类型，用于修改 Config 配置
type Option interface {
	apply(c *Config)
}

type OptionFunc func(c *Config)

func (f OptionFunc) apply(c *Config) {
	f(c)
}

// WithHttpProxy 配置 HttpProxy 字段
func WithHttpProxy(httpProxy string) Option {
	return OptionFunc(func(c *Config) {
		c.HttpProxy = httpProxy
	})
}

// WithBtTrackerUrl 配置 BtTrackerUrl 字段
func WithBtTrackerUrl(btTrackerUrl []string) Option {
	return OptionFunc(func(c *Config) {
		c.BtTrackerUrl = btTrackerUrl
	})
}

// WithFrequency 配置 Frequency 字段
func WithFrequency(frequency time.Duration) Option {
	return OptionFunc(func(c *Config) {
		c.Frequency = frequency
	})
}

// WithAria2c 配置 Aria2c 字段
func WithAria2c(aria2c Aria2c) Option {
	return OptionFunc(func(c *Config) {
		c.Aria2c = aria2c
	})
}
