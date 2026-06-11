package Dto

type Repository struct {
	Namespace   string `json:"namespace"`
	Name        string `json:"name"` // Image name
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	Kind        string `json:"kind"` // Image or Application
	State       string `json:"state"`
	IsStarred   bool   `json:"is_starred"`
}
