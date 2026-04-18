package Repositories

import (
	"gorm.io/gorm"
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

func UpdateMessage(message Models.Message) (Models.Message, error) {
	tx := Database.DB.Model(&Models.Message{}).
		Where("uuid = ?", message.UUID).
		Updates(map[string]interface{}{
			"content":       message.Content,
			"severity":      message.Severity,
			"media_type_id": message.MediaTypeId,
		})

	if tx.Error != nil {
		return Models.Message{}, tx.Error
	}

	if tx.RowsAffected == 0 {
		return Models.Message{}, gorm.ErrRecordNotFound
	}

	err := Database.DB.Preload("MediaType").Where("uuid = ?", message.UUID).First(&message).Error
	if err != nil {
		return Models.Message{}, err
	}

	return message, nil
}
