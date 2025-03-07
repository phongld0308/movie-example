package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/phongld0308/movie-example/metadata/internal/repository"
	"github.com/phongld0308/movie-example/metadata/pkg/model"
)

// Repository defines a PostgreSQL-based movie metadata repository.
type Repository struct {
	db *sql.DB
}

// New creates a new PostgreSQL-based repository.
func New(host string, port int, user, password, dbname string) (*Repository, error) {
	// Build PostgreSQL connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &Repository{db}, nil
}

// Get retrieves movie metadata by movie id.
func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	var title, description, director string

	row := r.db.QueryRowContext(ctx,
		"SELECT title, description, director FROM movies WHERE id = $1",
		id,
	)

	if err := row.Scan(&title, &description, &director); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan movie data: %v", err)
	}

	return &model.Metadata{
		ID:          id,
		Title:       title,
		Description: description,
		Director:    director,
	}, nil
}

// Put adds movie metadata for a given movie id.
func (r *Repository) Put(ctx context.Context, id string, metadata *model.Metadata) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO movies (id, title, description, director) 
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (id) DO UPDATE 
		 SET title = $2, description = $3, director = $4`,
		id, metadata.Title, metadata.Description, metadata.Director,
	)
	if err != nil {
		return fmt.Errorf("failed to insert movie: %v", err)
	}
	return nil
}

// Close closes the database connection.
func (r *Repository) Close() error {
	return r.db.Close()
}
