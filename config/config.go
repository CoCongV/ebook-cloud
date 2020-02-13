package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

//Config struct
type Config struct {
	Addr              string `toml:"addr"`
	DBURL             string `toml:"dbURL"`
	SecretKey         string `toml:"secret_key"`
	Mode              string `toml:"mode"`
	VerifyUserURL     string `toml:"verify_user_url"`
	UserServerTimeout int    `toml:"user_server_timeout"`
	DestPath          string `toml:"dest_path"`
	DefaultCoverImg   string `toml:"default_cover_img"`
	SearchIndexFile   string `toml:"search_index_file"`
	PerPageItem       int    `toml:"per_page_item"`
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
}
