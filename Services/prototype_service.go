package Services

import (
	"quay-go-api/Common"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
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
