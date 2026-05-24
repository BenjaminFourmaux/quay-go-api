package Services

import (
	"github.com/google/uuid"
	"quay-go-api/Common"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Entities/Models"
	"quay-go-api/Repositories"
	logger "quay-go-api/Services/Logger"
	"slices"
	"strings"
)

func ListMessages(filters map[string]string) ([]Dto.Message, error) {
	logger.Info("[Message Service] List Messages")
	logger.Debug("With filters: %#v", filters)

	// Validating filters
	logger.Debug("Validating filters")
	filterSeverities := []string{}
	if severityStr, ok := filters["severity"]; ok {
		// split severities (comma separated) and validate values
		for _, severity := range strings.Split(severityStr, ",") {
			if !Common.ValidateMessageSeverity(severity) {
				logger.Warning("Invalid severity filter value: %s", severity)
				return nil, Errors.InvalidParameterValue("severity", []string{"info", "warning", "error"})
			} else {
				filterSeverities = append(filterSeverities, severity)
			}
		}
	}

	logger.Info("Retrieving messages from database")
	messages, err := Repositories.GetAllMessages()
	if err != nil {
		logger.Error("Error retrieving messages from database: %s", err.Error())
		return []Dto.Message{}, err
	}

	logger.Debug("Messages find in database: %d", len(messages))

	messagesDto := []Dto.Message{}

	for _, message := range messages {
		if len(filterSeverities) == 0 || (len(filterSeverities) > 0 && slices.Contains(filterSeverities, message.Severity)) {
			messagesDto = append(messagesDto, Dto.Message{
				UUID:      message.UUID,
				Content:   message.Content,
				Severity:  message.Severity,
				MediaType: message.MediaType.Name,
			})
		}
	}

	logger.Debug("Messages after applying filters: %d", len(messagesDto))

	return messagesDto, nil
}

func CreateMessage(message Dto.CreateMessage) (Dto.Message, error) {
	logger.Info("[Message Service] Create Message")
	logger.Debug("With dto: %+v", message)

	// Check severity is valid
	if !Common.ValidateMessageSeverity(message.Severity) {
		logger.Warning("Invalid message severity: %s", message.Severity)
		return Dto.Message{}, Errors.MessageInvalidSeverity(message.Severity)
	}

	logger.Info("Creating message in database")
	messageToCreate := Models.Message{
		UUID:        uuid.New().String(),
		Content:     message.Content,
		Severity:    message.Severity,
		MediaTypeId: 3, // text/markdown
	}

	messageModel, err := Repositories.CreateMessage(messageToCreate)
	if err != nil {
		logger.Error("Error creating message in database: %s", err.Error())
		return Dto.Message{}, err
	}

	messageDto := Dto.Message{
		UUID:      messageModel.UUID,
		Content:   messageModel.Content,
		Severity:  messageModel.Severity,
		MediaType: messageModel.MediaType.Name,
	}

	logger.Success("Message created successfully: %s", messageDto.UUID)

	return messageDto, nil
}

func UpdateMessage(messageUUID string, message Dto.UpdateMessage) (Dto.Message, error) {
	logger.Info("[Message Service] Update Message")
	logger.Debug("Message UUID: %s", messageUUID)
	logger.Debug("With dto: %+v", message)

	// Check severity is valid
	if message.Severity != "" && !Common.ValidateMessageSeverity(message.Severity) {
		logger.Warning("Invalid message severity for update: %s", message.Severity)
		return Dto.Message{}, Errors.MessageInvalidSeverity(message.Severity)
	}

	logger.Info("Updating message in database")
	messageToUpdate := Models.Message{
		UUID:        messageUUID,
		Content:     message.Content,
		Severity:    message.Severity,
		MediaTypeId: 3, // text/markdown
	}

	messageModel, err := Repositories.UpdateMessage(messageToUpdate)
	if err != nil {
		logger.Error("Error updating message in database (uuid=%s): %s", messageUUID, err.Error())
		return Dto.Message{}, err
	}

	messageDto := Dto.Message{
		UUID:      messageModel.UUID,
		Content:   messageModel.Content,
		Severity:  messageModel.Severity,
		MediaType: messageModel.MediaType.Name,
	}

	logger.Success("Message updated successfully: %s", messageDto.UUID)

	return messageDto, nil
}

func DeleteMessage(messageUUID string) error {
	logger.Info("[Message Service] Delete Message")
	logger.Debug("Message UUID: %s", messageUUID)

	logger.Info("Deleting message in database")
	err := Repositories.DeleteMessage(messageUUID)
	if err != nil {
		logger.Error("Error deleting message in database (uuid=%s): %s", messageUUID, err.Error())
		return err
	}

	logger.Success("Message deleted successfully: %s", messageUUID)

	return nil
}
