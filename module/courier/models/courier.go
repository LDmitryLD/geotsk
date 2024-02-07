package models

type Courier struct {
	Score    int   `json:"score"`
	Location Point `json:"location"`
}

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
