package types

import "github.com/twpayne/go-geom"

type Region struct {
	ID          string           `json:"id,omitempty"`
	Name        string           `json:"name,omitempty"`
	Country     string           `json:"country,omitempty"`
	Coordinates [][][]geom.Coord `json:"coordinates,omitempty"`
	Point       geom.Coord       `json:"point,omitempty"`
}
