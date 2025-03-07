package model

// RecordID defines a record id. Together with RecordType
// identifies unique records across all types.
type RecordID string

// RecordType defines a record type. Together with RecordID
// identifies unique record across all types.
type RecordType string

// Existing record types.

const (
	RecordTypeMovie = RecordType("movie")
)

// UserID defines a user id.
type UserID string

// RatingValue defines a value of rating record.
type RatingValue int

// Rating defines a individual rating created by a user
// for some record.
type Rating struct {
	RecordID   RecordID    `json:"recordId"`
	RecordType RecordType  `json:"recordType"`
	UserID     UserID      `json:"userId"`
	Value      RatingValue `json:"value"`
}

// RatingEvent defines a event containing rating information.
type RatingEvent struct {
	UserID          UserID          `json:"userId"`
	RecordID        RecordID        `json:"recordId"`
	RecordType      RecordType      `json:"recordType"`
	Value           RatingValue     `json:"value"`
	RatingEventType RatingEventType `json:"eventType"`
}

// RatingEventType defines the type of a rating event.
type RatingEventType string

// Rating event types.
const (
	RatingEventTypePut    = "put"
	RatingEventTypeDelete = "delete"
)
