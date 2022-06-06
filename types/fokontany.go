package types

import "github.com/twpayne/go-geom"

type Fokontany struct {
	Point       geom.Coord       `json:"point,omitempty"`
	ID          string           `json:"id,omitempty"`
	Name        string           `json:"name,omitempty"`
	Commune     string           `json:"commune,omitempty"`
	Region      string           `json:"region,omitempty"`
	District    string           `json:"district,omitempty"`
	Country     string           `json:"country,omitempty"`
	Coordinates [][][]geom.Coord `json:"coordinates,omitempty"`
}
