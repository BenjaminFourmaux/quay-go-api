package Services

import (
	"quay-go-api/Common"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Repositories"
	"quay-go-api/Services/Auth"
	logger "quay-go-api/Services/Logger"
)

func ListUserRobots(filters map[string]string, currentUser *Auth.AuthenticatedUser) ([]Dto.Robot, error) {
	logger.Info("[Robot Service] List User Robots")
	logger.Debug("With filters: %+v", filters)

	// Validating filters
	var includeToken bool = false
	var includeRepositories bool = false
	if t, ok := filters["token"]; ok && t != "" {
		includeToken = t == "true"
	}
	if r, ok := filters["repositories"]; ok {
		includeRepositories = r == "true"
	}

	// Get user associated robot accounts from db
	robotUserModels, err := Repositories.GetUserRobots(currentUser.ID, includeToken, includeRepositories)
	if err != nil {
		logger.Error("Error retrieving user robots: %v", err)
		return nil, err
	}

	// Convert models to DTOs
	var robotDTOs []Dto.Robot
	for _, robotModel := range robotUserModels {
		robotDTO := Dto.Robot{
			Name:         robotModel.Username,
			Description:  robotModel.RobotAccountMetadata.Description,
			Created:      robotModel.CreationDate.Time,
			LastAccessed: Common.ConvertSQLNullTimeToTime(robotModel.LastAccessed),
		}

		// Include token
		if includeToken && robotModel.RobotAccountToken != nil {
			token, decryptErr := Common.DecryptAESCipherToken(robotModel.RobotAccountToken.Token)
			if decryptErr != nil {
				logger.Error("Error decrypting robot token for robot '%s': %v", robotModel.Username, decryptErr)
				return nil, decryptErr
			}

			robotDTO.Token = &token
		}

		// Include repository names
		if includeRepositories && robotModel.RepositoryPermissions != nil {
			var repoNames []string
			for _, repoPerm := range robotModel.RepositoryPermissions {
				repoNames = append(repoNames, repoPerm.Repository.Name)
			}
			robotDTO.Repositories = &repoNames
		}

		robotDTOs = append(robotDTOs, robotDTO)
	}

	return robotDTOs, nil
}
