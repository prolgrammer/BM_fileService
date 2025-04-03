package entities

type Folder struct {
	Name string `json:"name" bson:"name"`
}

func CreateFolder(name string) Folder {
	return Folder{
		Name: name,
	}
}
