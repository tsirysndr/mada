package interfaces

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/tsirysndr/mada/types"
)

type CommuneSvc interface {
	Count() (int, error)
	List(skip, limit int) (*bleve.SearchResult, error)
	ShowCommune(id string) (*types.Commune, error)
}

type DistrictSvc interface {
	Count() (int, error)
	List(skip, limit int) (*bleve.SearchResult, error)
	ShowDistrict(id string) (*types.District, error)
}

type FokontanySvc interface {
	Count() (int, error)
	List(skip, limit int) (*bleve.SearchResult, error)
	ShowFokontany(id string) (*types.Fokontany, error)
}

type RegionSvc interface {
	Count() (int, error)
	List(skip, limit int) (*bleve.SearchResult, error)
	ShowRegion(id string) (*types.Region, error)
}

type SearchSvc interface {
	Search(term string, opt types.SearchOptions) (*types.SearchResult, error)
}
