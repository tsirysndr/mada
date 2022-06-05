package mada

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/highlight/highlighter/ansi"
	olc "github.com/google/open-location-code/go"
	svc "github.com/tsirysndr/mada/interfaces"
	"github.com/tsirysndr/mada/types"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
)

type SearchService struct {
	db    *sql.DB
	index bleve.Index
}

func NewSearchService(db *sql.DB, index bleve.Index) svc.SearchSvc {
	return &SearchService{db: db, index: index}
}

func (s *SearchService) Search(term string, opt types.SearchOptions) (*types.SearchResult, error) {
	result, err := s.SearchPoint(term, opt)
	if err == nil {
		return result, err
	}

	query := bleve.NewQueryStringQuery(term)
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	search.Size = 100
	search.Highlight = bleve.NewHighlightWithStyle(ansi.Name)

	searchResults, err := s.index.Search(search)
	return &types.SearchResult{Result: searchResults}, err
}

func InitializeBleve(db *sql.DB) (bleve.Index, error) {
	if _, err := os.Stat(DATABASE_PATH); os.IsNotExist(err) {
		return Init(db)
	}
	return bleve.Open(DATABASE_PATH)
}

func (s *SearchService) SearchPoint(term string, opt types.SearchOptions) (result *types.SearchResult, err error) {
	area, err := olc.Decode(term)

	if err != nil {
		return nil, err
	}

	if opt.SearchForFokontany {
		f, _ := searchInFokontany(s.db, area, opt)
		fmt.Println(f)
		return &types.SearchResult{Fokontany: f}, nil
	}
	if opt.SearchForCommune {
		c, _ := searchInCommune(s.db, area, opt)
		return &types.SearchResult{Commune: c}, nil
	}
	if opt.SearchForDistrict {
		d, _ := searchInDistrict(s.db, area, opt)
		return &types.SearchResult{District: d}, nil
	}
	if opt.SearchForRegion {
		r, _ := searchInRegion(s.db, area, opt)
		return &types.SearchResult{Region: r}, nil
	}

	f, _ := searchInFokontany(s.db, area, opt)

	return &types.SearchResult{Fokontany: f}, nil
}

func searchInFokontany(db *sql.DB, area olc.CodeArea, opt types.SearchOptions) (*types.Fokontany, bool) {
	point := fmt.Sprintf("POINT(%f %f)", area.LngLo, area.LatLo)
	rows, err := db.Query("SELECT uid, name, commune, region, district, country, ST_AsText(geom) from fokontany f where st_contains(f.geom, st_geometryfromtext($1, 4326))", point)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var uid, name, commune, district, region, country, g string
	var noresults bool = true
	for rows.Next() {
		noresults = false
		rows.Scan(&uid, &name, &commune, &region, &district, &country, &g)

		p, _ := wkt.Unmarshal(g)

		return &types.Fokontany{
			ID:          uid,
			Name:        name,
			Commune:     commune,
			Region:      region,
			District:    district,
			Country:     country,
			Coordinates: p.(*geom.Polygon).Coords(),
			Point:       geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{area.LngLo, area.LatLo}).Coords(),
		}, noresults
	}
	return nil, noresults
}

func searchInCommune(db *sql.DB, area olc.CodeArea, opt types.SearchOptions) (*types.Commune, bool) {
	point := fmt.Sprintf("POINT(%f %f)", area.LngLo, area.LatLo)
	rows, err := db.Query("SELECT uid, name, region, district, country, ST_AsText(geom) from commune f where st_contains(f.geom, st_geometryfromtext($1, 4326))", point)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var uid, name, district, region, country, g string
	var noresults bool = true
	for rows.Next() {
		noresults = false
		rows.Scan(&uid, &name, &region, &district, &country, &g)

		p, _ := wkt.Unmarshal(g)

		return &types.Commune{
			ID:          uid,
			Name:        name,
			Region:      region,
			District:    district,
			Country:     country,
			Coordinates: p.(*geom.Polygon).Coords(),
			Point:       geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{area.LngLo, area.LatLo}).Coords(),
		}, noresults
	}
	return nil, noresults
}

func searchInDistrict(db *sql.DB, area olc.CodeArea, opt types.SearchOptions) (*types.District, bool) {
	point := fmt.Sprintf("POINT(%f %f)", area.LngLo, area.LatLo)
	rows, err := db.Query("SELECT uid, name, region, country, ST_AsText(geom) from district f where st_contains(f.geom, st_geometryfromtext($1, 4326))", point)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var uid, name, region, country, g string
	var noresults bool = true
	for rows.Next() {
		noresults = false
		rows.Scan(&uid, &name, &region, &country, &g)

		p, _ := wkt.Unmarshal(g)

		return &types.District{
			ID:          uid,
			Name:        name,
			Region:      region,
			Country:     country,
			Coordinates: p.(*geom.Polygon).Coords(),
			Point:       geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{area.LngLo, area.LatLo}).Coords(),
		}, noresults

	}
	return nil, noresults
}

func searchInRegion(db *sql.DB, area olc.CodeArea, opt types.SearchOptions) (*types.Region, bool) {
	point := fmt.Sprintf("POINT(%f %f)", area.LngLo, area.LatLo)
	rows, err := db.Query("SELECT uid, name, country, ST_AsText(geom) from region f where st_contains(f.geom, st_geometryfromtext($1, 4326))", point)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var uid, name, country, g string
	var noresults bool = true
	for rows.Next() {
		noresults = false
		rows.Scan(&uid, &name, &country, &g)

		p, _ := wkt.Unmarshal(g)

		return &types.Region{
			ID:          uid,
			Name:        name,
			Country:     country,
			Coordinates: p.(*geom.Polygon).Coords(),
			Point:       geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{area.LngLo, area.LatLo}).Coords(),
		}, noresults

	}
	return nil, noresults
}
