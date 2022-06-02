package mada

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/pkg/browser"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
)

type Commune struct {
	ID          string         `json:"id,omitempty"`
	Name        string         `json:"name,omitempty"`
	Region      string         `json:"region,omitempty"`
	District    string         `json:"district,omitempty"`
	Country     string         `json:"country,omitempty"`
	Coordinates [][]geom.Coord `json:"coordinates,omitempty"`
	Point       geom.Coord     `json:"point,omitempty"`
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

func (c *CommuneService) List(outputInJSON bool, skip, limit int, openInBrowser bool) {
	index, err := InitializeBleve()
	if err != nil {
		panic(err)
	}
	query := bleve.NewMatchQuery("commune")
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	search.From = skip
	search.Size = limit

	searchResults, err := index.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}

	if openInBrowser {
		err := browser.OpenURL("http://localhost:8010")
		if err != nil {
			fmt.Println("Open http://localhost:8010 in your browser")
		}
		StartHttpServer()
	}

	if !outputInJSON {
		fmt.Println(searchResults)
		return
	}

	b, _ := json.MarshalIndent(searchResults.Hits, "", "  ")

	fmt.Println(string(b))
}

func (c *CommuneService) ShowCommune(id string, outputInJSON, openInBrowser bool) {
	rows, _ := c.db.Query("SELECT uid, name, region, district, country, ST_AsText(geom) FROM commune WHERE uid = $1", id)
	defer rows.Close()
	var uid, name, region, district, country, g string
	for rows.Next() {
		rows.Scan(&uid, &name, &region, &district, &country, &g)

		if openInBrowser {
			err := browser.OpenURL("http://localhost:8010")
			if err != nil {
				fmt.Println("Open http://localhost:8010 in your browser")
			}
			StartHttpServer()
		}

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
}
