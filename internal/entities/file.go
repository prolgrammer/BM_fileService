package entities

import "time"

type File struct {
	Id        string    `json:"id" bson:"_id,omitempty" mapstructure:"id"`
	Name      string    `json:"name" bson:"name"`
	Path      string    `json:"path" bson:"path"`
	Size      int       `json:"size" bson:"size"`
	Type      string    `json:"type" bson:"type"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
