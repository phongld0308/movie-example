package gateway

import "errors"

// ErrNotFound is returned when the movie metadata is not found.
var ErrNotFound = errors.New("not found")
