package Services

import (
	"fmt"
	"github.com/google/uuid"
	"quay-go-api/Common"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Entities/Models"
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
		if validatedKind := Common.IsValidRepositoryKind(kind); !validatedKind {
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

func CreateRepository(repositoryMetadata Dto.CreateRepository, currentUser Auth.AuthenticatedUser) (Dto.Repository, error) {
	logger.Info("[Repository Service] Create Organization")
	logger.Debug("With dto: %+v", repositoryMetadata)

	// Validate repo
	err := Common.ValidateCreateRepository(repositoryMetadata)
	if err != nil {
		logger.Warning("Validation error: %s", err.Error())
		return Dto.Repository{}, err
	}

	// Check if the repository already exists
	_, err = Repositories.FindRepositoryByNameAndNamespace(repositoryMetadata.Name, repositoryMetadata.Namespace)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Info("No repositories found, continue")
		default:
			logger.Error("Error retrieving repository  from database: %s", err.Error())
			return Dto.Repository{}, err
		}
	} else {
		logger.Info("Repository already exists")
		return Dto.Repository{}, Errors.RepositoryAlreadyExists()
	}

	// Convert dto to model before inserting
	var repositoryToCreate = Models.Repository{
		Name:         repositoryMetadata.Name,
		Description:  repositoryMetadata.Description,
		VisibilityId: Common.InlineIf(repositoryMetadata.IsPublic, 1, 2),        // 1 = public, 2 = private
		BadgeToken:   uuid.New().String(),                                       // Generate a new UUID for the badge token
		KindId:       Common.InlineIf(repositoryMetadata.Kind == "image", 1, 2), // 1 = image, 2 = application
		TrustEnabled: false,
		State:        0, // NORMAL
	}

	// Add Namespace if specified
	if repositoryMetadata.Namespace != nil {
		// Get user/organization
		userOrOrganizationModel, err := Repositories.GetUserOrOrganizationByName(*repositoryMetadata.Namespace)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user or organization found with name: %s", *repositoryMetadata.Namespace)
				return Dto.Repository{}, Errors.RepositoryNamespaceNotFound(*repositoryMetadata.Namespace)
			default:
				logger.Error("Error retrieving repository  from database: %s", err.Error())
				return Dto.Repository{}, err
			}
		}

		repositoryToCreate.NamespaceUserId = &userOrOrganizationModel.ID
	}

	// Insert the new repository to the DB
	createdRepositoryModel, err := Repositories.CreateRepositoryTransaction(repositoryToCreate)
	if err != nil {
		logger.Error("Error creating repository: %s", err.Error())
		return Dto.Repository{}, err
	}
	logger.Info(fmt.Sprintf("Repository created with ID: %d", createdRepositoryModel.ID))

	// Get the new repository with details
	newRepositoryModel, err := Repositories.GetRepositoryById(createdRepositoryModel.ID, currentUser.ID)
	if err != nil {
		logger.Error("Error retrieving repository: %s", err.Error())
		return Dto.Repository{}, err
	}

	// Convert model to dto
	newRepository := Dto.Repository{
		Namespace:   Common.InlineIf(newRepositoryModel.NamespaceUserId != nil, newRepositoryModel.NamespaceUser.Username, ""),
		Name:        newRepositoryModel.Name,
		Description: newRepositoryModel.Description,
		IsPublic:    newRepositoryModel.VisibilityId == 1,
		Kind:        newRepositoryModel.Kind.Name,
		State:       Common.MapRepositoryStateStr(newRepositoryModel.State),
		IsStarred:   false, // not starred, it has just been created
	}

	return newRepository, nil
}
