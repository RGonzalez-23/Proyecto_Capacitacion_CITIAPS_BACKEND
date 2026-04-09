package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tag struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name      string             `bson:"name" json:"name"`
	Color     string             `bson:"color,omitempty" json:"color,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

// TagRequest para recibir tags por nombre en las requests
type TagRequest struct {
	Name  string `json:"name"`
	Color string `json:"color,omitempty"`
}
