package mada

import (
	"crypto/sha256"
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/everystreet/go-shapefile"
	_ "github.com/lib/pq"
	"github.com/mitchellh/go-homedir"
	"github.com/twpayne/go-geom"
)

var DATABASE_PATH string = filepath.Join(CreateConfigDir(), "mada.bleve")

type Shape struct {
	Type       string     `json:"type"`
	Bbox       []float32  `json:"bbox"`
	Properties Properties `json:"properties"`
	Geom       Geometry   `json:"geometry"`
}

type Properties struct {
	Adm0en    string `json:"ADM0_EN"`
	Adm0type  string `json:"ADM0_TYPE"`
	Adm1en    string `json:"ADM1_EN"`
	Adm1type  string `json:"ADM1_TYPE"`
	Adm2en    string `json:"ADM2_EN"`
	Adm2type  string `json:"ADM2_TYPE"`
	Adm3en    string `json:"ADM3_EN"`
	Adm3type  string `json:"ADM3_TYPE"`
	Adm4en    string `json:"ADM4_EN"`
	Adm4type  string `json:"ADM4_TYPE"`
	OldProvin string `json:"OLD_PROVIN"`
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
	Type        string         `json:"type"`
	Coordinates [][]geom.Coord `json:"coordinates"`
}

//go:embed shp/*
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
		parseShapefile(filename, index, db)
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

func parseShapefile(name string, index bleve.Index, db *sql.DB) {

	shp, err := Assets.Open(filepath.Join("shp", fmt.Sprintf("%s.shp", name)))
	if err != nil {
		log.Fatal(err)
	}
	dbf, err := Assets.Open(filepath.Join("shp", fmt.Sprintf("%s.dbf", name)))
	if err != nil {
		log.Fatal(err)
	}

	scanner := shapefile.NewScanner(shp, dbf)

	if err != nil {
		log.Fatal(err)
	}
	// Start the scanner
	scanner.Scan()

	// Call Record() to get each record in turn, until either the end of the file, or an error occurs
	for {
		record := scanner.Record()
		if record == nil {
			break
		}

		feature := record.GeoJSONFeature()

		jsonData, _ := json.Marshal(feature)

		h := sha256.New()
		h.Write(jsonData)
		id := fmt.Sprintf("%x", h.Sum(nil))

		var shape Shape

		json.Unmarshal(jsonData, &shape)

		polygon := Polygon{
			Geometry: shape.Geom,
		}

		if strings.Contains(name, "adm0") {
			polygon.Type = "country"
			polygon.Name = shape.Properties.Adm0en
			polygon.Country = shape.Properties.Adm0en
		}

		if strings.Contains(name, "adm1") {
			polygon.Type = "region"
			polygon.Name = shape.Properties.Adm1en
			polygon.Region = shape.Properties.Adm1en
			polygon.Country = shape.Properties.Adm0en
			polygon.Province = shape.Properties.OldProvin
		}
		if strings.Contains(name, "adm2") {
			polygon.Type = "district"
			polygon.Name = shape.Properties.Adm2en
			polygon.District = shape.Properties.Adm2en
			polygon.Region = shape.Properties.Adm1en
			polygon.Country = shape.Properties.Adm0en
			polygon.Province = shape.Properties.OldProvin
		}
		if strings.Contains(name, "adm3") {
			polygon.Type = "commune"
			polygon.Name = shape.Properties.Adm3en
			polygon.Commune = shape.Properties.Adm3en
			polygon.District = shape.Properties.Adm2en
			polygon.Region = shape.Properties.Adm1en
			polygon.Country = shape.Properties.Adm0en
			polygon.Province = shape.Properties.OldProvin
		}
		if strings.Contains(name, "adm4") {
			polygon.Type = "fokontany"
			polygon.Name = shape.Properties.Adm4en
			polygon.Fokontany = shape.Properties.Adm4en
			polygon.Commune = shape.Properties.Adm3en
			polygon.District = shape.Properties.Adm2en
			polygon.Region = shape.Properties.Adm1en
			polygon.Country = shape.Properties.Adm0en
			polygon.Province = shape.Properties.OldProvin
		}

		fmt.Printf("%s - %s - %s\n", id, polygon.Type, polygon.Name)

		saveToDatabase(db, index, id, polygon)

	}

	// Err() returns the first error encountered during calls to Record()
	scanner.Err()
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
		"SELECT AddGeometryColumn('country', 'geom', 4326, 'POLYGON', 2);",
		"SELECT AddGeometryColumn('region', 'geom', 4326, 'POLYGON', 2);",
		"SELECT AddGeometryColumn('district', 'geom', 4326, 'POLYGON', 2);",
		"SELECT AddGeometryColumn('commune', 'geom', 4326, 'POLYGON', 2);",
		"SELECT AddGeometryColumn('fokontany', 'geom', 4326, 'POLYGON', 2);",
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
		"SELECT AddGeometryColumn('public', 'country', 'geom', 4326, 'POLYGON', 2);",
		"SELECT AddGeometryColumn('public', 'region', 'geom', 4326, 'POLYGON', 2);",
		"SELECT AddGeometryColumn('public', 'district', 'geom', 4326, 'POLYGON', 2);",
		"SELECT AddGeometryColumn('public', 'commune', 'geom', 4326, 'POLYGON', 2);",
		"SELECT AddGeometryColumn('public', 'fokontany', 'geom', 4326, 'POLYGON', 2);",
	}
	for _, query := range queries {
		db.Exec(query)
	}

}

func coordinatesToText(polygon Polygon) string {
	coordstr := "("
	for _, coords := range polygon.Geometry.Coordinates[0] {
		coordstr += fmt.Sprintf("%f %f, ", coords[0], coords[1])
	}

	coordstr = strings.TrimSuffix(coordstr, ", ")

	coordstr += ")"

	return coordstr
}

func saveToDatabase(db *sql.DB, index bleve.Index, id string, polygon Polygon) {
	err := index.Index(id, polygon)

	if err != nil {
		log.Fatal(err)
	}

	coordsText := coordinatesToText(polygon)

	geom := fmt.Sprintf("GeomFromText('POLYGON(%s)', 4326)", coordsText)

	if os.Getenv("MADA_POSTGRES_URL") != "" {
		geom = fmt.Sprintf("ST_GeomFromText('POLYGON(%s)', 4326)", coordsText)
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
