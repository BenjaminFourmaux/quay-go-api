package Dto

type Avatar struct {
	Name  string `json:"name"`
	Hash  string `json:"hash"`
	Color string `json:"color"` // hex color code
	Kind  string `json:"kind"`  // e.g., "user", "org"
}
