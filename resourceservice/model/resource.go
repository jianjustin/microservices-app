package model

type Resource struct {
	Id           string
	ParentId     string
	ResourceType int64
	Name         string
	Content      string
	StdDelete    int64
}
