package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/phongld0308/movie-example/rating/internal/repository"
	"github.com/phongld0308/movie-example/rating/pkg/model"
)

// Repository defines a PostgreSQL-based rating repository.
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

// Get retrieves all ratings for a given record.
func (r *Repository) Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT user_id, value FROM ratings WHERE record_id = $1 AND record_type = $2",
		recordID, recordType,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query ratings: %v", err)
	}
	defer rows.Close()

	var ratings []model.Rating
	for rows.Next() {
		var userID string
		var value int32
		if err := rows.Scan(&userID, &value); err != nil {
			return nil, fmt.Errorf("failed to scan rating: %v", err)
		}

		ratings = append(ratings, model.Rating{
			UserID:     model.UserID(userID),
			RecordID:   recordID,
			RecordType: recordType,
			Value:      model.RatingValue(value),
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating ratings: %v", err)
	}

	if len(ratings) == 0 {
		return nil, repository.ErrNotFound
	}

	return ratings, nil
}

// Put adds a rating for a given record.
func (r *Repository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO ratings (record_id, record_type, user_id, value)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (record_id, record_type, user_id) DO UPDATE 
		 SET value = $4`,
		recordID, recordType, rating.UserID, rating.Value,
	)
	if err != nil {
		return fmt.Errorf("failed to insert rating: %v", err)
	}
	return nil
}

// Close closes the database connection.
func (r *Repository) Close() error {
	return r.db.Close()
}
