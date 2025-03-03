package entities

type Category struct {
	Id      string   `json:"id" bson:"_id,omitempty" mapstructure:"id"`
	Name    string   `json:"name" bson:"name"`
	UserId  string   `json:"user_id" bson:"user_id"`
	Folders []Folder `json:"folders" bson:"folders"`
}
