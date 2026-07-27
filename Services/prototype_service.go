package Services

import (
	"quay-go-api/Common"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Entities/Models"
	"quay-go-api/Repositories"
	"quay-go-api/Services/Auth"
	logger "quay-go-api/Services/Logger"
)

func ListOrganizationPrototypes(orgName string, currentUser Auth.AuthenticatedUser) ([]Dto.Prototype, error) {
	logger.Info("[Prototype Service] List Organization Prototypes")
	logger.Debug("Organization: %s", orgName)

	// Retrieve organization and check if exists
	logger.Info("Retrieving organization details from database")
	orgModel, err := Repositories.GetOrganizationDetailsByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("Organization not found: %s", orgName)
			return []Dto.Prototype{}, Errors.OrganizationNotFound(orgName)
		default:
			logger.Error("Error retrieving organization details from database: %s", err.Error())
			return []Dto.Prototype{}, err
		}
	}

	// Retrieve the list of organization's prototypes
	prototypesModel, err := Repositories.GetOrganizationPrototypes(orgModel.ID)
	if err != nil {
		logger.Error("Error retrieving organization prototypes from database: %s", err.Error())
		return []Dto.Prototype{}, err
	}

	// convert model to dto
	var prototypes []Dto.Prototype
	for _, prototypeModel := range prototypesModel {
		// Determine if user/team are org member
		var activatingUserOrgMember bool = false
		var delegateOrgMember bool = false

		if prototypeModel.ActivatingUserId != nil {
			activatingUserOrgMember = isUserOrgMember(*prototypeModel.ActivatingUserId, orgModel.ID)
		}
		if prototypeModel.DelegateUserId != nil {
			delegateOrgMember = isUserOrgMember(*prototypeModel.DelegateUserId, orgModel.ID)
		}
		if prototypeModel.DelegateTeamId != nil {
			delegateOrgMember = isTeamOrgMember(*prototypeModel.DelegateTeamId, orgModel.ID)
		}

		prototypes = append(prototypes, Common.ConvertPermissionPrototypeModelToDto(prototypeModel, activatingUserOrgMember, delegateOrgMember))
	}

	return prototypes, nil
}

func CreatePrototype(orgName string, prototypeDto Dto.CreatePrototype, currentUser Auth.AuthenticatedUser) (Dto.Prototype, error) {
	logger.Info("[Prototype Service] Create Prototype")
	logger.Info("Organization: %s", orgName)

	// Validate input
	if !Common.IsValidRepositoryPermissionRole(prototypeDto.Role) {
		return Dto.Prototype{}, Errors.RepositoryPermissionRoleInvalid(prototypeDto.Role)
	}
	if !Common.IsValidPrototypeDelegateKind(prototypeDto.DelegateKind) {
		return Dto.Prototype{}, Errors.PrototypeDelegateKindInvalid(prototypeDto.DelegateKind)
	}

	// Retrieve organization and check if exists
	logger.Info("Retrieving organization details from database")
	orgModel, err := Repositories.GetOrganizationDetailsByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("Organization not found: %s", orgName)
			return Dto.Prototype{}, Errors.OrganizationNotFound(orgName)
		default:
			logger.Error("Error retrieving organization details from database: %s", err.Error())
			return Dto.Prototype{}, err
		}
	}

	// Chek if the Activating user exists
	activatingUserModel, err := Repositories.GetUserByName(prototypeDto.ActivatingUserName)
	if err != nil {
		return Dto.Prototype{}, Errors.UserNotFound(prototypeDto.ActivatingUserName)
	}

	// Chek if the kind user and user exist
	var delegateUserId int
	if prototypeDto.DelegateKind == "user" {
		delegateUserModel, err := Repositories.GetUserByName(prototypeDto.DelegateName)
		if err != nil {
			return Dto.Prototype{}, Errors.UserNotFound(prototypeDto.DelegateName)
		}
		delegateUserId = delegateUserModel.ID
	}

	// Check if the kind teal and team exist
	var delegateTeamId int
	if prototypeDto.DelegateKind == "team" {
		delegateTeamModel, err := Repositories.GetTeamByName(prototypeDto.DelegateName)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("Team not found: %s", prototypeDto.DelegateName)
				return Dto.Prototype{}, Errors.TeamNotFound(prototypeDto.DelegateName)
			default:
				logger.Error("Error retrieving team details from database: %s", err.Error())
				return Dto.Prototype{}, err
			}
		}
		delegateTeamId = delegateTeamModel.ID
	}

	// Create model
	createPrototype := Models.PermissionPrototype{
		UUID:             Common.GenerateUUID(),
		OrgId:            orgModel.ID,
		RoleId:           Common.GetRoleIdFromRoleName(prototypeDto.Role),
		ActivatingUserId: &activatingUserModel.ID,
	}

	if prototypeDto.DelegateKind == "user" {
		createPrototype.DelegateUserId = &delegateUserId
	}

	if prototypeDto.DelegateKind == "team" {
		createPrototype.DelegateTeamId = &delegateTeamId
	}

	// Create the model into the database
	createdPrototypeModel, err := Repositories.CreatePermissionPrototype(&createPrototype)
	if err != nil {
		logger.Error("Error creating prototype in database: %s", err.Error())
		return Dto.Prototype{}, err
	}

	// Determine if user/team are org member
	var activatingUserOrgMember bool = false
	var delegateOrgMember bool = false

	if createdPrototypeModel.ActivatingUserId != nil {
		activatingUserOrgMember = isUserOrgMember(*createdPrototypeModel.ActivatingUserId, orgModel.ID)
	}
	if createdPrototypeModel.DelegateUserId != nil {
		delegateOrgMember = isUserOrgMember(*createdPrototypeModel.DelegateUserId, orgModel.ID)
	}
	if createdPrototypeModel.DelegateTeamId != nil {
		delegateOrgMember = isTeamOrgMember(*createdPrototypeModel.DelegateTeamId, orgModel.ID)
	}

	return Common.ConvertPermissionPrototypeModelToDto(*createdPrototypeModel, activatingUserOrgMember, delegateOrgMember), nil
}
