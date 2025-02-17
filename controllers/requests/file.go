package requests

import "mime/multipart"

type LoadFile struct {
	Dir   string           `json:"dir"`
	Files []multipart.File `json:"files"`
}
