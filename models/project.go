package models

// Project represents metadata associated with a project name.
type Project struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Category string `json:"category"`
}
