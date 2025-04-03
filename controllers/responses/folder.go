package responses

type Folder struct {
	Name string `json:"name" bson:"name"`
}

func NewFolder(name string) Folder {
	return Folder{
		Name: name,
	}

}
