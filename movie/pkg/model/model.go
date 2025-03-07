package model

import model "github.com/phongld0308/movie-example/metadata/pkg/model"

// MovieDetails includes movie meta data it aggregated rating.
type MovieDetails struct {
	Rating   *float64       `json:"rating"`
	Metadata model.Metadata `json:"metadata"`
}
