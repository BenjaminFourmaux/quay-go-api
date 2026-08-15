package Services

import (
	"database/sql"
	"encoding/json"
	"quay-go-api/Common"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Entities/Models"
	"quay-go-api/Repositories"
	"quay-go-api/Services/Auth"
	logger "quay-go-api/Services/Logger"
	"strconv"
	"time"
)

func ListUserRobots(filters map[string]string, currentUser *Auth.AuthenticatedUser) ([]Dto.Robot, error) {
	return listRobots(filters, "user", strconv.Itoa(currentUser.ID), currentUser)
}

func ListOrganizationRobots(filters map[string]string, orgName string, currentUser *Auth.AuthenticatedUser) ([]Dto.Robot, error) {
	return listRobots(filters, "organization", orgName, currentUser)
}

/*
listRobots wrapper
*/
func listRobots(filters map[string]string, kind string, kindIdOrName string, currentUser *Auth.AuthenticatedUser) ([]Dto.Robot, error) {
	logger.Info("[Robot Service] List %s Robots accounts", kind)
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

	// Check if the user or org exists
	var kindId int
	if kind == "user" {
		userExist, err := Repositories.GetUserById(Common.ParseStringToInt(kindIdOrName)) // The Current user is already authenticated, it will exist
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user found with id: %s", kindIdOrName)
				return []Dto.Robot{}, Errors.UserNotFoundById(Common.ParseStringToInt(kindIdOrName))
			default:
				logger.Error("Error retrieving repository  from database: %s", err.Error())
				return []Dto.Robot{}, err
			}
		}

		kindId = userExist.ID
	} else if kind == "organization" {
		orgExist, err := Repositories.GetOrganizationByName(kindIdOrName)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No organization found with name: %s", kindIdOrName)
				return []Dto.Robot{}, Errors.OrganizationNotFound(kindIdOrName)
			default:
				logger.Error("Error retrieving organization from database: %s", err.Error())
				return []Dto.Robot{}, err
			}
		}

		kindId = orgExist.ID
	}

	// Get user associated robot accounts from db
	robotUserModels, err := Repositories.GetUserOrOrgRobots(kindId, includeToken, includeRepositories)
	if err != nil {
		logger.Error("Error retrieving %s robots: %v", kind, err)
		return []Dto.Robot{}, err
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
				return []Dto.Robot{}, decryptErr
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

func CreateUserRobot(robotToCreate Dto.CreateRobot, currentUser *Auth.AuthenticatedUser) (Dto.Robot, error) {
	return createRobot(robotToCreate, "user", strconv.Itoa(currentUser.ID), currentUser)
}

func CreateOrganizationRobot(robotToCreate Dto.CreateRobot, orgName string, currentUser *Auth.AuthenticatedUser) (Dto.Robot, error) {
	return createRobot(robotToCreate, "organization", orgName, currentUser)
}

/*
createRobot wrapper
*/
func createRobot(robotToCreate Dto.CreateRobot, kind string, kindIdOrName string, currentUser *Auth.AuthenticatedUser) (Dto.Robot, error) {
	logger.Info("[Robot Service] Create %s Robot account", kind)
	logger.Debug("With data: %+v", robotToCreate)

	// Validate input
	if err := Common.ValidateCreateRobotAccount(robotToCreate); err != nil {
		logger.Error("Input validation error: %v", err)
		return Dto.Robot{}, err
	}

	// Check if the user or org exists
	var kindId int
	var kindName string // need to create the robot name in the format: <kindName>+<robotName>
	if kind == "user" {
		userExist, err := Repositories.GetUserById(Common.ParseStringToInt(kindIdOrName)) // The Current user is already authenticated, it will exist
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user found with id: %s", kindIdOrName)
				return Dto.Robot{}, Errors.UserNotFoundById(Common.ParseStringToInt(kindIdOrName))
			default:
				logger.Error("Error retrieving repository  from database: %s", err.Error())
				return Dto.Robot{}, err
			}
		}

		kindId = userExist.ID
		kindName = userExist.Username
	} else if kind == "organization" {
		orgExist, err := Repositories.GetOrganizationByName(kindIdOrName)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No organization found with name: %s", kindIdOrName)
				return Dto.Robot{}, Errors.OrganizationNotFound(kindIdOrName)
			default:
				logger.Error("Error retrieving organization from database: %s", err.Error())
				return Dto.Robot{}, err
			}
		}

		kindId = orgExist.ID
		kindName = orgExist.Username
	}

	// Check if a robot with the same name already exists for this user/org
	existingRobot, err := Repositories.GetRobotByName(Common.FormatRobotUsername(kindName, robotToCreate.Name), kindId)
	if err != nil {
		switch err.Error() {
		case "record not found":
			// pass
		default:
			logger.Error("Error checking existing robots: %s", err.Error())
			return Dto.Robot{}, err
		}
	} else {
		if existingRobot.ID != 0 {
			logger.Warning("A robot with the name %s already exists for %s", robotToCreate.Name, kindName)
			return Dto.Robot{}, Errors.RobotAlreadyExists(Common.FormatRobotUsername(kindName, robotToCreate.Name))
		}
	}

	// All validations passed
	logger.Info("All validations passed")

	// Create the models
	var robotToCreateModel = Models.User{
		UUID:                  Common.GenerateUUID(),
		Username:              Common.FormatRobotUsername(kindName, robotToCreate.Name), // Format the robot name as <kindName>+<robotName>
		Robot:                 true,                                                     // Yes, a user can be a real user, an org or a robot account
		PasswordHash:          sql.NullString{},                                         // Robot accounts do not have a password
		Email:                 Common.GenerateUUID(),                                    // Robot email is an uuid but idk where pointed, so I consider to be a random uuid
		RemovedTagExpirationS: Common.DefaultTagExpirationSeconds,
		Enabled:               true,
		CreationDate:          sql.NullTime{Time: time.Now(), Valid: true},
	}

	var robotToCreateMetadataModel = Models.RobotAccountMetadata{
		Description:      robotToCreate.Description,
		UnstructuredJson: metadataMapToJSON(robotToCreate.UnstructuredMetadata),
	}

	// Generate a random token for the robot account and encrypt it
	_, encryptedToken, encryptErr := generateEncryptedToken()
	if encryptErr != nil {
		return Dto.Robot{}, encryptErr
	}

	var robotToCreateTokenModel = Models.RobotAccountToken{
		Token:         encryptedToken,
		FullyMigrated: true,
	}

	var robotFederatedLoginModel = Models.FederatedLogin{
		ServiceId:    Common.MapLoginServiceId("quayrobot"),
		MetadataJson: "{}",
	}

	// create the robot account in the db
	robotCreatedModel, err := Repositories.CreateRobotAccount(robotToCreateModel, robotToCreateMetadataModel, robotToCreateTokenModel, robotFederatedLoginModel)
	if err != nil {
		logger.Error("Error when creating robot account: %s", err.Error())
		return Dto.Robot{}, err
	}

	// Convert the created model to DTO
	robotDTO := Dto.Robot{
		Name:         robotCreatedModel.Username,
		Description:  robotCreatedModel.RobotAccountMetadata.Description,
		Created:      robotCreatedModel.CreationDate.Time,
		LastAccessed: Common.ConvertSQLNullTimeToTime(robotCreatedModel.LastAccessed),
	}

	return robotDTO, nil
}

func GetUserRobot(robotShortname string, currentUser *Auth.AuthenticatedUser) (Dto.Robot, error) {
	return getRobot(robotShortname, "user", strconv.Itoa(currentUser.ID), currentUser)
}

func GetOrganizationRobot(orgName string, robotShortname string, currentUser *Auth.AuthenticatedUser) (Dto.Robot, error) {
	return getRobot(robotShortname, "organization", orgName, currentUser)
}

func getRobot(robotShortname string, kind string, kindIdOrName string, currentUser *Auth.AuthenticatedUser) (Dto.Robot, error) {
	logger.Info("[Robot Service] Get %s Robot account", kind)
	logger.Debug("With robot shortname: %s", robotShortname)

	// Check if the user or org exists
	var kindId int
	var kindName string // need to create the robot name in the format: <kindName>+<robotName>
	if kind == "user" {
		userExist, err := Repositories.GetUserById(Common.ParseStringToInt(kindIdOrName)) // The Current user is already authenticated, it will exist
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user found with id: %s", kindIdOrName)
				return Dto.Robot{}, Errors.UserNotFoundById(Common.ParseStringToInt(kindIdOrName))
			default:
				logger.Error("Error retrieving repository  from database: %s", err.Error())
				return Dto.Robot{}, err
			}
		}

		kindId = userExist.ID
		kindName = userExist.Username
	} else if kind == "organization" {
		orgExist, err := Repositories.GetOrganizationByName(kindIdOrName)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No organization found with name: %s", kindIdOrName)
				return Dto.Robot{}, Errors.OrganizationNotFound(kindIdOrName)
			default:
				logger.Error("Error retrieving organization from database: %s", err.Error())
				return Dto.Robot{}, err
			}
		}

		kindId = orgExist.ID
		kindName = orgExist.Username
	}

	// Check if a robot with the same name already exists for this user/org
	existingRobot, err := Repositories.GetRobotByName(Common.FormatRobotUsername(kindName, robotShortname), kindId)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No robot found with name: %s for %s", robotShortname, kindName)
			return Dto.Robot{}, Errors.RobotNotFound(Common.FormatRobotUsername(kindName, robotShortname))
		default:
			logger.Error("Error checking existing robots: %s", err.Error())
			return Dto.Robot{}, err
		}
	}

	// Retrieve robot with details
	robotModel, err := Repositories.GetRobotById(existingRobot.ID)
	if err != nil {
		switch err.Error() {
		case "record not found": // This should never exist
			logger.Warning("No robot found with id: %d", existingRobot.ID)
			return Dto.Robot{}, Errors.RobotNotFoundById(existingRobot.ID)
		default:
			logger.Error("Error checking existing robots: %s", err.Error())
			return Dto.Robot{}, err
		}
	}

	// Convert model to dto
	robotDTO := Dto.Robot{
		Name:         robotModel.Username,
		Description:  robotModel.RobotAccountMetadata.Description,
		Created:      robotModel.CreationDate.Time,
		LastAccessed: Common.ConvertSQLNullTimeToTime(robotModel.LastAccessed),
	}

	// Include token
	if robotModel.RobotAccountToken != nil {
		token, decryptErr := Common.DecryptAESCipherToken(robotModel.RobotAccountToken.Token)
		if decryptErr != nil {
			logger.Error("Error decrypting robot token for robot '%s': %v", robotModel.Username, decryptErr)
			return Dto.Robot{}, decryptErr
		}

		robotDTO.Token = &token
	}

	return robotDTO, nil
}

func DeleteUserRobot(robotShortname string, currentUser *Auth.AuthenticatedUser) error {
	return deleteRobot(robotShortname, "user", strconv.Itoa(currentUser.ID), &Auth.AuthenticatedUser{ID: currentUser.ID, Username: currentUser.Username})
}

func DeleteOrganizationRobot(robotShortname string, orgName string, currentUser *Auth.AuthenticatedUser) error {
	return deleteRobot(robotShortname, "organization", orgName, &Auth.AuthenticatedUser{ID: currentUser.ID, Username: currentUser.Username})
}

func deleteRobot(robotShortname string, kind string, kindIdOrName string, currentUser *Auth.AuthenticatedUser) error {
	logger.Info("[Robot Service] Delete %s Robot account", kind)
	logger.Debug("With robot shortname: %s", robotShortname)

	// Check if the user or org exists
	var kindId int
	var kindName string // need to create the robot name in the format: <kindName>+<robotName>
	if kind == "user" {
		userExist, err := Repositories.GetUserById(Common.ParseStringToInt(kindIdOrName)) // The Current user is already authenticated, it will exist
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user found with id: %s", kindIdOrName)
				return Errors.UserNotFoundById(Common.ParseStringToInt(kindIdOrName))
			default:
				logger.Error("Error retrieving repository  from database: %s", err.Error())
				return err
			}
		}

		kindId = userExist.ID
		kindName = userExist.Username
	} else if kind == "organization" {
		orgExist, err := Repositories.GetOrganizationByName(kindIdOrName)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No organization found with name: %s", kindIdOrName)
				return Errors.OrganizationNotFound(kindIdOrName)
			default:
				logger.Error("Error retrieving organization from database: %s", err.Error())
				return err
			}
		}

		kindId = orgExist.ID
		kindName = orgExist.Username
	}

	// Check if a robot with the same name already exists for this user/org
	existingRobot, err := Repositories.GetRobotByName(Common.FormatRobotUsername(kindName, robotShortname), kindId)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No robot found with name: %s for %s", robotShortname, kindName)
			return Errors.RobotNotFound(Common.FormatRobotUsername(kindName, robotShortname))
		default:
			logger.Error("Error checking existing robots: %s", err.Error())
			return err
		}
	}

	// Delete the robot account from the db
	err = Repositories.DeleteRobotAccount(existingRobot.ID)
	if err != nil {
		logger.Error("Error when deleting robot account: %s", err.Error())
		return err
	}

	logger.Info("Robot account '%s' deleted successfully", existingRobot.Username)

	return nil
}

func GetUserRobotPermissions(robotShortname string, currentUser *Auth.AuthenticatedUser) ([]Dto.RobotPermission, error) {
	return getRobotPermissions(robotShortname, "user", strconv.Itoa(currentUser.ID), currentUser)
}

func GetOrganizationRobotPermissions(orgName string, robotShortName string, currentUser *Auth.AuthenticatedUser) ([]Dto.RobotPermission, error) {
	return getRobotPermissions(robotShortName, "organization", orgName, currentUser)
}

func getRobotPermissions(robotShortname string, kind string, kindIdOrName string, currentUser *Auth.AuthenticatedUser) ([]Dto.RobotPermission, error) {
	logger.Info("[Robot Service] Get %s Robot account permissions", kind)
	logger.Debug("With robot shortname: %s", robotShortname)

	// Check if the user or org exists
	var kindId int
	var kindName string // need to create the robot name in the format: <kindName>+<robotName>
	if kind == "user" {
		userExist, err := Repositories.GetUserById(Common.ParseStringToInt(kindIdOrName)) // The Current user is already authenticated, it will exist
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user found with id: %s", kindIdOrName)
				return []Dto.RobotPermission{}, Errors.UserNotFoundById(Common.ParseStringToInt(kindIdOrName))
			default:
				logger.Error("Error retrieving repository  from database: %s", err.Error())
				return []Dto.RobotPermission{}, err
			}
		}

		kindId = userExist.ID
		kindName = userExist.Username
	} else if kind == "organization" {
		orgExist, err := Repositories.GetOrganizationByName(kindIdOrName)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No organization found with name: %s", kindIdOrName)
				return []Dto.RobotPermission{}, Errors.OrganizationNotFound(kindIdOrName)
			default:
				logger.Error("Error retrieving organization from database: %s", err.Error())
				return []Dto.RobotPermission{}, err
			}
		}

		kindId = orgExist.ID
		kindName = orgExist.Username
	}

	// Check if a robot with the same name already exists for this user/org
	existingRobot, err := Repositories.GetRobotByName(Common.FormatRobotUsername(kindName, robotShortname), kindId)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No robot found with name: %s for %s", robotShortname, kindName)
			return []Dto.RobotPermission{}, Errors.RobotNotFound(Common.FormatRobotUsername(kindName, robotShortname))
		default:
			logger.Error("Error checking existing robots: %s", err.Error())
			return []Dto.RobotPermission{}, err
		}
	}

	// Retrieve robot with details
	robotModel, err := Repositories.GetRobotById(existingRobot.ID)
	if err != nil {
		switch err.Error() {
		case "record not found": // This should never exist
			logger.Warning("No robot found with id: %d", existingRobot.ID)
			return []Dto.RobotPermission{}, Errors.RobotNotFoundById(existingRobot.ID)
		default:
			logger.Error("Error checking existing robots: %s", err.Error())
			return []Dto.RobotPermission{}, err
		}
	}

	// Convert model to dto
	permissions := []Dto.RobotPermission{}
	for _, repoPerm := range robotModel.RepositoryPermissions {
		permissions = append(permissions, Dto.RobotPermission{
			Repository: Dto.RobotPermissionRepository{
				Name:     repoPerm.Repository.Name,
				IsPublic: Common.MapRepositoryVisibilityById(repoPerm.Repository.VisibilityId) == "public",
			},
			Role: Common.MapRoleById(repoPerm.RoleId),
		})
	}

	return permissions, nil
}

func ListUserRobotFederations(robotShortname string, currentUser *Auth.AuthenticatedUser) ([]Dto.RobotFederation, error) {
	return listRobotFederations(robotShortname, "user", strconv.Itoa(currentUser.ID), currentUser)
}

func ListOrganizationRobotFederations(orgName string, robotShortname string, currentUser *Auth.AuthenticatedUser) ([]Dto.RobotFederation, error) {
	return listRobotFederations(robotShortname, "organization", orgName, currentUser)
}

func listRobotFederations(robotShortname string, kind string, kindIdOrName string, currentUser *Auth.AuthenticatedUser) ([]Dto.RobotFederation, error) {
	logger.Info("[Robot Service] Get %s Robot account federation", kind)
	logger.Debug("With robot shortname: %s", robotShortname)

	// Check if the user or org exists
	var kindId int
	var kindName string // need to create the robot name in the format: <kindName>+<robotName>
	if kind == "user" {
		userExist, err := Repositories.GetUserById(Common.ParseStringToInt(kindIdOrName)) // The Current user is already authenticated, it will exist
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user found with id: %s", kindIdOrName)
				return []Dto.RobotFederation{}, Errors.UserNotFoundById(Common.ParseStringToInt(kindIdOrName))
			default:
				logger.Error("Error retrieving repository  from database: %s", err.Error())
				return []Dto.RobotFederation{}, err
			}
		}

		kindId = userExist.ID
		kindName = userExist.Username
	} else if kind == "organization" {
		orgExist, err := Repositories.GetOrganizationByName(kindIdOrName)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No organization found with name: %s", kindIdOrName)
				return []Dto.RobotFederation{}, Errors.OrganizationNotFound(kindIdOrName)
			default:
				logger.Error("Error retrieving organization from database: %s", err.Error())
				return []Dto.RobotFederation{}, err
			}
		}

		kindId = orgExist.ID
		kindName = orgExist.Username
	}

	// Check if a robot with the same name already exists for this user/org
	existingRobot, err := Repositories.GetRobotByName(Common.FormatRobotUsername(kindName, robotShortname), kindId)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No robot found with name: %s for %s", robotShortname, kindName)
			return []Dto.RobotFederation{}, Errors.RobotNotFound(Common.FormatRobotUsername(kindName, robotShortname))
		default:
			logger.Error("Error checking existing robots: %s", err.Error())
			return []Dto.RobotFederation{}, err
		}
	}

	// Retrieve robot with details
	robotModel, err := Repositories.GetRobotById(existingRobot.ID)
	if err != nil {
		switch err.Error() {
		case "record not found": // This should never exist
			logger.Warning("No robot found with id: %d", existingRobot.ID)
			return []Dto.RobotFederation{}, Errors.RobotNotFoundById(existingRobot.ID)
		default:
			logger.Error("Error checking existing robots: %s", err.Error())
			return []Dto.RobotFederation{}, err
		}
	}

	// Convert model to dto
	federations := []Dto.RobotFederation{}
	for _, federatedLogin := range robotModel.FederatedLogins {
		var metadataJson Models.FederatedLoginMetadata
		err = json.Unmarshal([]byte(federatedLogin.MetadataJson), &metadataJson)
		if err != nil {
			logger.Error("Error unmarshaling robot federated login metadata: %v", err)
			return []Dto.RobotFederation{}, err
		}

		for _, metadata := range metadataJson.FederationConfig {
			robotFederation := Dto.RobotFederation{
				Issuer:  metadata.Issuer,
				Subject: metadata.Subject,
			}
			federations = append(federations, robotFederation)
		}
	}

	return federations, nil
}

func CreateOrUpdateUserRobotFederations(robotShortname string, robotFederationsToCreate []Dto.RobotFederation, currentUser *Auth.AuthenticatedUser) ([]Dto.RobotFederation, error) {
	return createOrUpdateRobotFederations(robotShortname, robotFederationsToCreate, "user", strconv.Itoa(currentUser.ID), currentUser)
}

func CreateOrUpdateOrganizationRobotFederations(orgName string, robotShortname string, robotFederationsToCreate []Dto.RobotFederation, currentUser *Auth.AuthenticatedUser) ([]Dto.RobotFederation, error) {
	return createOrUpdateRobotFederations(robotShortname, robotFederationsToCreate, "organization", orgName, currentUser)
}

func createOrUpdateRobotFederations(robotShortname string, robotFederationsToCreate []Dto.RobotFederation, kind string, kindIdOrName string, currentUser *Auth.AuthenticatedUser) ([]Dto.RobotFederation, error) {
	logger.Info("[Robot Service] Create or update %s Robot account federations", kind)
	logger.Debug("With robot shortname: %s", robotShortname)

	// Validate input
	if err := Common.ValidateCreateRobotAccountFederations(robotFederationsToCreate); err != nil {
		logger.Error("Input validation error: %v", err)
		return []Dto.RobotFederation{}, err
	}

	// Check if the user or org exists
	var kindId int
	var kindName string // need to create the robot name in the format: <kindName>+<robotName>
	if kind == "user" {
		userExist, err := Repositories.GetUserById(Common.ParseStringToInt(kindIdOrName)) // The Current user is already authenticated, it will exist
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user found with id: %s", kindIdOrName)
				return []Dto.RobotFederation{}, Errors.UserNotFoundById(Common.ParseStringToInt(kindIdOrName))
			default:
				logger.Error("Error retrieving repository  from database: %s", err.Error())
				return []Dto.RobotFederation{}, err
			}
		}

		kindId = userExist.ID
		kindName = userExist.Username
	} else if kind == "organization" {
		orgExist, err := Repositories.GetOrganizationByName(kindIdOrName)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No organization found with name: %s", kindIdOrName)
				return []Dto.RobotFederation{}, Errors.OrganizationNotFound(kindIdOrName)
			default:
				logger.Error("Error retrieving organization from database: %s", err.Error())
				return []Dto.RobotFederation{}, err
			}
		}

		kindId = orgExist.ID
		kindName = orgExist.Username
	}

	// Check if a robot with the same name already exists for this user/org
	existingRobot, err := Repositories.GetRobotByName(Common.FormatRobotUsername(kindName, robotShortname), kindId)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No robot found with name: %s for %s", robotShortname, kindName)
			return []Dto.RobotFederation{}, Errors.RobotNotFound(Common.FormatRobotUsername(kindName, robotShortname))
		default:
			logger.Error("Error checking existing robots: %s", err.Error())
			return []Dto.RobotFederation{}, err
		}
	}

	// Retrieve robot with details
	robotModel, err := Repositories.GetRobotById(existingRobot.ID)
	if err != nil {
		switch err.Error() {
		case "record not found": // This should never exist
			logger.Warning("No robot found with id: %d", existingRobot.ID)
			return []Dto.RobotFederation{}, Errors.RobotNotFoundById(existingRobot.ID)
		default:
			logger.Error("Error checking existing robots: %s", err.Error())
			return []Dto.RobotFederation{}, err
		}
	}

	// All validations passed
	logger.Info("All validations passed")

	// Convert dto to model
	var metadataToCreate Models.FederatedLoginMetadata
	for _, metadata := range robotFederationsToCreate {
		metadataToCreate.FederationConfig = append(metadataToCreate.FederationConfig, Models.FederationConfig{
			Issuer:  metadata.Issuer,
			Subject: metadata.Subject,
		})
	}

	// Parse struct ot json string
	jsonMetadata, err := json.Marshal(metadataToCreate)
	if err != nil {
		logger.Error("Error marshaling metadata to JSON: %v", err)
		return []Dto.RobotFederation{}, err
	}

	// Update the field MetadataJson from FederatedLogin
	err = Repositories.UpdateFederatedLoginMetadata(robotModel.FederatedLogins[0].ID, string(jsonMetadata))
	if err != nil {
		logger.Error("Error updating federated login metadata: %v", err)
		return []Dto.RobotFederation{}, err
	}

	return robotFederationsToCreate, nil
}

// <editor-fold desc="Private Methods">

func metadataMapToJSON(metadata map[string]interface{}) string {
	unstructuredMetadataJSON := "{}"
	if metadata != nil {
		metadataBytes, marshalErr := json.Marshal(metadata)
		if marshalErr != nil {
			logger.Error("Error marshaling robot unstructured metadata: %v", marshalErr)
			return "{}"
		}

		unstructuredMetadataJSON = string(metadataBytes)
	}
	return unstructuredMetadataJSON
}

func generateToken() string {
	return Common.RandomStringGenerator(64)
}

/*
generateEncryptedToken Generate a new token and encrypt it.
Returns the raw token, the encrypted token and encryption error
*/
func generateEncryptedToken() (string, string, error) {
	// Generate raw token
	rawToken := generateToken()

	// Encrypt the raw token
	encryptedToken, encryptErr := Common.EncryptAESCipherToken(rawToken)
	if encryptErr != nil {
		logger.Error("Error encrypting robot token: %v", encryptErr)
		return "", "", encryptErr
	}
	return rawToken, encryptedToken, nil
}

// </editor-fold>
