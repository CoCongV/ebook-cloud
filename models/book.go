package models

import "fmt"

// Book is Model
type Book struct {
	Id     int64
	Name   string
	Author Author
}

func (b Book) String() string {
	return fmt.Sprintf("Book<%d %s>", b.Id, b.Name)
}
