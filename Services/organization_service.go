package Services

import (
	"quay-go-api/Common"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Entities/Models"
	"quay-go-api/Repositories"
	"quay-go-api/Services/Auth"
)

func GetUserOrganizations(userId int, userScopes []Auth.Scope) ([]Dto.UserOrganization, error) {
	// Get the current user
	currentUser, err := Repositories.GetUserById(userId)
	if err != nil {
		switch err.Error() {
		case "record not found":
			return nil, Errors.CurrentUserNotFound()
		default:
			return nil, err
		}
	}

	orgsModel, err := Repositories.GetUserOrganizations(userId)
	if err != nil {
		switch err.Error() {
		case "record not found": // if no result, is not an error, just return empty list
			return []Dto.UserOrganization{}, nil
		default:
			return nil, err
		}
	}

	// Convert userModel into UserOrganization Dto
	organizations := Common.ConvertUserModelsToDto(orgsModel, currentUser, userScopes)
	return organizations, nil
}

func CreateOrganization(organizationMetadata Dto.CreateOrganization, userId int, userScopes []Auth.Scope) (Dto.Organization, error) {
	// TODO: check if user has role to create org and check Features flag (not implement yet)

	// Check if the org (or a user) already exists
	existingOrg, err := Repositories.GetUserOrOrganizationByName(organizationMetadata.Name)
	if err == nil && existingOrg.Username == organizationMetadata.Name {
		return Dto.Organization{}, Errors.UserOrOrganizationAlreadyExists()
	}

	// validating org
	err = validateCreateOrganization(organizationMetadata)
	if err != nil {
		return Dto.Organization{}, err
	}

	// Convert dto to user model
	var createOrgModel = Models.User{
		Username:     organizationMetadata.Name,
		Organization: true,
	}

	createdOrgModel, err := Repositories.CreateOrganizationWithOwnerTeamTransaction(createOrgModel, userId)
	if err != nil {
		return Dto.Organization{}, err
	}

	// Get the new org with details
	createdOrgModel, err = Repositories.GetOrganizationDetailsById(createdOrgModel.ID)
	if err != nil {
		return Dto.Organization{}, err
	}

	// Convert model to dto
	createdOrgDto := Common.ConvertUserModelToOrganizationDto(createdOrgModel, userId, userScopes)

	return createdOrgDto, nil
}

func GetOrganizationDetailsByName(orgName string, userId int, userScopes []Auth.Scope) (Dto.Organization, error) {
	// Get the current user
	currentUser, err := Repositories.GetUserById(userId)
	if err != nil {
		switch err.Error() {
		case "record not found":
			return Dto.Organization{}, Errors.CurrentUserNotFound()
		default:
			return Dto.Organization{}, err
		}
	}

	// Get the organization with details
	orgModel, err := Repositories.GetOrganizationDetailsByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			return Dto.Organization{}, Errors.OrganizationNotFound(orgName)
		default:
			return Dto.Organization{}, err
		}
	}

	orgDetailDto := Common.ConvertUserModelToOrganizationDto(orgModel, currentUser.ID, userScopes)

	return orgDetailDto, nil
}

/*
Validate organization metadata for creating a new organization. Rules:
The organization name must:
1. Not be empty
2. Be between 2 and 255 characters long
3. Contain only alphanumeric characters, dashes, or underscores
*/
func validateCreateOrganization(organizationMetadata Dto.CreateOrganization) error {
	// Validate org name
	// 1. empty value
	if organizationMetadata.Name == "" {
		return Errors.OrganizationNameInvalid()
	}

	// 2. value length
	if len(organizationMetadata.Name) < 2 || len(organizationMetadata.Name) > 255 {
		return Errors.OrganizationNameInvalid()
	}

	// 3. valid characters (alphanumeric, dash and underscore)
	for _, char := range organizationMetadata.Name {
		if !(char >= 'a' && char <= 'z') && !(char >= '0' && char <= '9') && char != '-' && char != '_' {
			return Errors.OrganizationNameInvalid()
		}
	}

	return nil
}
