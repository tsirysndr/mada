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

type District struct {
	ID          string         `json:"id,omitempty"`
	Name        string         `json:"name,omitempty"`
	Region      string         `json:"region,omitempty"`
	Country     string         `json:"country,omitempty"`
	Coordinates [][]geom.Coord `json:"coordinates,omitempty"`
	Point       geom.Coord     `json:"point,omitempty"`
}

type DistrictService struct {
	db *sql.DB
}

func NewDistrictService() *DistrictService {
	if os.Getenv("MADA_POSTGRES_URL") != "" {
		db, err := sql.Open("postgres", os.Getenv("MADA_POSTGRES_URL"))
		if err != nil {
			panic(err)
		}
		return &DistrictService{db: db}
	}

	db, err := OpenSQLiteConnection()

	if err != nil {
		panic(err)
	}

	return &DistrictService{db: db}
}

func (d *DistrictService) List(outputInJSON bool, skip, limit int, openInBrowser bool) {
	index, err := InitializeBleve()
	if err != nil {
		panic(err)
	}
	query := bleve.NewMatchQuery("district")
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
		err := browser.OpenURL(fmt.Sprintf("http://localhost:%d", PORT))
		if err != nil {
			fmt.Printf("Open http://localhost:%d in your browser\n", PORT)
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

func (d *DistrictService) ShowDistrict(id string, outputInJSON, openInBrowser bool) {
	rows, _ := d.db.Query("SELECT uid, name, region, ST_AsText(geom) FROM district WHERE uid = $1", id)
	defer rows.Close()
	var uid, name, region, g string
	for rows.Next() {
		rows.Scan(&uid, &name, &region, &g)

		if openInBrowser {
			err := browser.OpenURL(fmt.Sprintf("http://localhost:%d", PORT))
			if err != nil {
				fmt.Printf("Open http://localhost:%d in your browser\n", PORT)
			}
			StartHttpServer()
		}

		p, _ := wkt.Unmarshal(g)

		if outputInJSON {
			b, _ := json.MarshalIndent(District{
				ID:          uid,
				Name:        name,
				Region:      region,
				Country:     "Madagascar",
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
        
        region
                %s
        type
                district
        country
                Madagascar
        geometry
                %v
	`, uid, name, region, p.(*geom.Polygon).Coords())
	}
}
