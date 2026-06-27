package Services

import (
	"quay-go-api/Common"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Entities/Models"
	"quay-go-api/Repositories"
	"quay-go-api/Services/Auth"
	"quay-go-api/Services/Avatar"
	logger "quay-go-api/Services/Logger"
	"slices"
	"strconv"
	"strings"
)

func GetUserOrganizations(currentUser Auth.AuthenticatedUser, filters map[string]string) ([]Dto.UserOrganization, error) {
	logger.Info("[Organization Service] Get User Organizations")
	logger.Debug("Filters: %+v", filters)

	// Validating filters
	var filterPublic bool
	if public, ok := filters["is_public"]; ok {
		isPublic, err := strconv.ParseBool(public)
		if err != nil {
			logger.Warning("Invalid filter is_public value: %s", public)
			return nil, Errors.InvalidParameterValue("is_public", []string{"true", "false"})
		}
		filterPublic = isPublic
	}

	logger.Info("Retrieving current user from database")
	user, err := Repositories.GetUserById(currentUser.ID)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("Current user not found in database: %d", currentUser.ID)
			return nil, Errors.CurrentUserNotFound()
		default:
			logger.Error("Error retrieving current user from database: %s", err.Error())
			return nil, err
		}
	}

	logger.Info("Retrieving user organizations from database")
	orgsModel, err := Repositories.GetUserOrganizations(user.ID)
	if err != nil {
		switch err.Error() {
		case "record not found": // if no result, is not an error, just return empty list
			logger.Debug("No organizations found for user ID: %d", user.ID)
			return []Dto.UserOrganization{}, nil
		default:
			logger.Error("Error retrieving user organizations from database: %s", err.Error())
			return nil, err
		}
	}

	// Convert userModel into UserOrganization Dto
	organizations := Common.ConvertUserModelsToDto(orgsModel, user, currentUser.Scopes)

	// Apply filters
	if _, ok := filters["is_public"]; ok {
		logger.Debug("Applying is_public filter: %t", filterPublic)
		organizations = slices.Collect(func(yield func(Dto.UserOrganization) bool) {
			for _, org := range organizations {
				if org.Public == filterPublic {
					if !yield(org) {
						return
					}
				}
			}
		})
		if organizations == nil {
			logger.Debug("No organizations left after filters")
			return []Dto.UserOrganization{}, nil
		}
	}

	logger.Debug("Organizations returned: %d", len(organizations))

	return organizations, nil
}

func CreateOrganization(organizationMetadata Dto.CreateOrganization, currentUser Auth.AuthenticatedUser) (Dto.Organization, error) {
	logger.Info("[Organization Service] Create Organization")
	logger.Debug("With dto: %+v", organizationMetadata)

	// TODO: check if user has role to create org and check Features flag (not implement yet)

	// Check if the org (or a user) already exists
	logger.Info("Checking if organization already exists in database")
	existingOrg, err := Repositories.GetUserOrOrganizationByName(organizationMetadata.Name)
	if err == nil && existingOrg.Username == organizationMetadata.Name {
		logger.Warning("Organization or user already exists with name: %s", organizationMetadata.Name)
		return Dto.Organization{}, Errors.UserOrOrganizationAlreadyExists()
	}

	// validating org
	err = Common.ValidateCreateOrganization(organizationMetadata)
	if err != nil {
		logger.Warning("Invalid create organization payload: %s", err.Error())
		return Dto.Organization{}, err
	}

	// Convert dto to user model
	var createOrgModel = Models.User{
		Username:     strings.ToLower(organizationMetadata.Name),
		Organization: true,
	}

	logger.Info("Creating organization in database")
	createdOrgModel, err := Repositories.CreateOrganizationWithOwnerTeamTransaction(createOrgModel, currentUser.ID)
	if err != nil {
		logger.Error("Error creating organization in database: %s", err.Error())
		return Dto.Organization{}, err
	}

	// Get the new org with details
	logger.Info("Retrieving created organization details from database")
	createdOrgModel, err = Repositories.GetOrganizationDetailsById(createdOrgModel.ID)
	if err != nil {
		logger.Error("Error retrieving created organization details from database: %s", err.Error())
		return Dto.Organization{}, err
	}

	// Convert model to dto
	createdOrgDto := Common.ConvertUserModelToOrganizationDto(createdOrgModel, currentUser.ID, currentUser.Scopes)
	logger.Success("Organization created successfully: %s", createdOrgDto.Name)

	return createdOrgDto, nil
}

func GetOrganizationDetailsByName(orgName string, currentUser Auth.AuthenticatedUser) (Dto.Organization, error) {
	logger.Info("[Organization Service] Get Organization Details")
	logger.Debug("Organization name: %s", orgName)

	// Retrieve organization and check if exists
	logger.Info("Retrieving organization details from database")
	orgModel, err := Repositories.GetOrganizationDetailsByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("Organization not found: %s", orgName)
			return Dto.Organization{}, Errors.OrganizationNotFound(orgName)
		default:
			logger.Error("Error retrieving organization details from database: %s", err.Error())
			return Dto.Organization{}, err
		}
	}

	orgDetailDto := Common.ConvertUserModelToOrganizationDto(orgModel, currentUser.ID, currentUser.Scopes)
	logger.Debug("Organization details retrieved successfully: %s", orgDetailDto.Name)

	return orgDetailDto, nil
}

func DeleteOrganization(orgName string, currentUser Auth.AuthenticatedUser) error {
	logger.Info("[Organization Service] Delete Organization")
	logger.Debug("Organization name: %s", orgName)

	// Get the org (with detail, for user role checking) if exists
	logger.Info("Retrieving organization details from database")
	organizationToDeleteModel, err := Repositories.GetOrganizationDetailsByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("Organization not found for deletion: %s", orgName)
			return Errors.OrganizationNotFound(orgName)
		default:
			logger.Error("Error retrieving organization to delete from database: %s", err.Error())
			return err
		}
	}

	// If the user is owners of the organization so it can delete this
	if isUserIsOrgOwner(currentUser.ID, organizationToDeleteModel) {
		logger.Info("Deleting organization in database")
		err = Repositories.DeleteOrganizationTransaction(organizationToDeleteModel.ID)
		if err != nil {
			logger.Error("Error deleting organization in database: %s", err.Error())
			return err
		}
		logger.Success("Organization deleted successfully: %s", orgName)
		// TODO: remove the namespace and Docker images ?
	} else {
		logger.Warning("Delete denied: user %d is not owner of organization %s", currentUser.ID, orgName)
	}
	return nil
}

func UpdateOrganization(orgName string, organizationMetadata Dto.UpdateOrganization, currentUser Auth.AuthenticatedUser) (Dto.Organization, error) {
	logger.Info("[Organization Service] Update Organization")
	logger.Debug("Organization name: %s", orgName)
	logger.Debug("With dto: %+v", organizationMetadata)

	// Get the org (with detail, for user role checking) if exists
	logger.Info("Retrieving organization details from database")
	organizationToUpdateModel, err := Repositories.GetOrganizationDetailsByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("Organization not found for update: %s", orgName)
			return Dto.Organization{}, Errors.OrganizationNotFound(orgName)
		default:
			logger.Error("Error retrieving organization to update from database: %s", err.Error())
			return Dto.Organization{}, err
		}
	}

	// If the user is owners of the organization so it can update this
	if isUserIsOrgOwner(currentUser.ID, organizationToUpdateModel) {
		if err = Common.ValidateUpdateOrganization(organizationMetadata); err != nil {
			logger.Warning("Invalid update organization payload: %s", err.Error())
			return Dto.Organization{}, err
		}

		// Select fields to update
		mappings := map[string]Common.UpdateFieldMapping{}

		if organizationMetadata.Email != nil {
			mappings["Email"] = Common.UpdateFieldMapping{
				ModelFieldName: "Email",
				Value:          strings.TrimSpace(*organizationMetadata.Email),
			}
		}

		if organizationMetadata.InvoiceEmail != nil {
			mappings["InvoiceEmail"] = Common.UpdateFieldMapping{
				ModelFieldName: "InvoiceEmail",
				Value:          *organizationMetadata.InvoiceEmail,
			}
		}

		if organizationMetadata.InvoiceEmailAddress != nil {
			mappings["InvoiceEmailAddress"] = Common.UpdateFieldMapping{
				ModelFieldName: "InvoiceEmailAddress",
				Value:          *organizationMetadata.InvoiceEmailAddress,
			}
		}

		if organizationMetadata.TagExpirationS != nil {
			mappings["TagExpirationS"] = Common.UpdateFieldMapping{
				ModelFieldName: "TagExpirationS",
				Value:          *organizationMetadata.TagExpirationS,
			}
		}

		updatedFields := Common.BuildUpdatedFields[Models.User](organizationMetadata, mappings)
		logger.Debug("Organization fields to update: %+v", updatedFields)

		logger.Info("Updating organization in database")
		if err = Repositories.UpdateOrganizationFieldsById(organizationToUpdateModel.ID, updatedFields); err != nil {
			logger.Error("Error updating organization in database: %s", err.Error())
			return Dto.Organization{}, err
		}

		logger.Info("Retrieving updated organization details from database")
		organizationToUpdateModel, err = Repositories.GetOrganizationDetailsById(organizationToUpdateModel.ID)
		if err != nil {
			logger.Error("Error retrieving updated organization details from database: %s", err.Error())
			return Dto.Organization{}, err
		}

		updatedOrgDto := Common.ConvertUserModelToOrganizationDto(organizationToUpdateModel, currentUser.ID, currentUser.Scopes)
		logger.Success("Organization updated successfully: %s", updatedOrgDto.Name)
		return updatedOrgDto, nil
	}

	logger.Warning("Update denied: user %d is not owner of organization %s", currentUser.ID, orgName)
	return Dto.Organization{}, Errors.UserNotOrganizationOwner()
}

func ListMembersOfOrganization(orgName string, currentUser Auth.AuthenticatedUser) ([]Dto.OrganizationMember, error) {
	logger.Info("[Organization Service] List Organization Members")
	logger.Debug("Organization name: %s", orgName)

	// Retrieve organization and check if exists
	logger.Info("Retrieving organization details from database")
	organizationModel, err := Repositories.GetOrganizationDetailsByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("Organization not found while listing members: %s", orgName)
			return nil, Errors.OrganizationNotFound(orgName)
		default:
			logger.Error("Error retrieving organization members from database: %s", err.Error())
			return nil, err
		}
	}

	// Chek if user has the right to see members
	if !Common.HasScope(currentUser.Scopes, Auth.OrgAdmin) &&
		!Common.HasScope(currentUser.Scopes, Auth.SuperUser) &&
		!isUserIsOrgOwner(currentUser.ID, organizationModel) {
		logger.Warning("List members denied: user %d has insufficient role on organization %s", currentUser.ID, orgName)
		return nil, Errors.UnauthorizedInsufficientRole()
	}

	membersByUserId := make(map[int]*Dto.OrganizationMember)
	memberOrder := make([]int, 0)

	for _, team := range organizationModel.Teams {
		for _, teamMember := range team.Members {
			user := teamMember.User
			member, exists := membersByUserId[user.ID]
			if !exists {
				avatar := Avatar.GetAvatarForUser(user)
				member = &Dto.OrganizationMember{
					Name:         user.Username,
					Kind:         "user",
					Avatar:       avatar,
					Teams:        []string{},
					Repositories: []string{},
				}
				membersByUserId[user.ID] = member
				memberOrder = append(memberOrder, user.ID)
			}

			alreadyInTeam := false
			for _, existingTeamName := range member.Teams {
				if existingTeamName == team.Name {
					alreadyInTeam = true
					break
				}
			}

			if !alreadyInTeam {
				member.Teams = append(member.Teams, team.Name)
			}
		}
	}

	members := make([]Dto.OrganizationMember, 0, len(memberOrder))
	for _, userId := range memberOrder {
		members = append(members, *membersByUserId[userId])
	}

	logger.Debug("Organization members returned: %d", len(members))

	return members, nil

}

// <editor-fold desc="Private Methods">

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
