package entities

type Account struct {
	Id         string     `json:"id" bson:"_id,omitempty" mapstructure:"id"`
	Categories []Category `json:"categories"`
}
