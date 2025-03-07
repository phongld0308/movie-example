package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/phongld0308/movie-example/movie/internal/repository"
	"github.com/phongld0308/movie-example/movie/pkg/model"
)

// Repository defines a PostgreSQL movie repository
type Repository struct {
	db *sql.DB
}

// New creates a new PostgreSQL movie repository
func New(host string, port int, user, password, dbname string) (*Repository, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &Repository{db: db}, nil
}

// Get retrieves movie details by ID
func (r *Repository) Get(ctx context.Context, id string) (*model.MovieDetails, error) {
	query := `
		SELECT m.id, m.title, m.description, m.director, 
			   COALESCE(AVG(r.value), 0) as avg_rating,
			   COUNT(r.value) as rating_count
		FROM movies m
		LEFT JOIN ratings r ON m.id = r.record_id AND r.record_type = 'movie'
		WHERE m.id = $1
		GROUP BY m.id, m.title, m.description, m.director`

	var movie model.MovieDetails
	var avgRating float64
	var ratingCount int

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&movie.Metadata.ID,
		&movie.Metadata.Title,
		&movie.Metadata.Description,
		&movie.Metadata.Director,
		&avgRating,
		&ratingCount,
	)

	if err == sql.ErrNoRows {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan movie: %v", err)
	}

	if ratingCount > 0 {
		movie.Rating = &avgRating
	}

	return &movie, nil
}

// Put stores new movie details
func (r *Repository) Put(ctx context.Context, movie *model.MovieDetails) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Insert movie metadata
	_, err = tx.ExecContext(ctx,
		`INSERT INTO movies (id, title, description, director)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (id) DO NOTHING`,
		movie.Metadata.ID,
		movie.Metadata.Title,
		movie.Metadata.Description,
		movie.Metadata.Director,
	)
	if err != nil {
		return fmt.Errorf("failed to insert movie: %v", err)
	}

	// If there's a rating, store it
	if movie.Rating != nil {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO ratings (record_id, record_type, value)
			 VALUES ($1, 'movie', $2)`,
			movie.Metadata.ID,
			*movie.Rating,
		)
		if err != nil {
			return fmt.Errorf("failed to insert rating: %v", err)
		}
	}

	return tx.Commit()
}

// Update updates existing movie details
func (r *Repository) Update(ctx context.Context, movie *model.MovieDetails) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Update movie metadata
	result, err := tx.ExecContext(ctx,
		`UPDATE movies 
		 SET title = $2, description = $3, director = $4
		 WHERE id = $1`,
		movie.Metadata.ID,
		movie.Metadata.Title,
		movie.Metadata.Description,
		movie.Metadata.Director,
	)
	if err != nil {
		return fmt.Errorf("failed to update movie: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	// If there's a rating, update it
	if movie.Rating != nil {
		_, err = tx.ExecContext(ctx,
			`UPDATE ratings 
			 SET value = $2
			 WHERE record_id = $1 AND record_type = 'movie'`,
			movie.Metadata.ID,
			*movie.Rating,
		)
		if err != nil {
			return fmt.Errorf("failed to update rating: %v", err)
		}
	}

	return tx.Commit()
}

// Delete removes a movie by ID
func (r *Repository) Delete(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Delete ratings first (due to foreign key constraints)
	_, err = tx.ExecContext(ctx,
		"DELETE FROM ratings WHERE record_id = $1 AND record_type = 'movie'",
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete ratings: %v", err)
	}

	// Delete movie
	result, err := tx.ExecContext(ctx,
		"DELETE FROM movies WHERE id = $1",
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete movie: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	return tx.Commit()
}

// List returns all movies with optional pagination
func (r *Repository) List(ctx context.Context, skip, take int) ([]model.MovieDetails, error) {
	query := `
		SELECT m.id, m.title, m.description, m.director, 
			   COALESCE(AVG(r.value), 0) as avg_rating,
			   COUNT(r.value) as rating_count
		FROM movies m
		LEFT JOIN ratings r ON m.id = r.record_id AND r.record_type = 'movie'
		GROUP BY m.id, m.title, m.description, m.director
		ORDER BY m.id
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, take, skip)
	if err != nil {
		return nil, fmt.Errorf("failed to query movies: %v", err)
	}
	defer rows.Close()

	var movies []model.MovieDetails
	for rows.Next() {
		var movie model.MovieDetails
		var avgRating float64
		var ratingCount int

		err := rows.Scan(
			&movie.Metadata.ID,
			&movie.Metadata.Title,
			&movie.Metadata.Description,
			&movie.Metadata.Director,
			&avgRating,
			&ratingCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movie: %v", err)
		}

		if ratingCount > 0 {
			movie.Rating = &avgRating
		}

		movies = append(movies, movie)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating movies: %v", err)
	}

	return movies, nil
}

// Close closes the database connection
func (r *Repository) Close() error {
	return r.db.Close()
}
