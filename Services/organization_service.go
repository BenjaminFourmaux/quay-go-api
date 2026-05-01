package Services

import (
	"quay-go-api/Common"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Repositories"
	"quay-go-api/Services/Auth"
)

func GetUserOrganization(userId int, userScopes []Auth.Scope) ([]Dto.Organization, error) {
	// Get the current user
	currentUser := Repositories.SelectUserById(userId)

	orgsModel, err := Repositories.GetUserOrganizations(userId)
	if err != nil {
		// TODO throw custom error ?
		return nil, err
	}

	// Convert userModel into Organization Dto
	organizations := Common.ConvertUserModelsToDto(orgsModel, currentUser, userScopes)
	return organizations, nil
}
