package responses

import "app/internal/entities"

type Category struct {
	Name    string            `json:"name"`
	UserId  string            `json:"user_id"`
	Folders []entities.Folder `json:"folders"`
}
