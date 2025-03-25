package responses

import (
	"app/internal/entities"
	"time"
)

type File struct {
	Name      string    `json:"name" bson:"name"`
	Path      string    `json:"path" bson:"path"`
	Size      int       `json:"size" bson:"size"`
	Type      string    `json:"type" bson:"type"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

func NewFile(file entities.File) File {
	return File{
		Name:      file.Name,
		Path:      file.Path,
		Size:      file.Size,
		Type:      file.Type,
		CreatedAt: file.CreatedAt,
	}
}
