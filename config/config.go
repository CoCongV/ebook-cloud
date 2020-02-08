package config

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
)

//Config struct
type Config struct {
	Addr              string `toml:"addr"`
	DBURL             string `toml:"dbURL"`
	SecretKey         string `toml:"secret_key"`
	ExpiresAt         int64  `toml:"expires_at"`
	Mode              string `toml:"mode"`
	UserServerURL     string `toml:"user_server_url"`
	UserServerTimeout int    `toml:"user_server_timeout"`
}

//Conf is struct Config point
var Conf *Config

//Setup is read toml config file for init
func Setup() {
	filepath := "./config.toml"
	_, err := toml.DecodeFile(filepath, &Conf)
	if err != nil {
		log.Fatal(err)
	}
	switch Conf.Mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}
