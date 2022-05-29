package mada

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/blevesearch/bleve/v2"
)

type District struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Region  string `json:"region"`
	Country string `json:"country"`
}

type DistrictService struct {
	db *sql.DB
}

func NewDistrictService() *DistrictService {
	db, err := OpenSQLiteConnection()

	if err != nil {
		panic(err)
	}

	return &DistrictService{db: db}
}

func (d *DistrictService) List(outputInJSON bool, limit int) {
	index, err := InitializeBleve()
	if err != nil {
		panic(err)
	}
	query := bleve.NewMatchQuery("district")
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

func (d *DistrictService) ShowDistrict(id string, outputInJSON bool) {
	rows, _ := d.db.Query("SELECT uid, name, region FROM district WHERE uid = ?", id)
	defer rows.Close()
	var uid, name, region string
	rows.Next()
	rows.Scan(&uid, &name, &region)

	if outputInJSON {
		b, _ := json.MarshalIndent(District{ID: uid, Name: name, Region: region, Country: "Madagascar"}, "", "  ")
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
	`, uid, name, region)
}
