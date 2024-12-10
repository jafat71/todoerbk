package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskStatus string

const (
	TODO  TaskStatus = "TODO"
	DOING TaskStatus = "DOING"
	DONE  TaskStatus = "DONE"
)

func (s TaskStatus) IsValid() bool {
	return s == TODO || s == DOING || s == DONE
}

type TaskPriority string

const (
	LOW    TaskPriority = "LOW"
	MEDIUM TaskPriority = "MEDIUM"
	HIGH   TaskPriority = "HIGH"
)

func (p TaskPriority) IsValid() bool {
	return p == LOW || p == MEDIUM || p == HIGH
}

type Task struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Title     string             `json:"title" bson:"title" validate:"min=4"`
	Status    TaskStatus         `json:"status" bson:"status"`
	Priority  TaskPriority       `json:"priority" bson:"priority"`
}
