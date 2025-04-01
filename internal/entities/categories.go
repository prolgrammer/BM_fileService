package entities

type Category struct {
	Id      string   `json:"id" bson:"_id,omitempty" mapstructure:"id"`
	Name    string   `json:"name" bson:"name"`
	Folders []Folder `json:"folders" bson:"folders"`
}

type FileCategory struct {
	CategoryId string   `json:"category_id" bson:"category_id"`
	Folders    []Folder `json:"folders" bson:"folders"`
}
