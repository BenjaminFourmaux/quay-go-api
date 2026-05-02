package Services

import (
	"net/mail"
	"quay-go-api/Common"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Entities/Models"
	"quay-go-api/Repositories"
	"quay-go-api/Services/Auth"
	"strings"
)

func GetUserOrganizations(currentUser Auth.AuthenticatedUser) ([]Dto.UserOrganization, error) {
	user, err := Repositories.GetUserById(currentUser.ID)
	if err != nil {
		switch err.Error() {
		case "record not found":
			return nil, Errors.CurrentUserNotFound()
		default:
			return nil, err
		}
	}

	orgsModel, err := Repositories.GetUserOrganizations(user.ID)
	if err != nil {
		switch err.Error() {
		case "record not found": // if no result, is not an error, just return empty list
			return []Dto.UserOrganization{}, nil
		default:
			return nil, err
		}
	}

	// Convert userModel into UserOrganization Dto
	organizations := Common.ConvertUserModelsToDto(orgsModel, user, currentUser.Scopes)
	return organizations, nil
}

func CreateOrganization(organizationMetadata Dto.CreateOrganization, currentUser Auth.AuthenticatedUser) (Dto.Organization, error) {
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
		Username:     strings.ToLower(organizationMetadata.Name),
		Organization: true,
	}

	createdOrgModel, err := Repositories.CreateOrganizationWithOwnerTeamTransaction(createOrgModel, currentUser.ID)
	if err != nil {
		return Dto.Organization{}, err
	}

	// Get the new org with details
	createdOrgModel, err = Repositories.GetOrganizationDetailsById(createdOrgModel.ID)
	if err != nil {
		return Dto.Organization{}, err
	}

	// Convert model to dto
	createdOrgDto := Common.ConvertUserModelToOrganizationDto(createdOrgModel, currentUser.ID, currentUser.Scopes)

	return createdOrgDto, nil
}

func GetOrganizationDetailsByName(orgName string, currentUser Auth.AuthenticatedUser) (Dto.Organization, error) {
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

	orgDetailDto := Common.ConvertUserModelToOrganizationDto(orgModel, currentUser.ID, currentUser.Scopes)

	return orgDetailDto, nil
}

func DeleteOrganization(orgName string, currentUser Auth.AuthenticatedUser) error {
	// Get the org (with detail, for user role checking) if exists
	organizationToDeleteModel, err := Repositories.GetOrganizationDetailsByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			return Errors.OrganizationNotFound(orgName)
		default:
			return err
		}
	}

	// If the user is owners of the organization so it can delete this
	if isUserIsOrgOwner(currentUser.ID, organizationToDeleteModel) {
		err = Repositories.DeleteOrganizationTransaction(organizationToDeleteModel.ID)
		if err != nil {
			return err
		}
		// TODO: remove the namespace and Docker images ?
	}
	return nil
}

func UpdateOrganization(orgName string, organizationMetadata Dto.UpdateOrganization, currentUser Auth.AuthenticatedUser) (Dto.Organization, error) {
	// Get the org (with detail, for user role checking) if exists
	organizationToUpdateModel, err := Repositories.GetOrganizationDetailsByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			return Dto.Organization{}, Errors.OrganizationNotFound(orgName)
		default:
			return Dto.Organization{}, err
		}
	}

	// If the user is owners of the organization so it can update this
	if isUserIsOrgOwner(currentUser.ID, organizationToUpdateModel) {
		if err = validateUpdateOrganization(organizationMetadata); err != nil {
			return Dto.Organization{}, err
		}

		updatedFields := make(map[string]interface{})

		if organizationMetadata.Email != nil {
			updatedFields["email"] = strings.TrimSpace(*organizationMetadata.Email)
		}

		if organizationMetadata.InvoiceEmail != nil {
			updatedFields["invoice_email"] = *organizationMetadata.InvoiceEmail
		}

		if organizationMetadata.InvoiceEmailAddress != nil {
			updatedFields["invoice_email_address"] = *organizationMetadata.InvoiceEmailAddress
		}

		if organizationMetadata.TagExpirationS != nil {
			updatedFields["removed_tag_expiration_s"] = *organizationMetadata.TagExpirationS
		}

		if err = Repositories.UpdateOrganizationFieldsById(organizationToUpdateModel.ID, updatedFields); err != nil {
			return Dto.Organization{}, err
		}

		organizationToUpdateModel, err = Repositories.GetOrganizationDetailsById(organizationToUpdateModel.ID)
		if err != nil {
			return Dto.Organization{}, err
		}

		updatedOrgDto := Common.ConvertUserModelToOrganizationDto(organizationToUpdateModel, currentUser.ID, currentUser.Scopes)
		return updatedOrgDto, nil
	}

	return Dto.Organization{}, Errors.UserNotOrganizationOwner()
}

// <editor-fold desc="Private Methods">

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

func validateUpdateOrganization(organizationMetadata Dto.UpdateOrganization) error {
	if organizationMetadata.Email != nil {
		email := strings.TrimSpace(*organizationMetadata.Email)
		if email == "" {
			return Errors.OrganizationEmailInvalid()
		}

		if _, err := mail.ParseAddress(email); err != nil {
			return Errors.OrganizationEmailInvalid()
		}
	}

	if organizationMetadata.InvoiceEmailAddress != nil {
		email := strings.TrimSpace(*organizationMetadata.InvoiceEmailAddress)
		if email == "" {
			return Errors.OrganizationEmailInvalid()
		}

		if _, err := mail.ParseAddress(email); err != nil {
			return Errors.OrganizationEmailInvalid()
		}
	}

	if organizationMetadata.TagExpirationS != nil && *organizationMetadata.TagExpirationS < 0 {
		return Errors.OrganizationTagExpirationInvalid()
	}

	return nil
}

/*
isUserIsOrgOwner checks if the user is on the 'owners' team of the organization
*/
func isUserIsOrgOwner(userId int, organization Models.User) bool {
	for _, team := range organization.Teams {
		if team.Name == "owners" || team.RoleId == 1 {
			for _, member := range team.Members {
				if member.UserId == userId {
					return true
				}
			}
		}
	}
	return false // the user isn't in 'owners'
}

// </editor-fold>
