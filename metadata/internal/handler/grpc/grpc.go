package grpc

import (
	"context"
	"errors"

	"github.com/phongld0308/movie-example/gen"
	"github.com/phongld0308/movie-example/metadata/internal/controller/metadata"
	"github.com/phongld0308/movie-example/metadata/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a movie metadata gRPC handler.
type Handler struct {
	gen.UnimplementedMetadataServiceServer
	ctrl *metadata.Controller
}

// New creates a new movie metadata gRPC handler.
func New(ctrl *metadata.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// Getmetadata returns movie metadata.
func (h *Handler) Getmetadata(ctx context.Context, req *gen.GetMetadataRequest) (*gen.GetMetadataResponse, error) {
	if req == nil || req.MovieId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	m, err := h.ctrl.Get(ctx, req.MovieId)
	if err != nil && errors.Is(err, metadata.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.GetMetadataResponse{
		Metadata: model.MetadataToProto(m),
	}, nil
}

func (h *Handler) PutMetadata(ctx context.Context, req *gen.PutMetadataRequest) (*gen.PutMetadataResponse, error) {
	if req == nil || req.Metadata == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or metadata")
	}

	metadata := &model.Metadata{
		ID:          req.Metadata.Id,
		Title:       req.Metadata.Title,
		Description: req.Metadata.Description,
		Director:    req.Metadata.Director,
	}

	if err := h.ctrl.Put(ctx, metadata); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.PutMetadataResponse{}, nil
}
