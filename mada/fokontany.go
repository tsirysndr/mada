package mada

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/blevesearch/bleve/v2"
)

type Fokontany struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Commune  string `json:"commune"`
	Region   string `json:"region"`
	District string `json:"district"`
	Country  string `json:"country"`
}

type FokontanyService struct {
	db *sql.DB
}

func NewFokontanyService() *FokontanyService {
	db, err := OpenSQLiteConnection()

	if err != nil {
		panic(err)
	}

	return &FokontanyService{db: db}
}

func (f *FokontanyService) List(outputInJSON bool, limit int) {
	index, err := InitializeBleve()
	if err != nil {
		panic(err)
	}
	query := bleve.NewMatchQuery("fokontany")
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

func (f *FokontanyService) ShowFokontany(id string, outputInJSON bool) {
	rows, _ := f.db.Query("SELECT uid, name, commune, region, district, country FROM fokontany WHERE uid = ?", id)
	defer rows.Close()
	var uid, name, commune, district, region, country string
	rows.Next()
	rows.Scan(&uid, &name, &commune, &region, &district, &country)

	if outputInJSON {
		b, _ := json.MarshalIndent(Fokontany{
			ID:       uid,
			Name:     name,
			Commune:  commune,
			Region:   region,
			District: district,
			Country:  country,
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
        country
                Madagascar
	`, uid, name, commune, region, district, country)
}
