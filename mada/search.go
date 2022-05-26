package mada

import (
	"encoding/json"
	"fmt"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/highlight/highlighter/ansi"
)

func Search(term string, outputInJSON bool) {
	index, err := CreateOrOpenBleve()

	if err != nil {
		panic(err)
	}

	query := bleve.NewQueryStringQuery(term)
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	search.Size = 100
	search.Highlight = bleve.NewHighlightWithStyle(ansi.Name)

	searchResults, err := index.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !outputInJSON {
		fmt.Println(searchResults)
		return
	}

	b, _ := json.MarshalIndent(searchResults.Hits, "", "  ")

	fmt.Println(string(b))
}
