package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateNoteRequestBody struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type UpdateNoteRequestBody struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}
