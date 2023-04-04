package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Username    string             `json:"username" bson:"username"`
	Title       string             `json:"title" `
	Text        string             `json:"text" `
	Answers     []Answers
	Tags        []string `json:"tags"`
	CreatedAt   time.Time
	AnswerCount int
	ViewCount   int
}

type Answers struct {
	Text      string `json:"text"`
	Username  string `json:"username"`
	CreatedAt time.Time
}
