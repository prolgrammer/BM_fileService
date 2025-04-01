package requests

type Category struct {
	Name string `json:"name" form:"category[name]"`
}
