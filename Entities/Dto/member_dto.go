package Dto

type TeamMember struct {
	Name    string `json:"name"`
	Kind    string `json:"kind"`
	IsRobot bool   `json:"is_robot"`
	Avatar  Avatar `json:"avatar"`
	Invited bool   `json:"invited"`
}
