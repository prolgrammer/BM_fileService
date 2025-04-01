package requests

type Folder struct {
	Category `json:"category" form:"category"`
	Name     string `json:"name" form:"folder[name]"`
}
