package types

import "github.com/twpayne/go-geom"

type Commune struct {
	ID          string           `json:"id,omitempty"`
	Name        string           `json:"name,omitempty"`
	Region      string           `json:"region,omitempty"`
	District    string           `json:"district,omitempty"`
	Country     string           `json:"country,omitempty"`
	Coordinates [][][]geom.Coord `json:"coordinates,omitempty"`
	Point       geom.Coord       `json:"point,omitempty"`
}
