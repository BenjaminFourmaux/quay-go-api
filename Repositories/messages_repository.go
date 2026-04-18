package Repositories

import (
	"quay-go-api/Database"
	"quay-go-api/Entities/Models"
)

func GetAllMessages() ([]Models.Message, error) {
	var messages []Models.Message
	err := Database.DB.Preload("MediaType").Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func CreateMessage(message Models.Message) (Models.Message, error) {
	err := Database.DB.Create(&message).Error
	if err != nil {
		return Models.Message{}, err
	}

	err = Database.DB.Preload("MediaType").First(&message, message.ID).Error
	if err != nil {
		return Models.Message{}, err
	}
	return message, nil
}
