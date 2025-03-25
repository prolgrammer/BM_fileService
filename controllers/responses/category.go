package responses

import "app/internal/entities"

type Category struct {
	Name    string   `json:"name"`
	UserId  string   `json:"user_id"`
	Folders []Folder `json:"folders"`
}

func NewCategory(Name, UserId string, folders []entities.Folder) Category {
	respFolders := make([]Folder, len(folders))
	for i, folder := range folders {
		respFolders[i] = NewFolder(folder.Name, folder.Files)
	}

	return Category{
		Name:    Name,
		UserId:  UserId,
		Folders: respFolders,
	}
}
