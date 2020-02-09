package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Book is Model
type Book struct {
	gorm.Model
	Name   string `gorm:"type:varchar(32);unique_index;not null"`
	File   string `gorm:"not null"`
	UserID uint
	Author []*Author `gorm:"many2many:author_books;"`
}

func (b Book) String() string {
	return fmt.Sprintf("Book<%d %s>", b.ID, b.Name)
}
