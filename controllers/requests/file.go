package requests

import "mime/multipart"

type LoadFile struct {
	Dir   string                  `form:"dir"`
	Files []*multipart.FileHeader `form:"files"`
}
