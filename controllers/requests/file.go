package requests

import "mime/multipart"

type LoadFile struct {
	Folder `form:"folder"`
	Files  []*multipart.FileHeader `form:"files"`
}
