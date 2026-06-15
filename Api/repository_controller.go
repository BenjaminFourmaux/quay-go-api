package Api

import (
	"github.com/gin-gonic/gin"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func repositoryController() {
	repository := engine.Group("/api/v1/repository")
	{
		repository.Use(authorizedMiddleware)
		repository.GET("", listRepositories)
	}
}

// listRepositories List repositories
// @Description List repositories
// @Summary List repositories
// @Tags Repository
// @Param namespace query string false "Organization or Username as namespace to filter repositories"
// @Param is_public query bool false "Filter on public repositories"
// @Param is_starred query bool false "Filter on stared repositories"
// @Success 200 {object} []Dto.Repository
// @Failure 400 {object} Errors.ErrorResponse "Bad Request"
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/repository [get]
func listRepositories(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.ReadRepo})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	// Get filters from query params
	filters := extractFilters(c)

	listTeams, err := Services.ListRepositories(filters, &currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(200, listTeams)
}
