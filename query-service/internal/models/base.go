package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Base contains common fields for all models
type Base struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created" json:"created"`
	UpdatedAt time.Time          `bson:"updated" json:"updated"`
}