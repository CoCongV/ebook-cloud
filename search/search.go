package search

import (
	"ebook-cloud/config"

	"github.com/blevesearch/bleve"
)

//BookIndex is bleve.BookIndex
var (
	BookIndex bleve.Index
	err       error
)

type BookIndexData struct {
	Name string
}

//Setup is init bleve index
func Setup() {
	mapping := bleve.NewIndexMapping()
	BookIndex, err = bleve.New(config.Conf.BookSearchIndexFile, mapping)
	if err == bleve.ErrorIndexPathExists {
		BookIndex, _ = bleve.Open(config.Conf.BookSearchIndexFile)
	} else if err != nil {
		panic(err)
	}
}
