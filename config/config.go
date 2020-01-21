package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

//Config struct
type Config struct {
	Addr  string `toml:"addr"`
	DBURL string `toml:"dbURL"`
}

//ReadConfig is read toml config file
func ReadConfig(filePath string) *Config {
	var config Config

	_, err := toml.DecodeFile(filePath, &config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}
