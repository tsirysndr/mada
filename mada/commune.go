package mada

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
)

type Commune struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Region      string         `json:"region"`
	District    string         `json:"district"`
	Country     string         `json:"country"`
	Coordinates [][]geom.Coord `json:"coordinates"`
}

type CommuneService struct {
	db *sql.DB
}

func NewCommuneService() *CommuneService {
	if os.Getenv("MADA_POSTGRES_URL") != "" {
		db, err := sql.Open("postgres", os.Getenv("MADA_POSTGRES_URL"))
		if err != nil {
			panic(err)
		}
		return &CommuneService{db: db}
	}

	db, err := OpenSQLiteConnection()

	if err != nil {
		panic(err)
	}

	return &CommuneService{db: db}
}

func (c *CommuneService) List(outputInJSON bool, limit int) {
	index, err := InitializeBleve()
	if err != nil {
		panic(err)
	}
	query := bleve.NewMatchQuery("commune")
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

func (c *CommuneService) ShowCommune(id string, outputInJSON bool) {
	rows, _ := c.db.Query("SELECT uid, name, region, district, country, ST_AsText(geom) FROM commune WHERE uid = $1", id)
	defer rows.Close()
	var uid, name, region, district, country, g string
	rows.Next()
	rows.Scan(&uid, &name, &region, &district, &country, &g)

	p, _ := wkt.Unmarshal(g)

	if outputInJSON {
		b, _ := json.MarshalIndent(Commune{
			ID:          uid,
			Name:        name,
			Region:      region,
			District:    district,
			Country:     country,
			Coordinates: p.(*geom.Polygon).Coords(),
		}, "", "  ")

		fmt.Println(string(b))
		return
	}

	fmt.Printf(`
        id
                %s
        name
                %s
        district
                %s
        
        region
                %s
        type
                commune
        country
                Madagascar
        geometry
                %v
	`, uid, name, district, region, p.(*geom.Polygon).Coords())
}
