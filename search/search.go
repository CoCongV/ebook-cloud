package search

import (
	"ebook-cloud/config"

	"github.com/blevesearch/bleve"
)

//Index is bleve.Index
var (
	Index bleve.Index
	err   error
)

type BookIndex struct {
	Name string
}

//Setup is init bleve index
func Setup() {
	mapping := bleve.NewIndexMapping()
	Index, err = bleve.New(config.Conf.SearchIndexFile, mapping)
	if err == bleve.ErrorIndexPathExists {
		Index, _ = bleve.Open(config.Conf.SearchIndexFile)
	} else if err != nil {
		panic(err)
	}
}
