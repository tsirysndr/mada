package mada

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/pkg/browser"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
)

type Fokontany struct {
	Point       geom.Coord     `json:"point,omitempty"`
	ID          string         `json:"id,omitempty"`
	Name        string         `json:"name,omitempty"`
	Commune     string         `json:"commune,omitempty"`
	Region      string         `json:"region,omitempty"`
	District    string         `json:"district,omitempty"`
	Country     string         `json:"country,omitempty"`
	Coordinates [][]geom.Coord `json:"coordinates,omitempty"`
}

type FokontanyService struct {
	db *sql.DB
}

func NewFokontanyService() *FokontanyService {
	if os.Getenv("MADA_POSTGRES_URL") != "" {
		db, err := sql.Open("postgres", os.Getenv("MADA_POSTGRES_URL"))
		if err != nil {
			panic(err)
		}
		return &FokontanyService{db: db}
	}

	db, err := OpenSQLiteConnection()

	if err != nil {
		panic(err)
	}

	return &FokontanyService{db: db}
}

func (f *FokontanyService) List(outputInJSON bool, skip, limit int, openInBrowser bool) {
	index, err := InitializeBleve()
	if err != nil {
		panic(err)
	}
	query := bleve.NewMatchQuery("fokontany")
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

func (f *FokontanyService) ShowFokontany(id string, outputInJSON, openInBrowser bool) {
	rows, err := f.db.Query("SELECT uid, name, commune, region, district, country, ST_AsText(geom) FROM fokontany WHERE uid = $1", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var uid, name, commune, district, region, country, g string
	for rows.Next() {
		rows.Scan(&uid, &name, &commune, &region, &district, &country, &g)

		if openInBrowser {
			err := browser.OpenURL(fmt.Sprintf("http://localhost:%d", PORT))
			if err != nil {
				fmt.Printf("Open http://localhost:%d in your browser\n", PORT)
			}
			StartHttpServer()
		}

		p, _ := wkt.Unmarshal(g)

		if outputInJSON {
			b, _ := json.MarshalIndent(Fokontany{
				ID:          uid,
				Name:        name,
				Commune:     commune,
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
					commune
									%s
					district
									%s
					region
									%s
          country
									%s
					type
									fokontany
					geometry
									%v
		`, uid, name, commune, district, region, country, p.(*geom.Polygon).Coords())

	}
}
