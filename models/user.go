package models

import (
	"github.com/jinzhu/gorm"
)

const (
	//BLACK for Forbid user visit website
	BLACK uint = 0

	//COMMON for user download resource
	COMMON uint = 1

	//RESOURCE for user change resource
	RESOURCE uint = 8
)

//Role for control user permission
type Role struct {
	gorm.Model
	Name        string `gorm:"type:varchar(32);unique_index`
	Permmission uint
	Users       []User
}

//User for store user info
type User struct {
	gorm.Model
	UID    uint
	RoleID uint
}
