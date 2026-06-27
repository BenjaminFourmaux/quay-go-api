package Api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quay-go-api/Entities/Dto"
	"quay-go-api/Services"
	"quay-go-api/Services/Auth"
)

func repositoryController() {
	repository := engine.Group("/api/v1/repository")
	{
		repository.Use(authorizedMiddleware)
		repository.GET("", listRepositories)
		repository.POST("", createRepository)
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

// createRepository Create a repository
// @Description Create a repository
// @Summary Create a repository
// @Tags Repository
// @Accept json
// @Param message body Dto.CreateRepository true "Repository metadata"
// @Success 201 {object} Dto.Repository
// @Failure 401 {object} Errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} Errors.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /api/v1/repository [post]
func createRepository(c *gin.Context) {
	currentUser, hasScopeErr := retrieveCurrentUser(c, []Auth.Scope{Auth.CreateRepo})
	if hasScopeErr != nil {
		throwError(c, hasScopeErr)
		return
	}

	var repositoryToCreate Dto.CreateRepository
	_ = c.BindJSON(&repositoryToCreate)

	newRepository, err := Services.CreateRepository(repositoryToCreate, currentUser)
	if err != nil {
		throwError(c, err)
		return
	}
	c.JSON(http.StatusCreated, newRepository)
}
