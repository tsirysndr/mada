package mada

import (
	"database/sql"

	"github.com/blevesearch/bleve/v2"
	svc "github.com/tsirysndr/mada/interfaces"
	"github.com/tsirysndr/mada/types"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
)

type DistrictService struct {
	db    *sql.DB
	index bleve.Index
}

func NewDistrictService(db *sql.DB, index bleve.Index) svc.DistrictSvc {
	return &DistrictService{db: db, index: index}
}

func (d *DistrictService) List(skip, limit int) (*bleve.SearchResult, error) {
	query := bleve.NewMatchQuery("district")
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	search.From = skip
	search.Size = limit

	return d.index.Search(search)
}

func (d *DistrictService) Count() (c int, err error) {
	row := d.db.QueryRow("SELECT count(*) FROM district")
	err = row.Scan(&c)
	return c, err
}

func (d *DistrictService) ShowDistrict(id string) (*types.District, error) {
	rows, _ := d.db.Query("SELECT uid, name, region, ST_AsText(geom) FROM district WHERE uid = $1", id)
	defer rows.Close()
	var uid, name, region, g string
	for rows.Next() {
		rows.Scan(&uid, &name, &region, &g)

		p, _ := wkt.Unmarshal(g)

		return &types.District{
			ID:          uid,
			Name:        name,
			Region:      region,
			Country:     "Madagascar",
			Coordinates: p.(*geom.Polygon).Coords(),
		}, nil
	}
	return nil, nil
}
