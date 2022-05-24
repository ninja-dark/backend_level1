package model

type File struct {
	Name      string `json:"name"`
	Extension string `json:"ext"`
	Size      int64  `json:"size"`
}
