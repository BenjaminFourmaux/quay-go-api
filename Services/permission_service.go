package Services

import (
	"quay-go-api/Common"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Repositories"
	"quay-go-api/Services/Auth"
	"quay-go-api/Services/Avatar"
	logger "quay-go-api/Services/Logger"
)

func ListRepositoryTeamPermission(repositoryNamespaced string, currentUser *Auth.AuthenticatedUser) ([]Dto.RepositoryPermission, error) {
	return listRepositoryPermission(repositoryNamespaced, "team", currentUser)
}

func ListRepositoryUserPermission(repositoryNamespaced string, currentUser *Auth.AuthenticatedUser) ([]Dto.RepositoryPermission, error) {
	return listRepositoryPermission(repositoryNamespaced, "user", currentUser)
}

/*
listRepositoryPermission Wrapper for list both team or user permissions on a repository
*/
func listRepositoryPermission(repositoryNamespaced string, kind string, currentUser *Auth.AuthenticatedUser) ([]Dto.RepositoryPermission, error) {
	logger.Info("[Permission Service] List Repository Permissions")
	logger.Debug("Repository: %s", repositoryNamespaced)
	logger.Debug("With kind filters: %s", kind)

	// Split repositoryNamespaced into namespace and name
	namespace, name, err := Common.SplitRepositoryNamespaced(repositoryNamespaced)
	if err != nil {
		logger.Warning("Invalid repository namespaced: %s", repositoryNamespaced)
		return []Dto.RepositoryPermission{}, Errors.RepositoryInvalid(repositoryNamespaced)
	}

	// Check if the namespace (org or user) exists
	if namespace != nil {
		_, err = Repositories.GetUserOrOrganizationByName(*namespace)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user or organization found with name: %s", *namespace)
				return []Dto.RepositoryPermission{}, Errors.RepositoryNamespaceNotFound(*namespace)
			default:
				logger.Error("Error retrieving repository  from database: %s", err.Error())
				return []Dto.RepositoryPermission{}, err
			}
		}
	}

	// Check if the repository exits
	repoExist, err := Repositories.FindRepositoryByNameAndNamespace(name, namespace)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No repository '%s' found", repositoryNamespaced)
			return []Dto.RepositoryPermission{}, Errors.RepositoryNotFound(repositoryNamespaced)
		default:
			logger.Error("Error retrieving repository  from database: %s", err.Error())
			return []Dto.RepositoryPermission{}, err
		}
	}

	// Get Team or User permissions
	permissionsModel, err := Repositories.ListRepositoryPermissions(repoExist.ID, kind)
	if err != nil {
		logger.Error("Error retrieving repository permissions from database: %s", err.Error())
		return []Dto.RepositoryPermission{}, err
	}

	// Convert models into dto
	permissions := []Dto.RepositoryPermission{}
	for _, permissionModel := range permissionsModel {
		permission := Dto.RepositoryPermission{
			Role: permissionModel.Role.Name,
		}

		if kind == "user" {
			permission.Name = permissionModel.User.Username
			permission.Avatar = Avatar.GetAvatarForUser(*permissionModel.User)
		} else if kind == "team" {
			permission.Name = permissionModel.Team.Name
			permission.Avatar = Avatar.GetAvatarForTeam(*permissionModel.Team)
		}

		permissions = append(permissions, permission)
	}

	return permissions, nil
}
