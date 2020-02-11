package models

import (
	"fmt"

	"ebook-cloud/config"

	"github.com/jinzhu/gorm"
)

type Tag struct {
	gorm.Model
	Name  string  `gorm:"type:varchar(32);not null;unique_index"`
	Books []*Book `gorm:"many2many:book_tags"`
}

// Book is Model
type Book struct {
	gorm.Model
	Name     string `gorm:"type:varchar(64);not null"`
	File     string `gorm:"not null;type:varchar(128)"`
	UserID   uint
	CoverImg string `gorm:"type:varchar(128)"`
	Describe string
	Authors  []*Author `gorm:"many2many:author_books;"`
	Tags     []*Tag    `gorm:"many2many:book_tags;"`
}

func (b Book) String() string {
	return fmt.Sprintf("Book<%d %s>", b.ID, b.Name)
}

func (b *Book) BeforeSave() error {
	if b.CoverImg == "" {
		b.CoverImg = config.Conf.DefaultCoverImg
	}
	return nil
}
