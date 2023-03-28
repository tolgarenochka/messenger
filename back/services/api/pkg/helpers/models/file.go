package models

type File struct {
	Path string `json:"path" db:"path"`
	Name string `json:"name" db:"name"`
}
