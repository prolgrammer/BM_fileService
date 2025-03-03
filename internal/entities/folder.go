package entities

type Folder struct {
	Id    string `json:"id" bson:"_id,omitempty" mapstructure:"id"`
	Name  string `json:"name" bson:"name"`
	Files []File `json:"files" bson:"files"`
}
