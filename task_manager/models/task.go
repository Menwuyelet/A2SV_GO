package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Task struct {
	ID          bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	DueDate     string        `json:"due_date"`
	Status      string        `json:"status"`
}
