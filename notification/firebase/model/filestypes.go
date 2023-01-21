package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Files struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FileName string             `bson:"filename"`
	Metadata map[string]string  `bson:"metadata"`
	Service  string             `bson:"service"`
	Private  bool               `bson:"private"`
	Url      string             `bson:"url"`
	UpdateAt time.Time          `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time          `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
