package models

type Chunk struct {
	Project string `json:"project"`
	InUse   bool   `json:"in_use"`
	Data    string `json:"data"`
	Version int    `json:"version"`
}
