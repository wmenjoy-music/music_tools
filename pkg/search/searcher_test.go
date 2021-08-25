package search

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)


type T struct {
	Value string
}
func TestSimpleSearch(t *testing.T) {
	searcher, err := NewSearcher("./tmp/data")

	assert.Nil(t, err)

	message := struct{
		Id   string
		From string
		Body string
		Values []string
		T      T
	}{
		Id:   "1",
		From: "marty.schoch@gmail.com",
		Body: "bleve indexing is easy",
		Values: []string{"key", "key2"},
		T: T{
			Value: "232",
		},
	}

	err = searcher.Index("1", message)
	assert.Nil(t, err)

	result, err := searcher.Search("key")
	assert.Nil(t, err)
	fmt.Println(result)
	doc,err := searcher.Document(result.Hits[0].ID)

	assert.Equal(t, 6, len(doc.Fields))
	//assert.Equal(t, 1, result.Size())

}
