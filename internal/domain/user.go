package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type User struct {
	ID    bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name  string        `json:"name" bson:"name"`
	Phone string        `json:"phone,omitempty" bson:"phone,omitempty"`
}
