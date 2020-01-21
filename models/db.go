package models

import (
	"github.com/jinzhu/gorm"
)

//DB is ...
var DB *gorm.DB

//InitDB is init gorm db connection
func InitDB(url string) *gorm.DB {
	DB, err := gorm.Open("postgres", url)
	if err != nil {
		panic(err)
	}
	return DB
}
