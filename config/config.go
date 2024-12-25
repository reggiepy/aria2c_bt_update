package config

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"time"
)

type Config struct {
	Aria2c Aria2c `yaml:"aria2c" json:"aria2c"` // Aria2c 配置项
	Log    Log    `yaml:"log" json:"log"`       // 日志配置项
	System System `yaml:"system" json:"system"` // 系统配置项
}

type Aria2c struct {
	Host  string `yaml:"host" json:"host"`   // Aria2c 服务的主机地址
	Port  int    `yaml:"port" json:"port"`   // Aria2c 服务的端口
	Token string `yaml:"token" json:"token"` // Aria2c 服务的 Token 用于鉴权
}

type System struct {
	HttpProxy    string        `yaml:"http_proxy" json:"http_proxy"`         // 系统的 HTTP 代理
	BtTrackerUrl []string      `yaml:"bt_tracker_url" json:"bt_tracker_url"` // BT Tracker URL 列表
	Frequency    time.Duration `yaml:"frequency" json:"frequency"`           // 配置更新频率（单位：秒）
}

type Log struct {
	File       string `json:"file" yaml:"file"`               // 日志文件路径
	MaxSize    int    `json:"max_size" yaml:"max_size"`       // 日志文件的最大大小（单位：MB）
	MaxBackups int    `json:"max_backups" yaml:"max_backups"` // 最大保留的日志文件数量
	MaxAge     int    `json:"max_age" yaml:"max_age"`         // 旧日志文件的保留天数
	Compress   bool   `json:"compress" yaml:"compress"`       // 是否压缩旧日志文件
	Level      string `json:"log_level" yaml:"log_level"`     // 日志记录的级别（如：debug、info、warn、error）
	Format     string `json:"log_format" yaml:"log_format"`   // 日志格式（如：json、logfmt）
}

// ToJson 将配置转换为JSON格式
func (c *Config) ToJson() string {
	jsonData, _ := json.MarshalIndent(c, "", "  ")
	return string(jsonData)
}

// LoadJson 从JSON文件加载配置
func (c *Config) LoadJson(data string) error {
	if err := json.Unmarshal([]byte(data), c); err != nil {
		return err
	}
	return nil
}

// ToYaml 将配置转换为YAML格式
func (c *Config) ToYaml() (string, error) {
	yamlData, err := yaml.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(yamlData), nil
}

// LoadYaml 从YAML文件加载配置
func (c *Config) LoadYaml(data string) error {
	if err := yaml.Unmarshal([]byte(data), c); err != nil {
		return err
	}
	return nil
}
