package entities

type Category struct {
	Id      string   `json:"id" bson:"_id,omitempty" mapstructure:"id"`
	Name    string   `bson:"name"`
	UserId  string   `bson:"user_id"`
	Folders []Folder `bson:"folders"`
}
