package Api

import (
	"github.com/gin-gonic/gin"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func permissionController() {
	registerRepositorySubRoute("/permissions/team", listRepositoryTeamPermission)
	registerRepositorySubRoute("/permissions/user", listRepositoryUserPermission)
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
