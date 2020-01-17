package models

import (
	"fmt"
)

type Country struct {
	Id   int64
	Name string
}

func (c Country) String() string {
	return fmt.Sprintf("Country<%d %s>", c.Id, c.Name)
}

type Author struct {
	Id      int64
	Name    string
	Country *Country
}

func (a Author) String() string {
	return fmt.Sprintf("Author<%d %s %s>", a.Id, a.Name, a.Country.Name)
}
