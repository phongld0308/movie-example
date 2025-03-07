package metadata

import (
	"context"
	"errors"

	"github.com/phongld0308/movie-example/metadata/internal/repository"
	model "github.com/phongld0308/movie-example/metadata/pkg/model"
)

// ErrNotFound is returned when request record is not found.
var ErrNotFound = errors.New("not found")

type metadataRepository interface {
	Get(ctx context.Context, id string) (*model.Metadata, error)
	Put(ctx context.Context, id string, metadata *model.Metadata) error
}

// Controller defines a metadata service controller.
type Controller struct {
	repo metadataRepository
}

// New creates a new metadata service controller.
func New(repo metadataRepository) *Controller {
	return &Controller{repo}
}

// Get returns movie metada by id.
func (c *Controller) Get(ctx context.Context, id string) (*model.Metadata, error) {
	res, err := c.repo.Get(ctx, id)

	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	return res, nil
}

// Put creates or updates movie metadata.
func (c *Controller) Put(ctx context.Context, metadata *model.Metadata) error {
	return c.repo.Put(ctx, metadata.ID, metadata)
}
