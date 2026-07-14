package Api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quay-go-api/Common/Errors"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func permissionController() {
	registerRepositorySubRoute(http.MethodGet, "permissions/team", listRepositoryTeamPermission)
	registerRepositorySubRoute(http.MethodGet, "permissions/user", listRepositoryUserPermission)
	registerRepositorySubRoute(http.MethodGet, "permissions/user/:username", getUserRepositoryPermission)
	registerRepositorySubRoute(http.MethodGet, "permissions/team/:teamname", getTeamRepositoryPermission)
	registerRepositorySubRoute(http.MethodPatch, "permissions/user/:username", updateUserRepositoryPermission)
	registerRepositorySubRoute(http.MethodPatch, "permissions/team/:teamname", updateTeamRepositoryPermission)
}

// listRepositoryTeamPermission List teams permission on a repository
// @Description List teams permission on a repository
// @Summary List teams permission on a repository
// @Tags Permission
// @Param repository path string true "Repository name in the format namespace/repository"
// @Success 200 {object} []Dto.RepositoryPermission
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/repository/{repository}/permissions/team [get]
func listRepositoryTeamPermission(c *gin.Context, repositoryNamespaced string) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.ReadRepo})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	repository, err := Services.ListRepositoryTeamPermission(repositoryNamespaced, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(200, repository)
}

// listRepositoryUserPermission List users permission on a repository
// @Description List users permission on a repository
// @Summary List users permission on a repository
// @Tags Permission
// @Param repository path string true "Repository name in the format namespace/repository"
// @Success 200 {object} []Dto.RepositoryPermission
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/repository/{repository}/permissions/user [get]
func listRepositoryUserPermission(c *gin.Context, repositoryNamespaced string) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.ReadRepo})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	repository, err := Services.ListRepositoryUserPermission(repositoryNamespaced, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(200, repository)
}

// getUserRepositoryPermission Get a user's permission on a repository
// @Description Get a user's permission on a repository
// @Summary Get a user's permission on a repository
// @Tags Permission
// @Param repository path string true "Repository name in the format namespace/repository"
// @Param username path string true "Username"
// @Success 200 {object} Dto.RepositoryPermission
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/repository/{repository}/permissions/user/{username} [get]
func getUserRepositoryPermission(c *gin.Context, repositoryNamespaced string) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.ReadRepo})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	username := c.Param("username")

	permission, err := Services.GetUserRepositoryPermission(repositoryNamespaced, username, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(200, permission)
}

// getTeamRepositoryPermission Get a team's permission on a repository
// @Description Get a team's permission on a repository
// @Summary Get a team's permission on a repository
// @Tags Permission
// @Param repository path string true "Repository name in the format namespace/repository"
// @Param teamname path string true "Team name"
// @Success 200 {object} Dto.RepositoryPermission
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/repository/{repository}/permissions/team/{teamname} [get]
func getTeamRepositoryPermission(c *gin.Context, repositoryNamespaced string) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.ReadRepo})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	teamname := c.Param("teamname")

	permission, err := Services.GetTeamRepositoryPermission(repositoryNamespaced, teamname, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(200, permission)
}

// updateUserRepositoryPermission Update user repository permission
// @Description Update user repository permission
// @Summary Update user repository permission
// @Tags Permission
// @Accept json
// @Param repository path string true "Repository name in the format namespace/repository"
// @Param username path string true "Username"
// @Param update body Dto.UpdateRepositoryPermission true "Repository details to change"
// @Success 200 {object} Dto.Repository
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/repository/{repository}/permissions/user/{username} [patch]
func updateUserRepositoryPermission(c *gin.Context, repositoryNamespaced string) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.ReadRepo})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	username := c.Param("username")

	var updatePermission Dto.UpdateRepositoryPermission
	if err := c.BindJSON(&updatePermission); err != nil {
		throwError(c, Errors.RequestBodyInvalid())
		return
	}

	updatedPermission, err := Services.UpdateUserRepositoryPermission(repositoryNamespaced, username, updatePermission, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(http.StatusOK, updatedPermission)
}

// updateTeamRepositoryPermission Update team repository permission
// @Description Update team repository permission
// @Summary Update team repository permission
// @Tags Permission
// @Accept json
// @Param repository path string true "Repository name in the format namespace/repository"
// @Param teamname path string true "Team name"
// @Param update body Dto.UpdateRepositoryPermission true "Repository details to change"
// @Success 200 {object} Dto.Repository
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/repository/{repository}/permissions/team/{teamname} [patch]
func updateTeamRepositoryPermission(c *gin.Context, repositoryNamespaced string) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.ReadRepo})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	teamname := c.Param("teamname")

	var updatePermission Dto.UpdateRepositoryPermission
	if err := c.BindJSON(&updatePermission); err != nil {
		throwError(c, Errors.RequestBodyInvalid())
		return
	}

	updatedPermission, err := Services.UpdateTeamRepositoryPermission(repositoryNamespaced, teamname, updatePermission, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(http.StatusOK, updatedPermission)
}
