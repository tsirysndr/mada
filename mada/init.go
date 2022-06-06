package mada

import (
	"crypto/sha256"
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/blevesearch/bleve/v2"
	_ "github.com/lib/pq"
	"github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
)

var DATABASE_PATH string = filepath.Join(CreateConfigDir(), "mada.bleve")

type Feature struct {
	Geometry   Geometry               `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Properties struct {
	Adm0en    string `mapstructure:"ADM0_EN"`
	Adm0type  string `mapstructure:"ADM0_TYPE"`
	Adm1en    string `mapstructure:"ADM1_EN"`
	Adm1type  string `mapstructure:"ADM1_TYPE"`
	Adm2en    string `mapstructure:"ADM2_EN"`
	Adm2type  string `mapstructure:"ADM2_TYPE"`
	Adm3en    string `mapstructure:"ADM3_EN"`
	Adm3type  string `mapstructure:"ADM3_TYPE"`
	Adm4en    string `mapstructure:"ADM4_EN"`
	Adm4type  string `mapstructure:"ADM4_TYPE"`
	OldProvin string `mapstructure:"OLD_PROVIN"`
}

type Polygon struct {
	Type      string   `json:"type"`
	Name      string   `json:"name"`
	Country   string   `json:"country"`
	Region    string   `json:"region"`
	District  string   `json:"district"`
	Commune   string   `json:"commune"`
	Fokontany string   `json:"fokontany"`
	Geometry  Geometry `json:"geometry"`
	Province  string   `json:"province"`
}

type Geometry struct {
	Type        string           `json:"type"`
	Coordinates [][][]geom.Coord `json:"coordinates"`
}

//go:embed geojson/*
var Assets embed.FS

func Init(db *sql.DB) (bleve.Index, error) {

	index, err := CreateOrOpenBleve()

	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv("MADA_POSTGRES_URL") != "" {
		createPostgresTables(db)
		addPostgresGeometryColumns(db)
	} else {
		initializeSpatialMetadata(db)
		createTables(db)
		addGeometryColumns(db)
	}

	defer db.Close()

	/*
		level 0 - country
		level 1 - region
		level 2 - district
		level 3 - commune
		level 4 - fokontany
	*/
	var filenames = []string{
		"mdg_admbnda_adm0_BNGRC_OCHA_20181031",
		"mdg_admbnda_adm1_BNGRC_OCHA_20181031",
		"mdg_admbnda_adm2_BNGRC_OCHA_20181031",
		"mdg_admbnda_adm3_BNGRC_OCHA_20181031",
		"mdg_admbnda_adm4_BNGRC_OCHA_20181031",
	}
	for _, filename := range filenames {
		parseGeoJSONfile(filename, index, db)
	}

	return index, nil
}

func CreateOrOpenBleve() (bleve.Index, error) {
	if _, err := os.Stat(DATABASE_PATH); os.IsNotExist(err) {
		geometryMapping := bleve.NewDocumentDisabledMapping()

		mapping := bleve.NewIndexMapping()
		mapping.DefaultMapping.AddSubDocumentMapping("geometry", geometryMapping)
		return bleve.New(DATABASE_PATH, mapping)
	}
	return bleve.Open(DATABASE_PATH)
}

func parseGeoJSONfile(name string, index bleve.Index, db *sql.DB) {
	geojsonfile, err := Assets.Open(filepath.Join("geojson", fmt.Sprintf("%s.json", name)))

	if err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadAll(geojsonfile)

	if err != nil {
		log.Fatal(err)
	}

	var fc FeatureCollection
	err = json.Unmarshal(content, &fc)

	if err != nil {
		log.Fatal(err)
	}

	for _, feature := range fc.Features {
		var properties Properties
		mapstructure.Decode(feature.Properties, &properties)
		polygon := Polygon{
			Geometry: feature.Geometry,
		}
		if strings.Contains(name, "adm0") {
			polygon.Type = "country"
			polygon.Name = properties.Adm0en
			polygon.Country = properties.Adm0en
		}

		if strings.Contains(name, "adm1") {
			polygon.Type = "region"
			polygon.Name = properties.Adm1en
			polygon.Region = properties.Adm1en
			polygon.Country = properties.Adm0en
			polygon.Province = properties.OldProvin
		}
		if strings.Contains(name, "adm2") {
			polygon.Type = "district"
			polygon.Name = properties.Adm2en
			polygon.District = properties.Adm2en
			polygon.Region = properties.Adm1en
			polygon.Country = properties.Adm0en
			polygon.Province = properties.OldProvin
		}
		if strings.Contains(name, "adm3") {
			polygon.Type = "commune"
			polygon.Name = properties.Adm3en
			polygon.Commune = properties.Adm3en
			polygon.District = properties.Adm2en
			polygon.Region = properties.Adm1en
			polygon.Country = properties.Adm0en
			polygon.Province = properties.OldProvin
		}
		if strings.Contains(name, "adm4") {
			polygon.Type = "fokontany"
			polygon.Name = properties.Adm4en
			polygon.Fokontany = properties.Adm4en
			polygon.Commune = properties.Adm3en
			polygon.District = properties.Adm2en
			polygon.Region = properties.Adm1en
			polygon.Country = properties.Adm0en
			polygon.Province = properties.OldProvin
		}

		saveToDatabase(db, index, polygon)
	}
}

func CreateConfigDir() string {
	home, _ := homedir.Dir()
	path := filepath.Join(home, ".mada")
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	return path
}

func runQuery(db *sql.DB, query string) (sql.Result, error) {
	tx, err := db.Begin()

	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(query)

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	result, err := stmt.Exec()

	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()

	return result, nil
}

func initializeSpatialMetadata(db *sql.DB) {
	rows, err := db.Query("SELECT COUNT(*) as count FROM sqlite_master WHERE type='table' AND name='spatial_ref_sys';")

	if err != nil {
		log.Fatal(err)
	}

	count := 0
	for rows.Next() {
		rows.Scan(&count)
	}

	if count == 1 {
		return
	}

	runQuery(db, "SELECT InitSpatialMetaData();")
}

func createTables(db *sql.DB) {
	queries := []string{
		"CREATE TABLE IF NOT EXISTS country (id INTEGER PRIMARY KEY AUTOINCREMENT, uid TEXT NOT NULL UNIQUE, name TEXT);",
		"CREATE TABLE IF NOT EXISTS region (id INTEGER PRIMARY KEY AUTOINCREMENT, uid TEXT NOT NULL UNIQUE, name TEXT, country TEXT);",
		"CREATE TABLE IF NOT EXISTS district (id INTEGER PRIMARY KEY AUTOINCREMENT, uid TEXT NOT NULL UNIQUE, name TEXT, region TEXT, country TEXT);",
		"CREATE TABLE IF NOT EXISTS commune (id INTEGER PRIMARY KEY AUTOINCREMENT, uid TEXT NOT NULL UNIQUE, name TEXT, district TEXT, region TEXT, country TEXT);",
		"CREATE TABLE IF NOT EXISTS fokontany (id INTEGER PRIMARY KEY AUTOINCREMENT, uid TEXT NOT NULL UNIQUE, name TEXT, commune TEXT, district TEXT, region TEXT, country TEXT);",
		"CREATE UNIQUE INDEX IF NOT EXISTS country_uid_idx ON country (uid);",
		"CREATE UNIQUE INDEX IF NOT EXISTS region_uid_idx ON region (uid);",
		"CREATE UNIQUE INDEX IF NOT EXISTS district_uid_idx ON district (uid);",
		"CREATE UNIQUE INDEX IF NOT EXISTS commune_uid_idx ON commune (uid);",
		"CREATE UNIQUE INDEX IF NOT EXISTS fokontany_uid_idx ON fokontany (uid);",
	}
	for _, query := range queries {
		runQuery(db, query)
	}
}

func createPostgresTables(db *sql.DB) {
	queries := []string{
		"CREATE TABLE IF NOT EXISTS country (id SERIAL PRIMARY KEY, uid TEXT NOT NULL UNIQUE, name TEXT);",
		"CREATE TABLE IF NOT EXISTS region (id SERIAL PRIMARY KEY, uid TEXT NOT NULL UNIQUE, name TEXT, country TEXT);",
		"CREATE TABLE IF NOT EXISTS district (id SERIAL PRIMARY KEY, uid TEXT NOT NULL UNIQUE, name TEXT, region TEXT, country TEXT);",
		"CREATE TABLE IF NOT EXISTS commune (id SERIAL PRIMARY KEY, uid TEXT NOT NULL UNIQUE, name TEXT, district TEXT, region TEXT, country TEXT);",
		"CREATE TABLE IF NOT EXISTS fokontany (id SERIAL PRIMARY KEY, uid TEXT NOT NULL UNIQUE, name TEXT, commune TEXT, district TEXT, region TEXT, country TEXT);",
		"CREATE UNIQUE INDEX IF NOT EXISTS country_uid_idx ON country (uid);",
		"CREATE UNIQUE INDEX IF NOT EXISTS region_uid_idx ON region (uid);",
		"CREATE UNIQUE INDEX IF NOT EXISTS district_uid_idx ON district (uid);",
		"CREATE UNIQUE INDEX IF NOT EXISTS commune_uid_idx ON commune (uid);",
		"CREATE UNIQUE INDEX IF NOT EXISTS fokontany_uid_idx ON fokontany (uid);",
	}
	for _, query := range queries {
		runQuery(db, query)
	}
}

func addGeometryColumns(db *sql.DB) {
	queries := []string{
		"SELECT AddGeometryColumn('country', 'geom', 4326, 'MULTIPOLYGON', 2);",
		"SELECT AddGeometryColumn('region', 'geom', 4326, 'MULTIPOLYGON', 2);",
		"SELECT AddGeometryColumn('district', 'geom', 4326, 'MULTIPOLYGON', 2);",
		"SELECT AddGeometryColumn('commune', 'geom', 4326, 'MULTIPOLYGON', 2);",
		"SELECT AddGeometryColumn('fokontany', 'geom', 4326, 'MULTIPOLYGON', 2);",
		"SELECT CreateSpatialIndex('country', 'geom');",
		"SELECT CreateSpatialIndex('region', 'geom');",
		"SELECT CreateSpatialIndex('district', 'geom');",
		"SELECT CreateSpatialIndex('commune', 'geom');",
		"SELECT CreateSpatialIndex('fokontany', 'geom');",
	}
	for _, query := range queries {
		db.Exec(query)
	}

}

func addPostgresGeometryColumns(db *sql.DB) {
	queries := []string{
		"SELECT AddGeometryColumn('public', 'country', 'geom', 4326, 'MULTIPOLYGON', 2);",
		"SELECT AddGeometryColumn('public', 'region', 'geom', 4326, 'MULTIPOLYGON', 2);",
		"SELECT AddGeometryColumn('public', 'district', 'geom', 4326, 'MULTIPOLYGON', 2);",
		"SELECT AddGeometryColumn('public', 'commune', 'geom', 4326, 'MULTIPOLYGON', 2);",
		"SELECT AddGeometryColumn('public', 'fokontany', 'geom', 4326, 'MULTIPOLYGON', 2);",
	}
	for _, query := range queries {
		db.Exec(query)
	}

}

func saveToDatabase(db *sql.DB, index bleve.Index, polygon Polygon) {
	g := geom.NewMultiPolygon(geom.XY).MustSetCoords(polygon.Geometry.Coordinates).SetSRID(4326)
	coordsText, err := wkt.Marshal(g)

	if err != nil {
		log.Fatal(err)
	}

	h := sha256.New()
	h.Write([]byte(coordsText))
	id := fmt.Sprintf("%x", h.Sum(nil))

	fmt.Printf("%s - %s - %s\n", id, polygon.Type, polygon.Name)

	err = index.Index(id, polygon)

	if err != nil {
		log.Fatal(err)
	}

	geom := fmt.Sprintf("GeomFromText('%s', 4326)", coordsText)

	if os.Getenv("MADA_POSTGRES_URL") != "" {
		geom = fmt.Sprintf("ST_GeomFromText('%s', 4326)", coordsText)
	}

	q := ""

	if polygon.Type == "country" {
		q = fmt.Sprintf(
			"INSERT INTO country (uid, name, geom) VALUES ('%s', '%s', %s);",
			id,
			escapeSingleQuote(polygon.Name),
			geom,
		)
	}

	if polygon.Type == "region" {
		q = fmt.Sprintf(
			"INSERT INTO region (uid, name, geom, country) VALUES ('%s', '%s', %s, '%s');",
			id,
			escapeSingleQuote(polygon.Name),
			geom,
			escapeSingleQuote(polygon.Country),
		)
	}

	if polygon.Type == "district" {
		q = fmt.Sprintf(
			"INSERT INTO district (uid, name, geom, region, country) VALUES ('%s', '%s', %s, '%s', '%s');",
			id,
			escapeSingleQuote(polygon.Name),
			geom,
			escapeSingleQuote(polygon.Region),
			escapeSingleQuote(polygon.Country),
		)
	}

	if polygon.Type == "commune" {
		q = fmt.Sprintf(
			"INSERT INTO commune (uid, name, geom, district, region, country) VALUES ('%s', '%s', %s, '%s', '%s', '%s');",
			id,
			escapeSingleQuote(polygon.Name),
			geom,
			escapeSingleQuote(polygon.District),
			escapeSingleQuote(polygon.Region),
			escapeSingleQuote(polygon.Country),
		)
	}

	if polygon.Type == "fokontany" {
		q = fmt.Sprintf(
			"INSERT INTO fokontany (uid, name, geom, commune, district, region, country) VALUES ('%s', '%s', %s, '%s', '%s', '%s', '%s');",
			id,
			escapeSingleQuote(polygon.Name),
			geom,
			escapeSingleQuote(polygon.Commune),
			escapeSingleQuote(polygon.District),
			escapeSingleQuote(polygon.Region),
			escapeSingleQuote(polygon.Country),
		)
	}

	if q != "" {
		_, err := db.Exec(q)
		if err != nil {
			log.Println(err)
		}
	}

}

func escapeSingleQuote(str string) string {
	return strings.ReplaceAll(str, "'", "''")
}
