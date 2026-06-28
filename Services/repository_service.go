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
	"sort"
	"time"
)

const maxDaysIn3Months = 90

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
			Kind:        repository.Kind.Name,
			IsPublic:    repository.VisibilityId == 1,
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
	createdRepositoryModel, err := Repositories.CreateRepositoryTransaction(repositoryToCreate, currentUser.ID)
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
		Kind:        newRepositoryModel.Kind.Name,
		IsPublic:    newRepositoryModel.VisibilityId == 1,
		State:       Common.MapRepositoryStateStr(newRepositoryModel.State),
		IsStarred:   false, // not starred, it has just been created
	}

	return newRepository, nil
}

func GetRepository(repositoryNamespaced string, filters map[string]string, currentUser *Auth.AuthenticatedUser) (Dto.RepositoryDetails, error) {
	logger.Info("[Repository Service] Get Repository")
	logger.Debug("Repository: %s", repositoryNamespaced)
	logger.Debug("With filters: %+v", filters)

	// Validating filters
	var filterIncludeTags bool = false
	var filterIncludeStats bool = false
	if it, ok := filters["include_tags"]; ok && it != "" {
		filterIncludeTags = it == "true"
	}
	if is, ok := filters["include_stats"]; ok && is != "" {
		filterIncludeStats = is == "true"
	}

	// Split repositoryNamespaced into namespace and name
	namespace, name, err := Common.SplitRepositoryNamespaced(repositoryNamespaced)
	if err != nil {
		logger.Warning("Invalid repository namespaced: %s", repositoryNamespaced)
		return Dto.RepositoryDetails{}, Errors.RepositoryInvalid(repositoryNamespaced)
	}

	// Check if the namespace (org or user) exists
	if namespace != nil {
		_, err = Repositories.GetUserOrOrganizationByName(*namespace)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user or organization found with name: %s", *namespace)
				return Dto.RepositoryDetails{}, Errors.RepositoryNamespaceNotFound(*namespace)
			default:
				logger.Error("Error retrieving repository  from database: %s", err.Error())
				return Dto.RepositoryDetails{}, err
			}
		}
	}

	// Check if the repository exits
	repoExist, err := Repositories.FindRepositoryByNameAndNamespace(name, namespace)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No repository '%s' found", repositoryNamespaced)
			return Dto.RepositoryDetails{}, Errors.RepositoryNotFound(repositoryNamespaced)
		default:
			logger.Error("Error retrieving repository  from database: %s", err.Error())
			return Dto.RepositoryDetails{}, err
		}
	}

	// Get repository with details
	repositoryModel, err := Repositories.GetRepositoryById(repoExist.ID, currentUser.ID)
	if err != nil {
		logger.Error("Error retrieving repository: %s", err.Error())
		return Dto.RepositoryDetails{}, err
	}

	hasWritePermission := checkRepositoryUserPermission(repositoryModel.ID, currentUser.ID, "write")
	hasWritePermission = hasWritePermission && repositoryModel.State == 0 // NORMAL
	hasAdminPermission := checkRepositoryUserPermission(repositoryModel.ID, currentUser.ID, "admin")

	// Convert model to dto
	repository := Dto.RepositoryDetails{
		Name:           repositoryModel.Name,
		Namespace:      Common.InlineIf(repositoryModel.NamespaceUserId != nil, repositoryModel.NamespaceUser.Username, ""),
		Description:    repositoryModel.Description,
		Kind:           repositoryModel.Kind.Name,
		IsPublic:       repositoryModel.VisibilityId == 1,
		IsOrganization: repositoryModel.NamespaceUserId != nil && repositoryModel.NamespaceUser.Organization,
		IsStarred:      len(repositoryModel.Stars) > 0,
		StatusToken:    repositoryModel.BadgeToken,
		TrustEnabled:   false,
		TagExpirationS: Common.InlineIf(repositoryModel.NamespaceUserId != nil, repositoryModel.NamespaceUser.RemovedTagExpirationS, 1209600), // 14 days default
		State:          Common.MapRepositoryStateStr(repositoryModel.State),
		CanWrite:       hasWritePermission,
		CanAdmin:       hasAdminPermission,
	}

	// Apply filters
	if filterIncludeTags {
		tagsModel, err := Repositories.GetTagsFromRepository(repositoryModel.ID)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Info("No tags found for repository: %s", repositoryModel.ID)
				repository.Tags = []Dto.RepositoryTag{}
			default:
				logger.Error("Error retrieving tags from database: %s", err.Error())
				return Dto.RepositoryDetails{}, err
			}
		} else {
			// Add and convert tags model to dto
			tags := []Dto.RepositoryTag{}
			for _, tag := range tagsModel {
				tags = append(tags, Dto.RepositoryTag{
					Name:           tag.Name,
					Size:           *tag.Manifest.LayersCompressedSize,
					LastModified:   time.UnixMilli(tag.LifetimeStartMs),
					ManifestDigest: tag.Manifest.Digest,
				})
			}
			repository.Tags = tags
		}
	}

	if filterIncludeStats {
		countsModel, err := Repositories.GetCountsFromRepository(repositoryModel.ID)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Info("No stats found for repository: %s", repositoryModel.ID)
				repository.Stats = []Dto.RepositoryStats{}
			default:
				logger.Error("Error retrieving stats from database: %s", err.Error())
				return Dto.RepositoryDetails{}, err
			}
		} else {
			stats := make([]Dto.RepositoryStats, 0, len(countsModel))
			foundDates := make(map[string]struct{}, len(countsModel))

			for _, count := range countsModel {
				countDate := time.Date(count.Date.Year(), count.Date.Month(), count.Date.Day(), 0, 0, 0, 0, time.UTC)
				stats = append(stats, Dto.RepositoryStats{
					Date:  countDate,
					Count: count.Count,
				})
				key := fmt.Sprintf("%d/%d", countDate.Month(), countDate.Day())
				foundDates[key] = struct{}{}
			}

			// Fill in any missing stats with zeros.
			now := time.Now().UTC()
			for day := 1; day < maxDaysIn3Months; day++ {
				dayDate := now.AddDate(0, 0, -day)
				key := fmt.Sprintf("%d/%d", dayDate.Month(), dayDate.Day())
				if _, ok := foundDates[key]; !ok {
					stats = append(stats, Dto.RepositoryStats{
						Date:  time.Date(dayDate.Year(), dayDate.Month(), dayDate.Day(), 0, 0, 0, 0, time.UTC),
						Count: 0,
					})
				}
			}

			// Order stats by date descending (the newest first)
			sort.Slice(stats, func(i, j int) bool {
				return stats[i].Date.After(stats[j].Date)
			})

			repository.Stats = stats
		}
	}

	return repository, nil
}

func UpdateRepository(repositoryNamespaced string, repositoryMetadata Dto.UpdateRepository, currentUser Auth.AuthenticatedUser) (Dto.Repository, error) {
	logger.Info("[Repository Service] Update Repository")
	logger.Debug("Updating repository: %s", repositoryNamespaced)

	// Split repositoryNamespaced into namespace and name
	namespace, name, err := Common.SplitRepositoryNamespaced(repositoryNamespaced)
	if err != nil {
		logger.Warning("Invalid repository namespaced: %s", repositoryNamespaced)
		return Dto.Repository{}, Errors.RepositoryInvalid(repositoryNamespaced)
	}

	// Check if the namespace (org or user) exists
	if namespace != nil {
		_, err = Repositories.GetUserOrOrganizationByName(*namespace)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user or organization found with name: %s", *namespace)
				return Dto.Repository{}, Errors.RepositoryNamespaceNotFound(*namespace)
			default:
				logger.Error("Error retrieving repository  from database: %s", err.Error())
				return Dto.Repository{}, err
			}
		}
	}

	// Check if the repository exits
	repoExist, err := Repositories.FindRepositoryByNameAndNamespace(name, namespace)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No repository '%s' found", repositoryNamespaced)
			return Dto.Repository{}, Errors.RepositoryNotFound(repositoryNamespaced)
		default:
			logger.Error("Error retrieving repository  from database: %s", err.Error())
			return Dto.Repository{}, err
		}
	}

	// Update model
	repoExist.Description = repositoryMetadata.Description

	updatedRepositoryModel, err := Repositories.UpdateRepository(repoExist)
	if err != nil {
		logger.Error("Error updating repository: %s", err.Error())
		return Dto.Repository{}, err
	}

	// Convert model to dto
	updatedRepository := Dto.Repository{
		Name:        updatedRepositoryModel.Name,
		Namespace:   Common.InlineIf(updatedRepositoryModel.NamespaceUserId != nil, updatedRepositoryModel.NamespaceUser.Username, ""),
		Description: updatedRepositoryModel.Description,
		Kind:        Common.InlineIf(updatedRepositoryModel.KindId == 1, "image", "application"),
		State:       Common.MapRepositoryStateStr(updatedRepositoryModel.State),
		IsPublic:    updatedRepositoryModel.VisibilityId == 1,
		IsStarred:   len(updatedRepositoryModel.Stars) > 0,
	}

	return updatedRepository, nil
}

// <editor-fold desc="Private Methods"

func checkRepositoryUserPermission(repositoryId int, userId int, role string) bool {
	// Get user permission on the repository
	userPermission, err := Repositories.GetRepositoryUserPermission(repositoryId, userId)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Debug("No permission found for user %d on repository %d", userId, repositoryId)
			return false
		default:
			logger.Error("Error retrieving permission from database: %s", err.Error())
			return false
		}
	}

	switch role {
	case "admin":
		return userPermission.Role.Name == "admin"
	case "write":
		return userPermission.Role.Name == "admin" || userPermission.Role.Name == "write"
	case "read":
		return userPermission.Role.Name == "admin" || userPermission.Role.Name == "read"
	default:
		logger.Warning("Invalid role: %s", role)
		return false
	}
}

// </editor-fold>
