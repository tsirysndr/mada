package mada

import (
	"database/sql"
	"log"

	"github.com/blevesearch/bleve/v2"
	svc "github.com/tsirysndr/mada/interfaces"
	"github.com/tsirysndr/mada/types"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
)

type FokontanyService struct {
	db    *sql.DB
	index bleve.Index
}

func NewFokontanyService(db *sql.DB, index bleve.Index) svc.FokontanySvc {
	return &FokontanyService{db: db, index: index}
}

func (f *FokontanyService) List(skip, limit int) (*bleve.SearchResult, error) {
	query := bleve.NewMatchQuery("fokontany")
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	search.From = skip
	search.Size = limit

	return f.index.Search(search)
}

func (f *FokontanyService) ShowFokontany(id string) (*types.Fokontany, error) {
	rows, err := f.db.Query("SELECT uid, name, commune, region, district, country, ST_AsText(geom) FROM fokontany WHERE uid = $1", id)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var uid, name, commune, district, region, country, g string
	for rows.Next() {
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
		}, nil
	}

	return nil, nil
}
