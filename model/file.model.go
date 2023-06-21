package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID              *primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	CreatedTime     *time.Time          `json:"createdTime,omitempty" bson:"created_time,omitempty"`
	LastUpdatedTime *time.Time          `json:"lastUpdatedTime,omitempty" bson:"last_updated_time,omitempty"`

	UploadID string `json:"uploadId,omitempty" bson:"upload_id,omitempty"`
	Key      string `json:"key,omitempty" bson:"key,omitempty"`
	Url      string `json:"url,omitempty" bson:"url,omitempty"`
}
