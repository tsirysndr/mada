package mada

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/highlight/highlighter/ansi"
	olc "github.com/google/open-location-code/go"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
)

type SearchOptions struct {
	OutputInJSON       bool
	SearchForFokontany bool
	SearchForCommune   bool
	SearchForDistrict  bool
	SearchForRegion    bool
}

func Search(term string, opt SearchOptions) {
	index, err := InitializeBleve()

	if err != nil {
		panic(err)
	}

	err = SearchPoint(term, opt)
	if err == nil {
		return
	}

	query := bleve.NewQueryStringQuery(term)
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	search.Size = 100
	search.Highlight = bleve.NewHighlightWithStyle(ansi.Name)

	searchResults, err := index.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !opt.OutputInJSON {
		fmt.Println(searchResults)
		return
	}

	b, _ := json.MarshalIndent(searchResults.Hits, "", "  ")

	fmt.Println(string(b))
}

func InitializeBleve() (bleve.Index, error) {
	if _, err := os.Stat(DATABASE_PATH); os.IsNotExist(err) {
		return Init()
	}
	return bleve.Open(DATABASE_PATH)
}

func SearchPoint(term string, opt SearchOptions) (err error) {
	var db *sql.DB
	if os.Getenv("MADA_POSTGRES_URL") != "" {
		db, err = OpenPostgresConnection()
	} else {
		db, err = OpenSQLiteConnection()
	}

	defer db.Close()

	area, err := olc.Decode(term)

	if err != nil {
		return err
	}

	if opt.SearchForFokontany {
		searchInFokontany(db, area, opt)
		return nil
	}
	if opt.SearchForCommune {
		searchInCommune(db, area, opt)
		return nil
	}
	if opt.SearchForDistrict {
		searchInDistrict(db, area, opt)
		return nil
	}
	if opt.SearchForRegion {
		searchInRegion(db, area, opt)
		return nil
	}

	searchInFokontany(db, area, opt)

	return nil
}

func searchInFokontany(db *sql.DB, area olc.CodeArea, opt SearchOptions) bool {
	point := fmt.Sprintf("POINT(%f %f)", area.LngLo, area.LatLo)
	rows, err := db.Query("SELECT uid, name, commune, region, district, country, ST_AsText(geom) from fokontany f where st_contains(f.geom, st_geometryfromtext($1, 4326))", point)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var uid, name, commune, district, region, country, g string
	var noresults bool = true
	for rows.Next() {
		noresults = false
		rows.Scan(&uid, &name, &commune, &region, &district, &country, &g)

		p, _ := wkt.Unmarshal(g)

		if opt.OutputInJSON {
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
			return noresults
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
	return noresults
}

func searchInCommune(db *sql.DB, area olc.CodeArea, opt SearchOptions) bool {
	point := fmt.Sprintf("POINT(%f %f)", area.LngLo, area.LatLo)
	rows, err := db.Query("SELECT uid, name, region, district, country, ST_AsText(geom) from commune f where st_contains(f.geom, st_geometryfromtext($1, 4326))", point)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var uid, name, district, region, country, g string
	var noresults bool = true
	for rows.Next() {
		noresults = false
		rows.Scan(&uid, &name, &region, &district, &country, &g)

		p, _ := wkt.Unmarshal(g)

		if opt.OutputInJSON {
			b, _ := json.MarshalIndent(Commune{
				ID:          uid,
				Name:        name,
				Region:      region,
				District:    district,
				Country:     country,
				Coordinates: p.(*geom.Polygon).Coords(),
			}, "", "  ")
			fmt.Println(string(b))
			return noresults
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
        country
                %s
        type
                commune
        geometry
                %v
	`, uid, name, district, region, country, p.(*geom.Polygon).Coords())

	}
	return noresults
}

func searchInDistrict(db *sql.DB, area olc.CodeArea, opt SearchOptions) bool {
	point := fmt.Sprintf("POINT(%f %f)", area.LngLo, area.LatLo)
	rows, err := db.Query("SELECT uid, name, region, country, ST_AsText(geom) from district f where st_contains(f.geom, st_geometryfromtext($1, 4326))", point)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var uid, name, region, country, g string
	var noresults bool = true
	for rows.Next() {
		noresults = false
		rows.Scan(&uid, &name, &region, &country, &g)

		p, _ := wkt.Unmarshal(g)

		if opt.OutputInJSON {
			b, _ := json.MarshalIndent(District{
				ID:          uid,
				Name:        name,
				Region:      region,
				Country:     country,
				Coordinates: p.(*geom.Polygon).Coords(),
			}, "", "  ")
			fmt.Println(string(b))
			return noresults
		}
		fmt.Printf(`
        id
                %s
        name
                %s
        region
                %s
        country
                %s
        type
                district
        geometry
                %v
	`, uid, name, region, country, p.(*geom.Polygon).Coords())

	}
	return noresults
}

func searchInRegion(db *sql.DB, area olc.CodeArea, opt SearchOptions) bool {
	point := fmt.Sprintf("POINT(%f %f)", area.LngLo, area.LatLo)
	rows, err := db.Query("SELECT uid, name, country, ST_AsText(geom) from region f where st_contains(f.geom, st_geometryfromtext($1, 4326))", point)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var uid, name, country, g string
	var noresults bool = true
	for rows.Next() {
		noresults = false
		rows.Scan(&uid, &name, &country, &g)

		p, _ := wkt.Unmarshal(g)

		if opt.OutputInJSON {
			b, _ := json.MarshalIndent(Region{
				ID:          uid,
				Name:        name,
				Country:     country,
				Coordinates: p.(*geom.Polygon).Coords(),
			}, "", "  ")
			fmt.Println(string(b))
			return noresults
		}
		fmt.Printf(`
        id
                %s
        name
                %s
        country
                %s
        type
                region
        geometry
                %v
	`, uid, name, country, p.(*geom.Polygon).Coords())

	}
	return noresults
}
