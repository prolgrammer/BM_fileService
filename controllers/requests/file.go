package requests

import "mime/multipart"

type File struct {
	Folder `json:"folder"`
	Name   string `json:"name"`
}

type CreateFile struct {
	Folder `form:"folder"`
	Files  []*multipart.FileHeader `form:"files"`
}
