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
	BoardID   primitive.ObjectID `json:"board_id" bson:"board_id" validate:"required"`
}

// Board Model -- Set of tasks for a specific time period
type Board struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Title     string             `json:"title" bson:"title" validate:"min=4"`
	FromDate  time.Time          `json:"from_date" bson:"from_date" validate:"required" ` //validar que sea una fecha valida
	ToDate    time.Time          `json:"to_date" bson:"to_date" validate:"required"`      //validar que sea una fecha valida y posterior a la fecha de inicio
	Completed bool               `json:"completed" bson:"completed" default:"false"`
	OwnerID   primitive.ObjectID `json:"owner_id" bson:"owner_id"`
}

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
	Username     string             `json:"username" bson:"username" validate:"required"`
	Password     string             `json:"-" bson:"password" validate:"required"` //don't return this field in the response
	Email        string             `json:"email" bson:"email" validate:"required,email"`
	ResetCode    string             `json:"-" bson:"reset_code,omitempty"`
	ResetCodeExp time.Time          `json:"-" bson:"reset_code_exp,omitempty"`
}
