package repository

import (
	"context"
	"errors"

	"github.com/phongld0308/movie-example/movie/pkg/model"
)

// Common errors
var (
	ErrNotFound = errors.New("movie not found")
)

// Repository defines a movie repository interface
type Repository interface {
	// Get retrieves movie details by ID
	Get(ctx context.Context, id string) (*model.MovieDetails, error)

	// Put stores movie details
	Put(ctx context.Context, movie *model.MovieDetails) error

	// Update updates existing movie details
	Update(ctx context.Context, movie *model.MovieDetails) error

	// Delete removes a movie by ID
	Delete(ctx context.Context, id string) error

	// List returns all movies with optional pagination
	List(ctx context.Context, skip, take int) ([]model.MovieDetails, error)
}
