package Dto

type UserLogin struct {
	Service           string `json:"service"`
	ServiceIdentifier string `json:"service_identifier"`
	Metadata          any    `json:"metadata"`
}
