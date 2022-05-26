package mada

import (
	"crypto/sha256"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/everystreet/go-shapefile"
	"github.com/twpayne/go-geom"
)

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

func Init() {

	index, err := CreateOrOpenBleve()

	if err != nil {
		log.Fatal(err)
	}

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
		parseShapefile(filename, index)
	}

}

func CreateOrOpenBleve() (bleve.Index, error) {
	if _, err := os.Stat("mada.bleve"); os.IsNotExist(err) {
		geometryMapping := bleve.NewDocumentDisabledMapping()
		mapping := bleve.NewIndexMapping()
		mapping.DefaultMapping.AddSubDocumentMapping("geometry", geometryMapping)
		return bleve.New("mada.bleve", mapping)
	}
	return bleve.Open("mada.bleve")
}

func parseShapefile(name string, index bleve.Index) {

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

		geom.NewPolygon(geom.XY).MustSetCoords(shape.Geom.Coordinates)

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

		err := index.Index(id, polygon)

		if err != nil {
			log.Fatal(err)
		}

	}

	// Err() returns the first error encountered during calls to Record()
	scanner.Err()
}
