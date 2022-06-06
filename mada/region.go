package mada

import (
	"database/sql"

	"github.com/blevesearch/bleve/v2"
	svc "github.com/tsirysndr/mada/interfaces"
	"github.com/tsirysndr/mada/types"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
)

type RegionService struct {
	db    *sql.DB
	index bleve.Index
}

func NewRegionService(db *sql.DB, index bleve.Index) svc.RegionSvc {
	return &RegionService{db: db, index: index}
}

func (r *RegionService) Count() (c int, err error) {
	row := r.db.QueryRow("SELECT count(*) FROM region")
	err = row.Scan(&c)
	return c, err
}

func (r *RegionService) List(skip, limit int) (*bleve.SearchResult, error) {
	query := bleve.NewMatchQuery("region")
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	search.From = skip
	search.Size = limit

	return r.index.Search(search)
}

func (r *RegionService) ShowRegion(id string) (*types.Region, error) {
	rows, _ := r.db.Query("SELECT uid, name, ST_AsText(geom) FROM region WHERE uid = $1", id)
	defer rows.Close()
	var uid, name, g string
	for rows.Next() {
		rows.Scan(&uid, &name, &g)
		p, _ := wkt.Unmarshal(g)
		return &types.Region{ID: uid, Name: name, Country: "Madagascar", Coordinates: p.(*geom.MultiPolygon).Coords()}, nil
	}
	return nil, nil
}
