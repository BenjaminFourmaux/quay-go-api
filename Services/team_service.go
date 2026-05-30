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

func ListTeamsOfOrganization(orgName string, filters map[string]string, currentUser Auth.AuthenticatedUser) ([]Dto.Team, error) {
	logger.Info("[Team Service] List Teams Of Organization")
	logger.Debug("Organization name: %s", orgName)
	logger.Debug("With filters: %+v", filters)

	// Validating filters
	var filterRole string
	var filterName string
	if role, ok := filters["role"]; ok {
		if validatedRole := Common.ValidateTeamRole(role); !validatedRole {
			logger.Warning("Invalid role filter value: %s", role)
			return nil, Errors.InvalidParameterValue("role", []string{"admin", "creator", "member"})
		}
		filterRole = role
	}
	if name, ok := filters["name"]; ok {
		filterName = name
	}

	// Retrieve organization and check if exists
	logger.Info("Retrieving organization details from database")
	organizationModel, err := Repositories.GetOrganizationDetailsByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("Organization not found: %s", orgName)
			return nil, Errors.OrganizationNotFound(orgName)
		default:
			logger.Error("Error retrieving organization details from database: %s", err.Error())
			return nil, err
		}
	}

	// Check if user has the right to see teams
	if !Common.HasScope(currentUser.Scopes, Auth.OrgAdmin) &&
		!Common.HasScope(currentUser.Scopes, Auth.SuperUser) &&
		!isUserIsOrgOwner(currentUser.ID, organizationModel) {
		logger.Warning("List teams denied: user %d has insufficient role on organization %s", currentUser.ID, orgName)
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

	logger.Debug("Teams returned: %d", len(teams))
	return teams, nil
}

func CreateTeam(teamMetadata Dto.CreateTeam, orgName string, currentUser Auth.AuthenticatedUser) (Dto.Team, error) {
	logger.Info("[Team Service] Create Team")
	logger.Debug("Organization name: %s", orgName)
	logger.Debug("With dto: %+v", teamMetadata)

	// Retrieve organization and check if exists
	logger.Info("Retrieving organization details from database")
	organizationModel, err := Repositories.GetOrganizationDetailsByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("Organization not found: %s", orgName)
			return Dto.Team{}, Errors.OrganizationNotFound(orgName)
		default:
			logger.Error("Error retrieving organization details from database: %s", err.Error())
			return Dto.Team{}, err
		}
	}

	// Validate team
	validateError := Common.ValidateTeam(teamMetadata)
	if validateError != nil {
		logger.Warning("Invalid create team payload: %s", validateError.Error())
		return Dto.Team{}, validateError
	}

	// Check if the team already exists
	for _, team := range organizationModel.Teams {
		if team.Name == *teamMetadata.Name {
			logger.Warning("Team already exists: %s", *teamMetadata.Name)
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
	logger.Info("Creating team in database")
	createdTeamModel, err := Repositories.CreateTeam(newTeam)
	if err != nil {
		logger.Error("Error creating team in database: %s", err.Error())
		return Dto.Team{}, err
	}

	// Convert model to dto and return
	createdTeam := Common.ConvertTeamModelToDto(createdTeamModel, currentUser.ID, currentUser.Scopes)

	createdTeam.Role = *teamMetadata.Role
	logger.Success("Team created successfully: %s", createdTeam.Name)

	return createdTeam, nil
}

func GetTeam(orgName string, teamName string, currentUser Auth.AuthenticatedUser) (Dto.Team, error) {
	logger.Info("[Team Service] Get Team")
	logger.Debug("Organization name: %s", orgName)
	logger.Debug("Team name: %s", teamName)

	// Retrieve organization and check if exists
	logger.Info("Retrieving organization details from database")
	organizationModel, err := Repositories.GetOrganizationDetailsByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("Organization not found: %s", orgName)
			return Dto.Team{}, Errors.OrganizationNotFound(orgName)
		default:
			logger.Error("Error retrieving organization details from database: %s", err.Error())
			return Dto.Team{}, err
		}
	}

	// Browse the organization's team and find the team to get
	for _, team := range organizationModel.Teams {
		if team.Name == teamName {
			// Convert model to dto and return
			teamDto := Common.ConvertTeamModelToDto(team, currentUser.ID, currentUser.Scopes)
			logger.Debug("Team found: %s", teamDto.Name)
			return teamDto, nil
		}
	}

	logger.Warning("Team not found: %s", teamName)
	return Dto.Team{}, Errors.TeamNotFound(teamName)
}

func UpdateTeam(teamToUpdate Dto.UpdateTeam, orgName string, teamName string, currentUser Auth.AuthenticatedUser) (Dto.Team, error) {
	logger.Info("[Team Service] Update Team")
	logger.Debug("Organization name: %s", orgName)
	logger.Debug("Team name: %s", teamName)
	logger.Debug("With dto: %+v", teamToUpdate)

	// Retrieve organization and check if exists
	logger.Info("Retrieving organization details from database")
	organizationModel, err := Repositories.GetOrganizationDetailsByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("Organization not found for update team: %s", orgName)
			return Dto.Team{}, Errors.OrganizationNotFound(orgName)
		default:
			logger.Error("Error retrieving organization details from database: %s", err.Error())
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
			logger.Warning("Team not found for update: %s", teamName)
			return Dto.Team{}, Errors.TeamNotFound(teamName)
		}

		// Validate input values
		if teamToUpdate.Role != nil {
			if roleOk := Common.ValidateTeamRole(*teamToUpdate.Role); !roleOk {
				logger.Warning("Invalid role value for update team: %s", *teamToUpdate.Role)
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
		logger.Debug("Team fields to update: %+v", updatedFields)

		logger.Info("Updating team in database")
		if err = Repositories.UpdateTeamFieldsById(teamIdToUpdate, updatedFields); err != nil {
			logger.Error("Error updating team in database: %s", err.Error())
			return Dto.Team{}, err
		} else {
			logger.Info("Retrieving updated team from database")
			updatedTeamModel, err := Repositories.GetTeamById(teamIdToUpdate)
			if err != nil {
				logger.Error("Error retrieving updated team from database: %s", err.Error())
				return Dto.Team{}, err
			}

			// Convert model to dto
			updatedTeam := Common.ConvertTeamModelToDto(updatedTeamModel, currentUser.ID, currentUser.Scopes)
			logger.Success("Team updated successfully: %s", updatedTeam.Name)

			return updatedTeam, nil
		}
	}

	logger.Warning("Update denied: user %d is not owner of organization %s", currentUser.ID, orgName)
	return Dto.Team{}, nil
}

func DeleteTeam(orgName string, teamName string, currentUser Auth.AuthenticatedUser) error {
	logger.Info("[Team Service] Delete Team")
	logger.Debug("Organization name: %s", orgName)
	logger.Debug("Team name: %s", teamName)

	// Get the org (with detail, for user role checking) if exists
	logger.Info("Retrieving organization details from database")
	organizationModel, err := Repositories.GetOrganizationDetailsByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("Organization not found for delete team: %s", orgName)
			return Errors.OrganizationNotFound(orgName)
		default:
			logger.Error("Error retrieving organization details from database: %s", err.Error())
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
				logger.Warning("Team not found for deletion: %s", teamName)
				return Errors.TeamNotFound(teamName)
			}
		}

		logger.Info("Deleting team in database")
		err = Repositories.DeleteTeamTransaction(teamIdToDelete)
		if err != nil {
			logger.Error("Error deleting team in database: %s", err.Error())
			return err
		}
		logger.Success("Team deleted successfully: %s", teamName)
	} else {
		logger.Warning("Delete denied: user %d is not owner of organization %s", currentUser.ID, orgName)
	}
	return nil
}
