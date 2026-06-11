package Services

import (
	"quay-go-api/Common"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Repositories"
	"quay-go-api/Services/Auth"
	logger "quay-go-api/Services/Logger"
)

func ListOrganizationRepositories(orgName string, filters map[string]string, currentUser *Auth.AuthenticatedUser) ([]Dto.Repository, error) {
	logger.Info("[Repository Service] List Repositories Of Organization")
	logger.Debug("Organization name: %s", orgName)
	logger.Debug("With filters: %+v", filters)

	// Validating filters
	var filterPublic bool
	var filterStarred bool
	var hasStarredFilter bool
	var filterKind string
	if public, ok := filters["is_public"]; ok {
		filterPublic = public == "true"
	}
	if starred, ok := filters["is_starred"]; ok {
		filterStarred = starred == "true"
		hasStarredFilter = true
	} else if stared, ok := filters["is_stared"]; ok {
		filterStarred = stared == "true"
		hasStarredFilter = true
	}
	if kind, ok := filters["kind"]; ok {
		if validatedKind := Common.ValidateRepositoryKind(kind); !validatedKind {
			logger.Warning("Invalid kind filter value: %s", kind)
			return nil, Errors.InvalidParameterValue("kind", []string{"image", "application"})
		}
		filterKind = kind
	}

	// Retrieve organization and check if exists
	logger.Info("Retrieving the organization from database")
	organizationModel, err := Repositories.GetOrganizationByName(orgName)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("Organization not found: %s", orgName)
			return nil, Errors.OrganizationNotFound(orgName)
		default:
			logger.Error("Error retrieving organization from database: %s", err.Error())
			return nil, err
		}
	}

	// Retrieve repositories of the organization
	logger.Info("Retrieving organization repositories from database")
	repositoriesModel, err := Repositories.GetOrganizationRepositoriesByOrgId(organizationModel.ID, currentUser.ID)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Debug("No repository for organization: %s", orgName)
			return []Dto.Repository{}, nil
		default:
			logger.Error("Error retrieving repositories from database: %s", err.Error())
			return []Dto.Repository{}, err
		}
	}

	logger.Debug("Found repositories: %d", len(repositoriesModel))

	// Convert Repository Model in DTO
	repositories := []Dto.Repository{}

	for _, repository := range repositoriesModel {
		isStarred := len(repository.Stars) > 0

		// Apply filters
		if _, ok := filters["is_public"]; ok {
			if filterPublic && repository.VisibilityId == 2 { // want public but repo has visibility 2 (private)
				continue
			}
			if !filterPublic && repository.VisibilityId == 1 { // want private but repo has visibility 1 (public)
				continue
			}
		}
		if hasStarredFilter {
			if filterStarred != isStarred {
				continue
			}
		}
		if _, ok := filters["kind"]; ok {
			if filterKind != repository.Kind.Name {
				continue
			}
		}

		repositories = append(repositories, Dto.Repository{
			Namespace:   orgName,
			Name:        repository.Name,
			Description: repository.Description,
			IsPublic:    repository.VisibilityId == 1,
			Kind:        repository.Kind.Name,
			State:       Common.MapRepositoryStateStr(repository.State),
			IsStarred:   isStarred,
		})
	}

	return repositories, nil
}
