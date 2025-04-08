package entities

import "time"

type File struct {
	Id          string         `json:"id" bson:"_id,omitempty"`
	Name        string         `json:"name" bson:"name"`
	Description string         `json:"description" bson:"description, omitempty"`
	Size        int            `json:"size" bson:"size"`
	Type        string         `json:"type" bson:"type"`
	Version     string         `json:"version" bson:"version"`
	CreatedAt   time.Time      `json:"created_at" bson:"created_at"`
	Categories  []FileCategory `json:"categories" bson:"categories"`
}
