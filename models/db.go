package models

import (
	"github.com/jinzhu/gorm"

	"ebook-cloud/config"
)

//DB is ...
var DB *gorm.DB

//Setup is init gorm db connection
func Setup() {
	var err error
	DB, err = gorm.Open("postgres", config.Conf.DBURL)
	if err != nil {
		panic(err)
	}
}
