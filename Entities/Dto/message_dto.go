package Dto

type Message struct {
	UUID      string `json:"uuid"`
	Content   string `json:"content"`
	Severity  string `json:"severity"`
	MediaType string `json:"media_type"`
}
