package Services

import (
	"fmt"
	"quay-go-api/Common"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Repositories"
	"quay-go-api/Services/Auth"
	logger "quay-go-api/Services/Logger"
)

func ListRepositories(filters map[string]string, currentUser *Auth.AuthenticatedUser) ([]Dto.Repository, error) {
	logger.Info("[Repository Service] List Repositories")
	logger.Debug("With filters: %+v", filters)

	// Validating filters
	var filterNamespace *string
	var filterPublic bool
	var filterStarred bool
	var hasStarredFilter bool
	var filterKind string
	if ns, ok := filters["namespace"]; ok && ns != "" {
		filterNamespace = &ns
	}
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

	var orgOrUserId *int
	if filterNamespace != nil {
		// Find the user/organization id
		logger.Info(fmt.Sprintf("Filter namespace provided: %s. Retrieving user/organization id from database", *filterNamespace))
		orgOrUserModel, err := Repositories.GetUserOrOrganizationByName(*filterNamespace)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Debug("No user or organization found with name: %s", *filterNamespace)
				return []Dto.Repository{}, Errors.UserOrOrganizationNotFound(*filterNamespace)
			default:
				logger.Error("Error retrieving user/organization from database: %s", err.Error())
				return []Dto.Repository{}, err
			}
		}
		orgOrUserId = &orgOrUserModel.ID
		logger.Debug("Found user/organization id: %d for name: %s", *orgOrUserId, *filterNamespace)
		logger.Debug("Is an organization ? %t", orgOrUserModel.Organization)
	}

	// Retrieve repositories
	logger.Info("Retrieving all repositories from database")
	repositoriesModel, err := Repositories.SelectRepositories(currentUser.ID, orgOrUserId)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Debug("No repository find")
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
			Namespace:   Common.InlineIf(repository.NamespaceUserId != nil, repository.NamespaceUser.Username, ""),
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
