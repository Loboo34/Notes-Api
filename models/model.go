package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type Note struct{
// 	ID int    `json:"id"`
// 	Title string  `json:"title"`
// 	Content string 	`json:"content"`
// 	Tags []string  `json:"tags"`
// 	CreatedAt time.Time  `json:"createdAt"`

// }

type Note struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	Tags      []string           `bson:"tags" json:"tags"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

//var notes []Note
