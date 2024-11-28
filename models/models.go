package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Title     string             `json:"title" bson:"title" validate:"min=4"`
	Doing     bool               `json:"doing" bson:"doing"`
	Done      bool               `json:"done" bson:"done"`
}
