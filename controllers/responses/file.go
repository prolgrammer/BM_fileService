package responses

import (
	"app/internal/entities"
	"time"
)

type File struct {
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Size        int       `json:"size" bson:"size"`
	Type        string    `json:"type" bson:"type"`
	Version     string    `json:"version" bson:"version"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}

func NewFile(file entities.File) File {
	return File{
		Name:        file.Name,
		Description: file.Description,
		Size:        file.Size,
		Type:        file.Type,
		Version:     file.Version,
		CreatedAt:   file.CreatedAt,
	}
}
