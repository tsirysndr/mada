package mada

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/blevesearch/bleve/v2"
)

type Commune struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Region   string `json:"region"`
	District string `json:"district"`
	Country  string `json:"country"`
}

type CommuneService struct {
	db *sql.DB
}

func NewCommuneService() *CommuneService {
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
	rows, _ := c.db.Query("SELECT uid, name, region, district, country FROM commune WHERE uid = ?", id)
	defer rows.Close()
	var uid, name, region, district, country string
	rows.Next()
	rows.Scan(&uid, &name, &region, &district, &country)

	if outputInJSON {
		b, _ := json.MarshalIndent(Commune{
			ID:       uid,
			Name:     name,
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
        district
                %s
        
        region
                %s
        type
                commune
        country
                Madagascar
	`, uid, name, district, region)
}
