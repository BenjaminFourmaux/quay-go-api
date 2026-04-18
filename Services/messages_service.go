package Services

import (
	"github.com/google/uuid"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Entities/Models"
	"quay-go-api/Repositories"
)

func ListMessages() ([]Dto.Message, error) {
	messages, err := Repositories.GetAllMessages()
	if err != nil {
		return nil, err
	}

	messagesDto := []Dto.Message{}

	for _, message := range messages {
		messagesDto = append(messagesDto, Dto.Message{
			UUID:      message.UUID,
			Content:   message.Content,
			Severity:  message.Severity,
			MediaType: message.MediaType.Name,
		})
	}

	return messagesDto, nil
}

func CreateMessage(message Dto.CreateMessage) (Dto.Message, error) {
	// Check severity is valid
	if message.Severity != "info" && message.Severity != "warning" && message.Severity != "error" {
		return Dto.Message{}, Errors.MessageInvalidSeverity(message.Severity)
	}

	messageToCreate := Models.Message{
		UUID:        uuid.New().String(),
		Content:     message.Content,
		Severity:    message.Severity,
		MediaTypeId: 3, // text/markdown
	}

	messageModel, err := Repositories.CreateMessage(messageToCreate)
	if err != nil {
		return Dto.Message{}, err
	}

	messageDto := Dto.Message{
		UUID:      messageModel.UUID,
		Content:   messageModel.Content,
		Severity:  messageModel.Severity,
		MediaType: messageModel.MediaType.Name,
	}

	return messageDto, nil
}

func UpdateMessage(messageUUID string, message Dto.UpdateMessage) (Dto.Message, error) {
	// Check severity is valid
	if message.Severity != "" && message.Severity != "info" && message.Severity != "warning" && message.Severity != "error" {
		return Dto.Message{}, Errors.MessageInvalidSeverity(message.Severity)
	}

	messageToUpdate := Models.Message{
		UUID:        messageUUID,
		Content:     message.Content,
		Severity:    message.Severity,
		MediaTypeId: 3, // text/markdown
	}

	messageModel, err := Repositories.UpdateMessage(messageToUpdate)
	if err != nil {
		return Dto.Message{}, err
	}

	messageDto := Dto.Message{
		UUID:      messageModel.UUID,
		Content:   messageModel.Content,
		Severity:  messageModel.Severity,
		MediaType: messageModel.MediaType.Name,
	}

	return messageDto, nil
}

func DeleteMessage(messageUUID string) error {
	return Repositories.DeleteMessage(messageUUID)
}
