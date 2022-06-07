package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Mysql  MysqlConfig
	Redis  RedisConfig
	Wx     WxConfig
	System SystemConfig
}

type MysqlConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	MaxIdle  int
	MaxOpen  int
}

type RedisConfig struct {
	Host        string
	Port        string
	Username    string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
	Index       int
}

type WxConfig struct {
	Appid     string
	Appsecret string
}

type SystemConfig struct {
	Secret  string
	Start   int64
	End     int64
	Maxvote int
	Port int
}

var Mysql MysqlConfig
var Redis RedisConfig
var Wx WxConfig
var System SystemConfig

func init() {
	var config Config
	viper.SetConfigName("config")
	viper.AddConfigPath("./static/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	viper.Unmarshal(&config)
	Mysql = config.Mysql
	Redis = config.Redis
	Wx = config.Wx
	System = config.System
}
