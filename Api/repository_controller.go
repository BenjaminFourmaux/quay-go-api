package Api

import (
	"github.com/gin-gonic/gin"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func repositoryController() {
	repository := engine.Group("/api/v1/organization/:orgname/repository")
	{
		repository.Use(authorizedMiddleware)
		repository.GET("/", listOrganizationRepositories)
	}
}

// listOrganizationRepositories List organization's repositories
// @Description List organization's repositories
// @Summary List organization's repositories
// @Tags Repository
// @Param orgname path string true "Name of the organization"
// @Param is_public query bool false "Filter on public repositories"
// @Param is_starred query bool false "Filter on stared repositories"
// @Success 200 {object} []Dto.Repository
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/organization/{orgname}/repository [get]
func listOrganizationRepositories(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.ReadRepo})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	orgname := c.Param("orgname")

	// Get filters from query params
	filters := extractFilters(c)

	listTeams, err := Services.ListOrganizationRepositories(orgname, filters, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(200, listTeams)
}
