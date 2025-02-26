package entities

import "time"

type File struct {
	Name      string
	Path      string
	Size      int
	Type      string
	CreatedAt time.Time
}
