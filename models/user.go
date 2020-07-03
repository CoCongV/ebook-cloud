package models

import (
	"github.com/jinzhu/gorm"
)

const (
	//BLACK for Forbid user visit website
	BLACK     uint   = 0
	BlackUser string = "BlackUser"

	//COMMON for user download resource
	COMMON     uint   = 1
	CommonUser string = "CommonUser"

	//RESOURCE for user change resource
	MODERATE  uint   = 8
	Moderator string = "Moderator"

	//ADMIN for administrator
	ADMIN         uint   = 16
	Administrator string = "Administrator"
)

//Role for control user permission
type Role struct {
	gorm.Model
	Name        string `gorm:"type:varchar(32);unique_index"`
	Permmission uint
	Users       []User
}

//HasPermission for check user permission
func (r *Role) HasPermission(perm uint) bool {
	return r.Permmission&perm == perm
}

//NewRoles for create role and administrator
func NewRoles(uid uint) error {
	DB.Create(&Role{
		Name:        BlackUser,
		Permmission: BLACK,
	})
	DB.Create(&Role{
		Name:        CommonUser,
		Permmission: COMMON,
	})
	DB.Create(&Role{
		Name:        Moderator,
		Permmission: MODERATE + COMMON,
	})
	DB.Create(&Role{
		Name:        Administrator,
		Permmission: ADMIN + MODERATE + COMMON,
	}).Association("Users").Append(&User{
		UID: uid,
	})
	return nil
}

//User for store user info
type User struct {
	gorm.Model
	UID    uint `gorm:"unique_index"`
	RoleID uint
}

//AfterSave is hook
func (u *User) AfterSave(scope *gorm.Scope) error {
	if u.RoleID == 0 {
		var role Role
		scope.DB().Where("Name = ?", CommonUser).First(&role).Association("Users").Append(u)
	}
	return nil
}
