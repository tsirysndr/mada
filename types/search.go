package types

import "github.com/blevesearch/bleve/v2"

type SearchOptions struct {
	OutputInJSON       bool
	SearchForFokontany bool
	SearchForCommune   bool
	SearchForDistrict  bool
	SearchForRegion    bool
	OpenInBrowser      bool
}

type SearchResult struct {
	Fokontany *Fokontany
	Commune   *Commune
	District  *District
	Region    *Region
	Result    *bleve.SearchResult
}
