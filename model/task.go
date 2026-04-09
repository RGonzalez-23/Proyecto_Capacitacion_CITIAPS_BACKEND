package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title       string               `bson:"title" json:"title"`
	Description string               `bson:"description" json:"description"`
	Completed   bool                 `bson:"completed" json:"completed"`
	Tags        []primitive.ObjectID `bson:"tags" json:"tags"`
	TagNames    []string             `bson:"-" json:"tagNames"` // Campo poblado para respuestas
	CreatedAt   primitive.DateTime   `bson:"createdAt" json:"createdAt"`
}

// TaskRequest para recibir tags por nombre en POST/PUT requests
type TaskRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}
