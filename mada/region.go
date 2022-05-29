package mada

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/blevesearch/bleve/v2"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
)

type Region struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Country     string         `json:"country"`
	Coordinates [][]geom.Coord `json:"coordinates"`
}

type RegionService struct {
	db *sql.DB
}

func NewRegionService() *RegionService {
	db, err := OpenSQLiteConnection()

	if err != nil {
		panic(err)
	}

	return &RegionService{db: db}
}

func (r *RegionService) List(outputInJSON bool, limit int) {
	index, err := InitializeBleve()
	if err != nil {
		panic(err)
	}
	query := bleve.NewMatchQuery("region")
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	search.Size = limit

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

func (r *RegionService) ShowRegion(id string, outputInJSON bool) {
	rows, _ := r.db.Query("SELECT uid, name, ST_AsText(geom) FROM region WHERE uid = ?", id)
	defer rows.Close()
	var uid, name, g string
	rows.Next()
	rows.Scan(&uid, &name, &g)

	p, _ := wkt.Unmarshal(g)

	if outputInJSON {
		b, _ := json.MarshalIndent(Region{ID: uid, Name: name, Country: "Madagascar", Coordinates: p.(*geom.Polygon).Coords()}, "", "  ")
		fmt.Println(string(b))
		return
	}

	fmt.Printf(`
        id
                %s
        name
                %s
        type
                region
        country
                Madagascar
        geometry
               %v
	`, uid, name, p.(*geom.Polygon).Coords())
}
