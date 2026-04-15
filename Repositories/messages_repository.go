package Repositories

import (
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

func GetAllMessages() ([]Models.Messages, error) {
	var messages []Models.Messages
	err := Database.DB.Preload("MediaType").Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}
