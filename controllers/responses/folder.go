package responses

import "app/internal/entities"

type Folder struct {
	Name  string `json:"name" bson:"name"`
	Files []File `json:"files" bson:"files"`
}

func NewFolder(name string, files []entities.File) Folder {
	filesResponse := make([]File, len(files))
	for j, file := range files {
		filesResponse[j] = NewFile(file)
	}

	return Folder{
		Name:  name,
		Files: filesResponse,
	}

}
