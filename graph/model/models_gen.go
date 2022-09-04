// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Commune struct {
	ID          *string         `json:"id"`
	Name        *string         `json:"name"`
	Province    *string         `json:"province"`
	Code        *string         `json:"code"`
	District    *string         `json:"district"`
	Region      *string         `json:"region"`
	Country     *string         `json:"country"`
	Coordinates [][][][]float64 `json:"coordinates"`
	Point       []float64       `json:"point"`
	Geometry    *Geometry       `json:"geometry"`
}

type CommuneList struct {
	Data  []*Commune `json:"data"`
	After *Commune   `json:"after"`
}

type Country struct {
	ID          *string         `json:"id"`
	Name        *string         `json:"name"`
	Code        *string         `json:"code"`
	Coordinates [][][][]float64 `json:"coordinates"`
	Geometry    *MultiPolygon   `json:"geometry"`
}

type District struct {
	ID          *string         `json:"id"`
	Name        *string         `json:"name"`
	Province    *string         `json:"province"`
	Code        *string         `json:"code"`
	Region      *string         `json:"region"`
	Country     *string         `json:"country"`
	Coordinates [][][][]float64 `json:"coordinates"`
	Point       []float64       `json:"point"`
	Geometry    *Geometry       `json:"geometry"`
}

type DistrictList struct {
	Data  []*District `json:"data"`
	After *District   `json:"after"`
}

type Fields struct {
	Commune   *string `json:"commune"`
	Country   *string `json:"country"`
	District  *string `json:"district"`
	Fokontany *string `json:"fokontany"`
	Name      *string `json:"name"`
	Province  *string `json:"province"`
	Region    *string `json:"region"`
	Type      *string `json:"type"`
}

type Fokontany struct {
	ID          *string         `json:"id"`
	Name        *string         `json:"name"`
	Province    *string         `json:"province"`
	Code        *string         `json:"code"`
	Commune     *string         `json:"commune"`
	District    *string         `json:"district"`
	Region      *string         `json:"region"`
	Country     *string         `json:"country"`
	Coordinates [][][][]float64 `json:"coordinates"`
	Point       []float64       `json:"point"`
	Geometry    *Geometry       `json:"geometry"`
}

type FokontanyList struct {
	Data  []*Fokontany `json:"data"`
	After *Fokontany   `json:"after"`
}

type Geometry struct {
	Type         *string       `json:"type"`
	Polygon      *Polygon      `json:"polygon"`
	Multipolygon *MultiPolygon `json:"multipolygon"`
}

type Hit struct {
	ID     *string  `json:"id"`
	Score  *float64 `json:"score"`
	Fields *Fields  `json:"fields"`
}

type MultiPolygon struct {
	Type        *string          `json:"type"`
	Coordinates [][][][]*float64 `json:"coordinates"`
}

type Polygon struct {
	Type        *string        `json:"type"`
	Coordinates [][][]*float64 `json:"coordinates"`
}

type Region struct {
	ID          *string         `json:"id"`
	Name        *string         `json:"name"`
	Province    *string         `json:"province"`
	Code        *string         `json:"code"`
	Geometry    *Geometry       `json:"geometry"`
	Country     *string         `json:"country"`
	Coordinates [][][][]float64 `json:"coordinates"`
	Point       []float64       `json:"point"`
}

type RegionList struct {
	Data  []*Region `json:"data"`
	After *Region   `json:"after"`
}

type Results struct {
	Region    *Region    `json:"region"`
	District  *District  `json:"district"`
	Commune   *Commune   `json:"commune"`
	Fokontany *Fokontany `json:"fokontany"`
	Hits      []*Hit     `json:"hits"`
}

type Todo struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
	User *User  `json:"user"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}