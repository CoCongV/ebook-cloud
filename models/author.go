package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

//Country Model
type Country struct {
	gorm.Model
	Name string `gorm:"type:varchar(32);unique_index"`
}

func (c Country) String() string {
	return fmt.Sprintf("Country<%d %s>", c.ID, c.Name)
}

// Author Model
type Author struct {
	gorm.Model
	Name     string  `gorm:"type:varchar(32);unique_index;not null"`
	Country  Country `gorm:"foreignkey:Name"`
	AuthorID uint
}

func (a Author) String() string {
	return fmt.Sprintf("Author<%d %s %s>", a.ID, a.Name, a.Country.Name)
}
