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

func GetUserRepositoryPermission(repositoryNamespaced string, username string, currentUser *Auth.AuthenticatedUser) (Dto.RepositoryPermission, error) {
	return GetRepositoryPermission(repositoryNamespaced, username, "user", currentUser)
}
func GetTeamRepositoryPermission(repositoryNamespaced string, teamname string, currentUser *Auth.AuthenticatedUser) (Dto.RepositoryPermission, error) {
	return GetRepositoryPermission(repositoryNamespaced, teamname, "team", currentUser)
}

/*
GetRepositoryPermission Wrapper for getting a user or team permission on a repository
*/
func GetRepositoryPermission(repositoryNamespaced string, username string, kind string, currentUser *Auth.AuthenticatedUser) (Dto.RepositoryPermission, error) {
	logger.Info("[Permission Service] Get Repository Permission")
	logger.Debug("Repository: %s", repositoryNamespaced)
	logger.Debug("With kind filters: %s", kind)
	logger.Debug("With name filters: %s", username)

	// Split repositoryNamespaced into namespace and name
	namespace, name, err := Common.SplitRepositoryNamespaced(repositoryNamespaced)
	if err != nil {
		logger.Warning("Invalid repository namespaced: %s", repositoryNamespaced)
		return Dto.RepositoryPermission{}, Errors.RepositoryInvalid(repositoryNamespaced)
	}

	// Check if the namespace (org or user) exists
	if namespace != nil {
		_, err = Repositories.GetUserOrOrganizationByName(*namespace)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user or organization found with name: %s", *namespace)
				return Dto.RepositoryPermission{}, Errors.RepositoryNamespaceNotFound(*namespace)
			default:
				logger.Error("Error retrieving repository  from database: %s", err.Error())
				return Dto.RepositoryPermission{}, err
			}
		}
	}

	// Check if the repository exits
	repoExist, err := Repositories.FindRepositoryByNameAndNamespace(name, namespace)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No repository '%s' found", repositoryNamespaced)
			return Dto.RepositoryPermission{}, Errors.RepositoryNotFound(repositoryNamespaced)
		default:
			logger.Error("Error retrieving repository  from database: %s", err.Error())
			return Dto.RepositoryPermission{}, err
		}
	}

	// Get User or Team Permission
	var permissionModel Models.RepositoryPermission
	if kind == "user" {
		permissionModel, err = Repositories.GetUserRepositoryPermissionByUsername(repoExist.ID, username)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user permission found for '%s': '%s'", kind, username)
				return Dto.RepositoryPermission{}, Errors.PermissionNotFound(kind, username)
			default:
				logger.Error("Error retrieving repository permission from database: %s", err.Error())
				return Dto.RepositoryPermission{}, err
			}
		}
	} else {
		permissionModel, err = Repositories.GetTeamRepositoryPermissionByTeamname(repoExist.ID, username)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No team permission found for '%s': '%s'", kind, username)
				return Dto.RepositoryPermission{}, Errors.PermissionNotFound(kind, username)
			default:
				logger.Error("Error retrieving repository permission from database: %s", err.Error())
				return Dto.RepositoryPermission{}, err
			}
		}
	}

	// Convert model to dto
	permission := Common.ConvertRepositoryPermissionModelToDto(permissionModel, kind)

	return permission, nil
}

func UpdateUserRepositoryPermission(repositoryNamespaced string, username string, updatePermission Dto.UpdateRepositoryPermission, currentUser *Auth.AuthenticatedUser) (Dto.RepositoryPermission, error) {
	return updateRepositoryPermission(repositoryNamespaced, username, updatePermission, "user", currentUser)
}

func UpdateTeamRepositoryPermission(repositoryNamespaced string, teamname string, updatePermission Dto.UpdateRepositoryPermission, currentUser *Auth.AuthenticatedUser) (Dto.RepositoryPermission, error) {
	return updateRepositoryPermission(repositoryNamespaced, teamname, updatePermission, "team", currentUser)
}

/*
updateRepositoryPermission Wrapper for update repository permission for a user or team
*/
func updateRepositoryPermission(repositoryNamespaced string, username string, updatePermission Dto.UpdateRepositoryPermission, kind string, currentUser *Auth.AuthenticatedUser) (Dto.RepositoryPermission, error) {
	logger.Info("[Permission Service] Update Repository Permission")
	logger.Debug("Repository: %s", repositoryNamespaced)
	logger.Debug("With kind filters: %s", kind)
	logger.Debug("With name filters: %s", username)
	logger.Debug("With updatePermission: %s", updatePermission.Role)

	// Validate inputs
	err := Common.ValidateUpdateRepositoryPermission(updatePermission)
	if err != nil {
		return Dto.RepositoryPermission{}, err
	}

	// Split repositoryNamespaced into namespace and name
	namespace, name, err := Common.SplitRepositoryNamespaced(repositoryNamespaced)
	if err != nil {
		logger.Warning("Invalid repository namespaced: %s", repositoryNamespaced)
		return Dto.RepositoryPermission{}, Errors.RepositoryInvalid(repositoryNamespaced)
	}

	// Check if the namespace (org or user) exists
	if namespace != nil {
		_, err = Repositories.GetUserOrOrganizationByName(*namespace)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user or organization found with name: %s", *namespace)
				return Dto.RepositoryPermission{}, Errors.RepositoryNamespaceNotFound(*namespace)
			default:
				logger.Error("Error retrieving repository  from database: %s", err.Error())
				return Dto.RepositoryPermission{}, err
			}
		}
	}

	// Check if the repository exits
	repoExist, err := Repositories.FindRepositoryByNameAndNamespace(name, namespace)
	if err != nil {
		switch err.Error() {
		case "record not found":
			logger.Warning("No repository '%s' found", repositoryNamespaced)
			return Dto.RepositoryPermission{}, Errors.RepositoryNotFound(repositoryNamespaced)
		default:
			logger.Error("Error retrieving repository  from database: %s", err.Error())
			return Dto.RepositoryPermission{}, err
		}
	}

	// Get User or Team Permission
	var permissionModel Models.RepositoryPermission
	if kind == "user" {
		permissionModel, err = Repositories.GetUserRepositoryPermissionByUsername(repoExist.ID, username)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No user permission found for '%s': '%s'", kind, username)
				return Dto.RepositoryPermission{}, Errors.PermissionNotFound(kind, username)
			default:
				logger.Error("Error retrieving repository permission from database: %s", err.Error())
				return Dto.RepositoryPermission{}, err
			}
		}
	} else {
		permissionModel, err = Repositories.GetTeamRepositoryPermissionByTeamname(repoExist.ID, username)
		if err != nil {
			switch err.Error() {
			case "record not found":
				logger.Warning("No team permission found for '%s': '%s'", kind, username)
				return Dto.RepositoryPermission{}, Errors.PermissionNotFound(kind, username)
			default:
				logger.Error("Error retrieving repository permission from database: %s", err.Error())
				return Dto.RepositoryPermission{}, err
			}
		}
	}

	// Permission exists, update them
	permissionModel.RoleId = Common.GetRoleIdFromRoleName(updatePermission.Role)
	permissionModel.Role = Models.Role{ID: permissionModel.RoleId, Name: updatePermission.Role} // just for avoiding recall db for getting the updated model
	err = Repositories.UpdateRepositoryPermission(permissionModel)
	if err != nil {
		logger.Error("Error updating repository permission: %s", err.Error())
		return Dto.RepositoryPermission{}, err
	}

	// Convert model to dto
	permission := Common.ConvertRepositoryPermissionModelToDto(permissionModel, kind)

	return permission, nil
}
