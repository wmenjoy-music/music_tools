package search

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/document"
	"github.com/blevesearch/bleve/mapping"
	"os"
)


type Searcher struct {
	FileName string

	mapping  *mapping.IndexMappingImpl
	index    bleve.Index
}

func NewSearcher(path string) (*Searcher, error) {
	mapping := bleve.NewIndexMapping()
	var index bleve.Index
	var err error
	if _, errors := os.Stat(path); errors != nil {
		index, err = bleve.New(path, mapping)
		if err != nil {
			return nil , err
		}
	} else {
		index, err = bleve.Open(path)
		if err != nil {
			return nil , err
		}
	}


	return &Searcher{
		FileName: path,
		mapping: mapping,
		index: index,
	}, nil
}

func (s Searcher) Index(id string, message interface{}) error {
	return s.index.Index(id, message)
}

func (s Searcher) Search(message string) (*bleve.SearchResult, error) {
	query := bleve.NewQueryStringQuery(message)
	searchRequest := bleve.NewSearchRequest(query)
	return s.index.Search(searchRequest)
}

func (s Searcher) Document(id string) (*document.Document, error) {
	return s.index.Document(id)
}