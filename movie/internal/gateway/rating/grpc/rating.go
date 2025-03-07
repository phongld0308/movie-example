package grpc

import (
	"context"

	"github.com/phongld0308/movie-example/gen"
	"github.com/phongld0308/movie-example/internal/grpcutil"
	"github.com/phongld0308/movie-example/pkg/discovery"
	model "github.com/phongld0308/movie-example/rating/pkg/model"
)

// Gateway defines an gRPC gate for rating service.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new gRPC gateway for rating service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

// GetAggregatedRating returns the aggregated rating for a record or ErrNotFound if there are no ratings for it.
func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)

	if err != nil {
		return 0, err
	}

	defer conn.Close()

	client := gen.NewRatingServiceClient(conn)

	resp, err := client.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{RecordId: string(recordID), RecordType: string(recordType)})
	if err != nil {
		return 0, err
	}

	return resp.RatingValue, nil
}

// PutRating returns the aggregated rating for a record or ErrNotFound if there are no ratings for it.
func (g *Gateway) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)

	if err != nil {
		return err
	}

	defer conn.Close()

	client := gen.NewRatingServiceClient(conn)

	_, err = client.PutRating(ctx, &gen.PutRatingRequest{RecordId: string(recordID), RecordType: string(recordType), RatingValue: int32(rating.Value)})
	if err != nil {
		return err
	}

	return nil
}
