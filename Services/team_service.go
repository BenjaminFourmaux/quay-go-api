package Services

import (
	"quay-go-api/Common"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Entities/Models"
	"quay-go-api/Repositories"
	"quay-go-api/Services/Auth"
)

func ListTeamsOfOrganization(orgName string, filters map[string]string, currentUser Auth.AuthenticatedUser) ([]Dto.Team, error) {
	// Validating filters
	var filterRole string
	var filterName string
	if role, ok := filters["role"]; ok {
		if validatedRole := Common.ValidateTeamRole(role); !validatedRole {
			return nil, Errors.InvalidParameterValue("role", []string{"admin", "creator", "member"})
		}
		filterRole = role
	}
	if name, ok := filters["name"]; ok {
		filterName = name
	}

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

	// Check if user has the right to see teams
	if !Common.HasScope(currentUser.Scopes, Auth.OrgAdmin) &&
		!Common.HasScope(currentUser.Scopes, Auth.SuperUser) &&
		!isUserIsOrgOwner(currentUser.ID, organizationModel) {
		return nil, Errors.UnauthorizedInsufficientRole()
	}

	teams := []Dto.Team{}
	for _, team := range organizationModel.Teams {
		// Apply filters
		if (filterRole != "" && team.Role.Name != filterRole) || (filterName != "" && team.Name != filterName) {
			continue
		}

		teams = append(teams, Common.ConvertTeamModelToDto(team, currentUser.ID, currentUser.Scopes))
	}
	return teams, nil
}

func CreateTeam(teamMetadata Dto.CreateTeam, orgName string, currentUser Auth.AuthenticatedUser) (Dto.Team, error) {
	// Retrieve organization and check if exists
	organizationModel, err := Repositories.GetOrganizationDetailsByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			return Dto.Team{}, Errors.OrganizationNotFound(orgName)
		default:
			return Dto.Team{}, err
		}
	}

	// Validate team
	validateError := Common.ValidateTeam(teamMetadata)
	if validateError != nil {
		return Dto.Team{}, validateError
	}

	// Check if the team already exists
	for _, team := range organizationModel.Teams {
		if team.Name == *teamMetadata.Name {
			return Dto.Team{}, Errors.TeamAlreadyExists()
		}
	}

	// Convert dto to model
	newTeam := Models.Team{
		Name:           *teamMetadata.Name,
		OrganizationId: organizationModel.ID,
	}

	// Optional fields
	if teamMetadata.Description != nil {
		newTeam.Description = *teamMetadata.Description
	}
	if teamMetadata.Role != nil {
		newTeam.RoleId = Common.GetTeamRoleIdFromRoleName(*teamMetadata.Role)
	}

	// Insert into the DB
	createdTeamModel, err := Repositories.CreateTeam(newTeam)
	if err != nil {
		return Dto.Team{}, err
	}

	// Convert model to dto and return
	createdTeam := Common.ConvertTeamModelToDto(createdTeamModel, currentUser.ID, currentUser.Scopes)

	createdTeam.Role = *teamMetadata.Role

	return createdTeam, nil
}

func GetTeam(orgName string, teamName string, currentUser Auth.AuthenticatedUser) (Dto.Team, error) {
	// Retrieve organization and check if exists
	organizationModel, err := Repositories.GetOrganizationDetailsByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			return Dto.Team{}, Errors.OrganizationNotFound(orgName)
		default:
			return Dto.Team{}, err
		}
	}

	// Browse the organization's team and find the team to get
	for _, team := range organizationModel.Teams {
		if team.Name == teamName {
			// Convert model to dto and return
			teamDto := Common.ConvertTeamModelToDto(team, currentUser.ID, currentUser.Scopes)
			return teamDto, nil
		}
	}

	return Dto.Team{}, Errors.TeamNotFound(teamName)
}

func UpdateTeam(teamToUpdate Dto.UpdateTeam, orgName string, teamName string, currentUser Auth.AuthenticatedUser) (Dto.Team, error) {
	// Retrieve organization and check if exists
	organizationModel, err := Repositories.GetOrganizationDetailsByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			return Dto.Team{}, Errors.OrganizationNotFound(orgName)
		default:
			return Dto.Team{}, err
		}
	}

	// If the user is owners of the organization so it can delete teams of this organization
	if isUserIsOrgOwner(currentUser.ID, organizationModel) {
		// Check if the organization's team exists
		var teamIdToUpdate int
		for _, team := range organizationModel.Teams {
			if team.Name == teamName {
				teamIdToUpdate = team.ID
				break
			}
		}

		// If team not found
		if teamIdToUpdate == 0 {
			return Dto.Team{}, Errors.TeamNotFound(teamName)
		}

		// Validate input values
		if teamToUpdate.Role != nil {
			if roleOk := Common.ValidateTeamRole(*teamToUpdate.Role); !roleOk {
				return Dto.Team{}, Errors.InvalidParameterValue("role", []string{"admin", "creator", "member"})
			}
		}

		// Select fields to update
		mappings := map[string]Common.UpdateFieldMapping{}

		if teamToUpdate.Role != nil {
			mappings["Role"] = Common.UpdateFieldMapping{
				ModelFieldName: "RoleId",
				Value:          Common.GetTeamRoleIdFromRoleName(*teamToUpdate.Role),
			}
		}
		if teamToUpdate.Description != nil {
			mappings["Description"] = Common.UpdateFieldMapping{
				ModelFieldName: "Description",
				Value:          *teamToUpdate.Description,
			}
		}

		updatedFields := Common.BuildUpdatedFields[Models.Team](teamToUpdate, mappings)

		if err = Repositories.UpdateTeamFieldsById(teamIdToUpdate, updatedFields); err != nil {
			return Dto.Team{}, err
		} else {
			updatedTeamModel, err := Repositories.GetTeamById(teamIdToUpdate)
			if err != nil {
				return Dto.Team{}, err
			}

			// Convert model to dto
			updatedTeam := Common.ConvertTeamModelToDto(updatedTeamModel, currentUser.ID, currentUser.Scopes)

			return updatedTeam, nil
		}
	}

	return Dto.Team{}, nil
}

func DeleteTeam(orgName string, teamName string, currentUser Auth.AuthenticatedUser) error {
	// Get the org (with detail, for user role checking) if exists
	organizationModel, err := Repositories.GetOrganizationDetailsByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			return Errors.OrganizationNotFound(orgName)
		default:
			return err
		}
	}

	// If the user is owners of the organization so it can delete teams of this organization
	if isUserIsOrgOwner(currentUser.ID, organizationModel) {
		// Check if the organization's team exists
		var teamIdToDelete int
		for _, team := range organizationModel.Teams {
			if team.Name == teamName {
				teamIdToDelete = team.ID
				break
			} else {
				return Errors.TeamNotFound(teamName)
			}
		}

		err = Repositories.DeleteTeamTransaction(teamIdToDelete)
		if err != nil {
			return err
		}
	}
	return nil
}
