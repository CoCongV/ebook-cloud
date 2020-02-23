package search

import (
	"ebook-cloud/config"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
)

//BookIndex is bleve.BookIndex
var (
	BookIndex   bleve.Index
	AuthorIndex bleve.Index
	err         error
)

type IndexData struct {
	Name string
}

//Setup is init bleve index
func Setup() {
	mapping := bleve.NewIndexMapping()
	BookIndex = createIndex(config.Conf.BookSearchIndexFile, mapping)
	AuthorIndex = createIndex(config.Conf.AuthorSearchIndexFile, mapping)
}

func createIndex(path string, mapping *mapping.IndexMappingImpl) bleve.Index {
	index, err := bleve.New(path, mapping)
	if err == bleve.ErrorIndexPathExists {
		index, _ = bleve.Open(path)
	} else if err != nil {
		panic(err)
	}
	return index
}
