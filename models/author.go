package models

import (
	"fmt"
)

type Country struct {
	Id   int64
	Name string
}

func (c Country) String() string {
	return fmt.Sprintf("Country<%d %s>", u.Id, u.Name)
}

type Author struct {
	Id      int64
	Name    string
	Country *Country
}

func (A Author) String() string {
	return fmt.Sprintf("Author<%d %s %s>", u.Id, u.Name, u.Country.Name)
}
