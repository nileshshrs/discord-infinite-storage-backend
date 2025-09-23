package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
    ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
    Name      string              `bson:"name" json:"name"`
    Size      int64               `bson:"size" json:"size"`
    ChannelID string              `bson:"channel_id" json:"channel_id"`
    UserID    *primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"` // optional
    Chunks    []Chunk             `bson:"chunks" json:"chunks"`
    CreatedAt time.Time           `bson:"created_at" json:"created_at"`
}

type Chunk struct {
	Index     int    `bson:"index" json:"index"`
	MessageID string `bson:"message_id" json:"message_id"`
	URL       string `bson:"url" json:"url"`
	Filename  string `bson:"filename" json:"filename"`
	Size      int64  `bson:"size" json:"size"`
}
