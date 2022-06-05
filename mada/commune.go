package mada

import (
	"database/sql"

	"github.com/blevesearch/bleve/v2"
	svc "github.com/tsirysndr/mada/interfaces"
	"github.com/tsirysndr/mada/types"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
)

type CommuneService struct {
	db    *sql.DB
	index bleve.Index
}

func NewCommuneService(db *sql.DB, index bleve.Index) svc.CommuneSvc {
	return &CommuneService{db: db, index: index}
}

func (c *CommuneService) Count() (count int, err error) {
	row := c.db.QueryRow("SELECT count(*) FROM commune")
	err = row.Scan(&count)
	return count, err
}

func (c *CommuneService) List(skip, limit int) (*bleve.SearchResult, error) {
	query := bleve.NewMatchQuery("commune")
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	search.From = skip
	search.Size = limit

	return c.index.Search(search)
}

func (c *CommuneService) ShowCommune(id string) (*types.Commune, error) {
	rows, _ := c.db.Query("SELECT uid, name, region, district, country, ST_AsText(geom) FROM commune WHERE uid = $1", id)
	defer rows.Close()
	var uid, name, region, district, country, g string
	for rows.Next() {
		rows.Scan(&uid, &name, &region, &district, &country, &g)
		p, _ := wkt.Unmarshal(g)

		return &types.Commune{
			ID:          uid,
			Name:        name,
			Region:      region,
			District:    district,
			Country:     country,
			Coordinates: p.(*geom.Polygon).Coords(),
		}, nil
	}
	return nil, nil
}
