package Services

import (
	"quay-go-api/Entities/Dto"
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
