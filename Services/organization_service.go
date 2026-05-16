package Services

import (
	"quay-go-api/Common"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Entities/Models"
	"quay-go-api/Repositories"
	"quay-go-api/Services/Auth"
	"quay-go-api/Services/Avatar"
	"slices"
	"strconv"
	"strings"
)

func GetUserOrganizations(currentUser Auth.AuthenticatedUser, filters map[string]string) ([]Dto.UserOrganization, error) {
	// Validating filters
	var filterPublic bool
	if public, ok := filters["is_public"]; ok {
		isPublic, err := strconv.ParseBool(public)
		if err != nil {
			return nil, Errors.InvalidParameterValue("is_public", []string{"true", "false"})
		}
		filterPublic = isPublic
	}

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

	// Apply filters
	if _, ok := filters["is_public"]; ok {
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
			return []Dto.UserOrganization{}, nil
		}
	}

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
	err = Common.ValidateCreateOrganization(organizationMetadata)
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
	// Retrieve organization and check if exists
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
		if err = Common.ValidateUpdateOrganization(organizationMetadata); err != nil {
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

func ListMembersOfOrganization(orgName string, currentUser Auth.AuthenticatedUser) ([]Dto.OrganizationMember, error) {
	// Retrieve organization and check if exists
	organizationModel, err := Repositories.GetOrganizationDetailsByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			return nil, Errors.OrganizationNotFound(orgName)
		default:
			return nil, err
		}
	}

	// Chek if user has the right to see members
	if !Common.HasScope(currentUser.Scopes, Auth.OrgAdmin) &&
		!Common.HasScope(currentUser.Scopes, Auth.SuperUser) &&
		!isUserIsOrgOwner(currentUser.ID, organizationModel) {
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
